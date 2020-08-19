module kchtest

go 1.14

require (
	admiralty.io/multicluster-controller v0.1.0
	clientset v0.0.0-00010101000000-000000000000
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	resource-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.5.0
	sigs.k8s.io/kubefed v0.3.0

)

replace (
	clientset => ../clientset
	resource-controller => ../resource-controller
	sync-controller => ../sync-controller
)
