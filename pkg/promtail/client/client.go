package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/grafana/loki/pkg/promtail/api"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"

	"github.com/grafana/loki/pkg/helpers"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
)

const contentType = "application/x-protobuf"
const maxErrMsgLen = 1024

var (
	encodedBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "encoded_bytes_total",
		Help:      "Number of bytes encoded and ready to send.",
	}, []string{"id", "host"})
	sentBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "sent_bytes_total",
		Help:      "Number of bytes sent.",
	}, []string{"id", "host"})
	droppedBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "dropped_bytes_total",
		Help:      "Number of bytes dropped because failed to be sent to the ingester after all retries.",
	}, []string{"id", "host"})
	sentEntries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "sent_entries_total",
		Help:      "Number of log entries sent to the ingester.",
	}, []string{"id", "host"})
	droppedEntries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "promtail",
		Name:      "dropped_entries_total",
		Help:      "Number of log entries dropped because failed to be sent to the ingester after all retries.",
	}, []string{"id", "host"})
	queueLen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "promtail",
		Name:      "output_queue_len",
		Help:      "Number of entries waiting to be sent.",
	}, []string{"id", "host"})
	requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "promtail",
		Name:      "request_duration_seconds",
		Help:      "Duration of send requests.",
	}, []string{"id", "status_code", "host"})
)

func init() {
	prometheus.MustRegister(encodedBytes)
	prometheus.MustRegister(sentBytes)
	prometheus.MustRegister(droppedBytes)
	prometheus.MustRegister(sentEntries)
	prometheus.MustRegister(droppedEntries)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(queueLen)
}

// Client pushes entries to Loki and can be stopped
type Client interface {
	api.EntryHandler
	// Stop goroutine sending batch of entries.
	Stop()
	GetOutChanLen() int
}

// Client for pushing logs in snappy-compressed protos over HTTP.
type client struct {
	logger  log.Logger
	cfg     Config
	client  *http.Client
	quit    chan struct{}
	once    sync.Once
	entries chan entry
	wg      sync.WaitGroup
	out     chan outBatch

	externalLabels model.LabelSet
}

type entry struct {
	labels model.LabelSet
	logproto.Entry
}

type outBatch struct {
	data  []byte
	count int
}

// New makes a new Client.
func New(cfg Config, logger log.Logger) (Client, error) {
	c := &client{
		logger:  log.With(logger, "component", "client", "host", cfg.URL.Host),
		cfg:     cfg,
		quit:    make(chan struct{}),
		entries: make(chan entry),
		out:     make(chan outBatch, cfg.OutBufferCap),

		externalLabels: cfg.ExternalLabels.LabelSet,
	}

	err := cfg.Client.Validate()
	if err != nil {
		return nil, err
	}

	c.client, err = config.NewClientFromConfig(cfg.Client, "promtail", false)
	if err != nil {
		return nil, err
	}

	c.client.Timeout = cfg.Timeout

	c.wg.Add(1)
	go c.run()
	c.wg.Add(1)
	go c.sender()
	return c, nil
}

func (c *client) run() {
	batch := map[model.Fingerprint]*logproto.Stream{}
	batchSize := 0
	maxWait := time.NewTicker(c.cfg.BatchWait)

	defer func() {
		if len(batch) > 0 {
			c.sendBatch(batch)
		}

		c.wg.Done()
	}()

	for {
		select {
		case <-c.quit:
			return

		case e := <-c.entries:
			if batchSize+len(e.Line) > c.cfg.BatchSize {
				c.sendBatch(batch)
				batchSize = 0
				batch = map[model.Fingerprint]*logproto.Stream{}
			}

			batchSize += len(e.Line)
			fp := e.labels.FastFingerprint()
			stream, ok := batch[fp]
			if !ok {
				stream = &logproto.Stream{
					Labels: e.labels.String(),
				}
				batch[fp] = stream
			}
			stream.Entries = append(stream.Entries, e.Entry)

		case <-maxWait.C:
			if len(batch) > 0 {
				c.sendBatch(batch)
				batchSize = 0
				batch = map[model.Fingerprint]*logproto.Stream{}
			}
		}
	}
}

