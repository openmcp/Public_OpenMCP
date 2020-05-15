package CRM

import (
	"k8s.io/client-go/kubernetes"
)

type ClientSet struct {
	clientSet        *kubernetes.Clientset
}