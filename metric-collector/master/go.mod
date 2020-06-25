module openmcp-metric-collector

go 1.12

require (
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/influxdata/influxdb v1.8.0
	github.com/influxdata/influxdb-client-go v1.1.0 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.23.0
	k8s.io/apiextensions-apiserver v0.18.2 // indirect
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	resource-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.6.0 // indirect
	sigs.k8s.io/kubefed v0.1.0-rc6
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
)

replace (
	github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.0
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	openmcp-dns-controller => ../../openmcp-dns-controller
	resource-controller => ../../resource-controller
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
	sync-controller => ../../sync-controller
)
