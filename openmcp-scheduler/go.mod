module openmcpscheduler

go 1.12

require (
	admiralty.io/multicluster-controller v0.1.0
	admiralty.io/multicluster-service-account v0.5.0 // indirect
	cloud.google.com/go v0.38.0 // indirect
	github.com/emicklei/go-restful v2.11.0+incompatible // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a // indirect
	github.com/go-httpproxy/httpproxy v0.0.0-20180417134941-6977c68bf38e // indirect
	github.com/go-openapi/spec v0.19.4
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hth0919/resourcecollector v0.0.0-20200410074403-e831e209e52f // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/ip2location/ip2location-go v8.3.0+incompatible // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/kubernetes-sigs/federation-v2 v0.0.10 // indirect
	github.com/kubernetes/autoscaler v0.0.0-20191021073337-3f137fde4f66 // indirect
	github.com/operator-framework/operator-sdk v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0 // indirect
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cluster-registry v0.0.6 // indirect
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d
	k8s.io/sample-controller v0.0.0-20191017070449-ab9e95689d58
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
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)
