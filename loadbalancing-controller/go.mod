module loadbalancing-controller

go 1.12

require (
	admiralty.io/multicluster-controller v0.1.0
	github.com/aeden/traceroute v0.0.0-20181124220833-147686d9cb0f // indirect
	github.com/go-openapi/spec v0.19.8
	github.com/golang/protobuf v1.4.2
	github.com/oschwald/geoip2-golang v1.4.0 // indirect
	github.com/umahmood/haversine v0.0.0-20151105152445-808ab04add26 // indirect
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.23.0
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	k8s.io/sample-controller v0.0.0-20191017070449-ab9e95689d58
	resource-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/kubefed v0.1.0-rc6
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20190602132728-7075c07e78bf
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	openmcp-dns-controller => ../openmcp-dns-controller
	resource-controller => ../resource-controller
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
	sync-controller => ../sync-controller
)
