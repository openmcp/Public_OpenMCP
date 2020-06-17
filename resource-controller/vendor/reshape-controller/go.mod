module reshape-controller

go 1.12

require (
	admiralty.io/multicluster-controller v0.1.0
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/googleapis/gnostic v0.4.2 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	go.uber.org/zap v1.15.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	k8s.io/apiextensions-apiserver v0.18.3 // indirect
	k8s.io/client-go v12.0.0+incompatible // indirect
	sigs.k8s.io/controller-runtime v0.3.0
	sigs.k8s.io/kubefed v0.1.0-rc6
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
)
