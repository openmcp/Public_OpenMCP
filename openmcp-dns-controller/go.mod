module openmcp-dns-controller

go 1.12

require (
	admiralty.io/multicluster-controller v0.1.0
	admiralty.io/multicluster-service-account v0.5.0 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dmportella/powerdns v0.0.0-20170512083725-7cb21fd59fc6 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-openapi/spec v0.19.8
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a // indirect
	github.com/joeig/go-powerdns/v2 v2.3.3 // indirect
	github.com/kubernetes-sigs/federation-v2 v0.0.10 // indirect
	github.com/mittwald/go-powerdns v0.4.0
	github.com/operator-framework/operator-sdk v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/cluster-registry v0.0.6 // indirect
	k8s.io/kube-openapi v0.0.0-20200121204235-bf4fb3bd569c
	k8s.io/sample-controller v0.0.0-20191017070449-ab9e95689d58 // indirect

	sync-controller v0.0.0-00010101000000-000000000000
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
	resource-controller => ../resource-controller
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
