module resource-controller

go 1.12

require (
	admiralty.io/multicluster-controller v0.1.0
	github.com/emicklei/go-restful v2.12.0+incompatible // indirect
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/go-openapi/spec v0.19.8
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.2
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hth0919/resourcecollector v0.0.0-20200410074403-e831e209e52f // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kr/pty v1.1.5 // indirect
	github.com/kubernetes/autoscaler v0.0.0-20191021073337-3f137fde4f66
	github.com/mittwald/go-powerdns v0.4.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8 // indirect
	golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3 // indirect
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f // indirect
	golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.23.0
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	openmcp-dns-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/kubefed v0.1.0-rc6
	sigs.k8s.io/yaml v1.2.0 // indirect
	sync-controller v0.0.0-00010101000000-000000000000

)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20190602132728-7075c07e78bf
	openmcp-dns-controller => ../openmcp-dns-controller
	sync-controller => ../sync-controller
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20181117043124-c2090bec4d9b
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20180711000925-0cf8f7e6ed1d
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)
