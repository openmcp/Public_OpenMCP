module clientset

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v12.0.0+incompatible
	resource-controller v0.0.0-00010101000000-000000000000
	sigs.k8s.io/structured-merge-diff v0.0.0-20190525122527-15d366b2352e // indirect

)

replace (
	clientset => ../clientset
	resource-controller => ../resource-controller
	sync-controller => ../sync-controller
)
