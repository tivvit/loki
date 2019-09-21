module github.com/grafana/loki

go 1.12

require (
	cloud.google.com/go v0.34.0
	contrib.go.opencensus.io/exporter/ocagent v0.2.0
	github.com/Azure/azure-sdk-for-go v26.3.0+incompatible
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/Azure/go-autorest v11.5.1+incompatible
	github.com/Microsoft/go-winio v0.4.12
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf
	github.com/aws/aws-sdk-go v1.16.17
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/blang/semver v3.5.1+incompatible
	github.com/bmatcuk/doublestar v1.1.1
	github.com/bradfitz/gomemcache v0.0.0-20180710155616-bc664df96737
	github.com/census-instrumentation/opencensus-proto v0.1.0
	github.com/cespare/xxhash v0.0.0-20190104012619-3b82fb7d1867
	github.com/containerd/fifo v0.0.0-20190226154929-a9fb20d87448
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea
	github.com/cortexproject/cortex v0.0.0-20190708133539-ef492f6bbafb
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v0.0.0-20190607191414-238f8eaa31aa
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.0-20181218153428-b84716841b82
	github.com/docker/go-plugins-helpers v0.0.0-20181025120712-1e6269c305b8
	github.com/docker/go-units v0.4.0
	github.com/etcd-io/bbolt v1.3.1-etcd.8
	github.com/fatih/color v1.7.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/fsouza/fake-gcs-server v0.0.0-20190102220127-b364c791f57a
	github.com/go-kit/kit v0.8.0
	github.com/go-logfmt/logfmt v0.4.0
	github.com/gocql/gocql v0.0.0-20181124151448-70385f88b28b
	github.com/gogo/googleapis v1.1.0
	github.com/gogo/protobuf v1.2.0
	github.com/gogo/status v1.0.3
	github.com/golang/protobuf v1.2.0
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
	github.com/googleapis/gax-go v0.0.0-20181219185031-c8a15bac9b9f
	github.com/googleapis/gnostic v0.2.0
	github.com/gophercloud/gophercloud v0.1.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.6.2
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed
	github.com/hashicorp/consul v1.4.0
	github.com/hashicorp/go-cleanhttp v0.5.0
	github.com/hashicorp/go-rootcerts v0.0.0-20160503143440-6bb64b370b90
	github.com/hashicorp/golang-lru v0.5.0
	github.com/hashicorp/serf v0.8.1
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af
	github.com/json-iterator/go v1.1.5
	github.com/klauspost/compress v1.7.4
	github.com/klauspost/cpuid v1.2.1
	github.com/konsorten/go-windows-terminal-sequences v1.0.1
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.4
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/miekg/dns v1.0.4
	github.com/mitchellh/go-homedir v1.0.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223
	github.com/mwitkow/go-grpc-middleware v1.0.0
	github.com/oklog/ulid v1.3.1
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/opentracing-contrib/go-grpc v0.0.0-20180928155321-4b5a12d3ff02
	github.com/opentracing-contrib/go-stdlib v0.0.0-20181222025249-77df8e8e70b4
	github.com/opentracing/opentracing-go v1.0.2
	github.com/petar/GoLLRB v0.0.0-20130427215148-53be0d36a84c
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v0.0.0-20190307145337-c5e14697eaa7
	github.com/prometheus/client_model v0.0.0-20190109181635-f287a105a20e
	github.com/prometheus/common v0.1.0
	github.com/prometheus/procfs v0.0.0-20190104112138-b1a0a9a36d74
	github.com/prometheus/prometheus v0.0.0-20190607092147-e23fa22233cf
	github.com/prometheus/tsdb v0.6.1
	github.com/samuel/go-zookeeper v0.0.0-20180130194729-c4fab1ac1bec
	github.com/sercand/kuberesolver v2.1.0+incompatible
	github.com/shurcooL/httpfs v0.0.0-20181222201310-74dc9339e414
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.3.0
	github.com/stretchr/objx v0.1.1
	github.com/stretchr/testify v1.3.0
	github.com/tonistiigi/fifo v0.0.0-20190226154929-a9fb20d87448
	github.com/uber/jaeger-client-go v0.0.0-20190309061702-9774b4be4743
	github.com/uber/jaeger-lib v2.0.0+incompatible
	github.com/weaveworks/common v0.0.0-20190508083148-5bf824591a65
	github.com/weaveworks/promrus v1.2.0
	go.opencensus.io v0.18.0
	golang.org/x/crypto v0.0.0-20190103213133-ff983b9c42bc
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc
	golang.org/x/oauth2 v0.0.0-20190110195249-fd3eaa146cbb
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.0.0-20190109145017-48ac38b7c8cb
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c
	google.golang.org/api v0.0.0-20190607001116-5213b8090861
	google.golang.org/appengine v1.4.0
	google.golang.org/genproto v0.0.0-20190110221437-6909d8a4a91b
	google.golang.org/grpc v1.21.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/fsnotify.v1 v1.4.7
	gopkg.in/fsnotify/fsnotify.v1 v1.4.7
	gopkg.in/inf.v0 v0.9.1
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20190602205700-9b8cae951d65
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go v0.0.0-20190409021438-1a26190bd76a
	k8s.io/klog v0.1.0
	k8s.io/utils v0.0.0-20190529001817-6999998975a7
	sigs.k8s.io/yaml v1.1.0
)
