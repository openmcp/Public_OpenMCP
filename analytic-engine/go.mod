module openmcp-analytic-engine

go 1.14

require (
	admiralty.io/multicluster-controller v0.6.0 // indirect
	github.com/golang/protobuf v1.4.0
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/influxdata/influxdb v1.8.0
	github.com/oschwald/geoip2-golang v1.4.0
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.22.0
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	resource-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.6.0 // indirect
	sigs.k8s.io/kubefed v0.1.0-rc6
)

replace (
	admiralty.io/multicluster-controller => admiralty.io/multicluster-controller v0.1.0
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0
	k8s.io/api => k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	k8s.io/client-go => k8s.io/client-go v0.17.3
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	openmcp-dns-controller => ../openmcp-dns-controller
	resource-controller => ../resource-controller
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
	sigs.k8s.io/kubefed => sigs.k8s.io/kubefed v0.3.0
	sync-controller => ../sync-controller
)
