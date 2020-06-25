module openmcp/openmcp

go 1.14

require (
	admiralty.io/multicluster-controller v0.6.0
	admiralty.io/multicluster-service-account v0.5.0 // indirect
	bitbucket.org/ww/goautoneg v0.0.0-20120707110453-75cd24fc2f2c // indirect
	github.com/NYTimes/gziphandler v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dmportella/powerdns v0.0.0-20170512083725-7cb21fd59fc6 // indirect
	github.com/emicklei/go-restful v2.12.0+incompatible
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6
	github.com/fatih/structs v1.1.0 // indirect
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/go-openapi/spec v0.19.8
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.4.0
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hth0919/resourcecollector v0.0.0-20200410074403-e831e209e52f // indirect
	github.com/influxdata/influxdb v1.8.0
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/joeig/go-powerdns/v2 v2.3.3 // indirect
	github.com/kubernetes-incubator/custom-metrics-apiserver v0.0.0-20200618121405-54026617ec44
	github.com/kubernetes-sigs/federation-v2 v0.0.10 // indirect
	github.com/kubernetes/autoscaler v0.0.0-20191021073337-3f137fde4f66
	github.com/mittwald/go-powerdns v0.5.0
	github.com/operator-framework/operator-sdk v0.10.0 // indirect
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0
	gonum.org/v1/netlib v0.0.0-20190331212654-76723241ea4e // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0
	k8s.io/api v0.18.4
	k8s.io/apimachinery v0.18.4
	k8s.io/apiserver v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-registry v0.0.6 // indirect
	k8s.io/component-base v0.18.2
	k8s.io/klog v1.0.0
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/metrics v0.18.2
	k8s.io/sample-controller v0.18.4
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubefed v0.3.0
	sigs.k8s.io/structured-merge-diff v1.0.1-0.20191108220359-b1b620dd3f06 // indirect
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.10.0
	k8s.io/client-go => k8s.io/client-go v0.18.2
)
