module openmcp-metric-collector

go 1.12

require (
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.4.0
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/influxdata/influxdb v1.8.0
	github.com/influxdata/influxdb-client-go v1.1.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/grpc v1.28.0
	google.golang.org/protobuf v1.21.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/apiextensions-apiserver v0.18.2 // indirect
	k8s.io/apimachinery v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0 // indirect
	sigs.k8s.io/controller-runtime v0.6.0 // indirect
	sigs.k8s.io/kubefed v0.1.0-rc6
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
)
