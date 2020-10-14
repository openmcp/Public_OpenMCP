module openmcp/openmcp

go 1.14

require (
	admiralty.io/multicluster-controller v0.6.0
	cloud.google.com/go v0.63.0 // indirect
	github.com/auth0/go-jwt-middleware v0.0.0-20200507191422-d30d7b9ece63
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getlantern/deepcopy v0.0.0-20160317154340-7f45deb8130a
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1 // indirect
	github.com/go-openapi/spec v0.19.8
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/influxdata/influxdb v1.8.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/kubernetes/autoscaler v0.0.0-20191021073337-3f137fde4f66
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mittwald/go-powerdns v0.5.0
	github.com/olekukonko/tablewriter v0.0.0-20170122224234-a0225b3f23b5
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.17.3
	k8s.io/apimachinery v0.18.6
	k8s.io/autoscaler/vertical-pod-autoscaler v0.0.0-20200710055348-ddb8ec8980e1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/component-base v0.18.6
	k8s.io/klog v1.0.0
	k8s.io/kube-aggregator v0.18.6
	k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6
	k8s.io/metrics v0.18.2
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kubefed v0.3.0

)

replace (
	admiralty.io/multicluster-controller => admiralty.io/multicluster-controller v0.1.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.10.0
	github.com/pkg/errors => github.com/pkg/errors v0.9.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.3
	//k8s.io/client-go => k8s.io/client-go v0.18.2
	k8s.io/client-go => k8s.io/client-go v0.17.3

	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.0
)