func (c *client) sendBatch(batch map[model.Fingerprint]*logproto.Stream) {
	buf, entriesCount, err := encodeBatch(batch)
	if err != nil {
		level.Error(c.logger).Log("msg", "error encoding batch", "error", err)
		return
	}
	bufBytes := float64(len(buf))
	encodedBytes.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Add(bufBytes)
	c.out <- outBatch{
		data:  buf,
		count: entriesCount,
	}
}

func (c *client) sender() {
	var ob outBatch
	var bufBytes float64
	var err error
	var buf []byte
	var entriesCount int
	var ql int
	quit := false
	for {
		select {
		case <-c.quit:
			quit = true
		case ob = <-c.out:
			buf = ob.data
			entriesCount = ob.count
			bufBytes = float64(len(buf))
			ctx := context.Background()
			backOff := util.NewBackoff(ctx, c.cfg.BackoffConfig)
			var status int
			for backOff.Ongoing() {
				start := time.Now()
				status, err = c.send(ctx, buf)
				requestDuration.WithLabelValues(c.cfg.Id, strconv.Itoa(status), c.cfg.URL.Host).Observe(time.Since(start).Seconds())

				if err == nil {
					sentBytes.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Add(bufBytes)
					sentEntries.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Add(float64(entriesCount))
					break
				}

				// Only retry 500s and connection-level errors.
				if status > 0 && status/100 != 5 {
					break
				}

				level.Warn(c.logger).Log("msg", "error sending batch, will retry", "status", status, "error", err)
				backOff.Wait()
			}

			if err != nil {
				level.Error(c.logger).Log("msg", "final error sending batch", "status", status, "error", err)
				droppedBytes.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Add(bufBytes)
				droppedEntries.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Add(float64(entriesCount))
			}
		}
		ql = len(c.out)
		queueLen.WithLabelValues(c.cfg.Id, c.cfg.URL.Host).Set(float64(ql))

		// send all messages and then quit (when exiting)
		if quit && ql == 0 {
			c.wg.Done()
			return
		}
	}
}

func encodeBatch(batch map[model.Fingerprint]*logproto.Stream) ([]byte, int, error) {
	req := logproto.PushRequest{
		Streams: make([]*logproto.Stream, 0, len(batch)),
	}

	entriesCount := 0
	for _, stream := range batch {
		req.Streams = append(req.Streams, stream)
		entriesCount += len(stream.Entries)
	}

	buf, err := proto.Marshal(&req)
	if err != nil {
		return nil, 0, err
	}
	buf = snappy.Encode(nil, buf)
	return buf, entriesCount, nil
}

func (c *client) send(ctx context.Context, buf []byte) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()
	req, err := http.NewRequest("POST", c.cfg.URL.String(), bytes.NewReader(buf))
	if err != nil {
		return -1, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		return -1, err
	}
	defer helpers.LogError("closing response body", resp.Body.Close)

	if resp.StatusCode/100 != 2 {
		scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("server returned HTTP status %s (%d): %s", resp.Status, resp.StatusCode, line)
	}
	return resp.StatusCode, err
}

// Stop the client.
func (c *client) Stop() {
	c.once.Do(func() { close(c.quit) })
	c.wg.Wait()
}

// Handle implement EntryHandler; adds a new line to the next batch; send is async.
func (c *client) Handle(ls model.LabelSet, t time.Time, s string) error {
	if len(c.externalLabels) > 0 {
		ls = c.externalLabels.Merge(ls)
	}

	c.entries <- entry{ls, logproto.Entry{
		Timestamp: t,
		Line:      s,
	}}
	return nil
}

func (c *client) GetOutChanLen() int {
	return len(c.out)
}
