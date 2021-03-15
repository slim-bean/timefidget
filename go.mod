module timefidget

go 1.15

require (
	github.com/cortexproject/cortex v1.6.1-0.20210204145131-7dac81171c66
	github.com/go-kit/kit v0.10.0
	github.com/go-logfmt/logfmt v0.5.0
	github.com/grafana/loki v1.6.2-0.20210227183507-877f524c36bf
	github.com/magefile/mage v1.11.0
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.15.0
	github.com/weaveworks/common v0.0.0-20210112142934-23c8d7fa6120
)

// We can't upgrade to grpc 1.30.0 until go.etcd.io/etcd will support it.
// Solves google.golang.org/grpc/naming: module google.golang.org/grpc@latest found (v1.35.0), but does not contain package google.golang.org/grpc/naming
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

// Solves ambiguous import
replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

// Keeping this same as Cortex to avoid dependency issues.
replace k8s.io/client-go => k8s.io/client-go v0.19.4

replace k8s.io/api => k8s.io/api v0.19.4

// Use fork of gocql that has gokit logs and Prometheus metrics.
replace github.com/gocql/gocql => github.com/grafana/gocql v0.0.0-20200605141915-ba5dc39ece85

// >v1.2.0 has some conflict with prometheus/alertmanager. Hence prevent the upgrade till it's fixed.
replace github.com/satori/go.uuid => github.com/satori/go.uuid v1.2.0

// Same as Cortex
// Using a 3rd-party branch for custom dialer - see https://github.com/bradfitz/gomemcache/pull/86
replace github.com/bradfitz/gomemcache => github.com/themihai/gomemcache v0.0.0-20180902122335-24332e2d58ab

// Fix errors like too many arguments in call to "github.com/go-openapi/errors".Required
//   have (string, string)
//   want (string, string, interface {})
replace github.com/go-openapi/errors => github.com/go-openapi/errors v0.19.4

replace github.com/go-openapi/validate => github.com/go-openapi/validate v0.19.8
