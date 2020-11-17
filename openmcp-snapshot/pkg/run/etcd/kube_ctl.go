package etcd

import (

	//"context"

	"fmt"
	"openmcp/openmcp/openmcp-snapshot/pkg/run/etcd/resources"
	"path/filepath"

	// "openmcp/openmcp/snapshot-operator/pkg/run/etcd/resources"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// ExistsResource indicates an error occurred unescaping a field selector
type ExistsResource struct {
	resourceName string
}

func (e ExistsResource) Error() string {
	return fmt.Sprintf("deployments.apps %s already exists", e.resourceName)
}

// GetResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func GetResourceJSON(clientset *kubernetes.Clientset, resourceType string, resourceName string, resourceNamespace string) (string, error) {
	var ret string
	var err error
	var client resources.Resource
	switch resourceType {
	case "Deployment", "deployment", "deploy":
		client = resources.Deployment{}
	case "Service", "service", "svc":
		client = resources.Service{}
	case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
		client = resources.PersistentVolumeClaim{}
	case "PersistentVolume", "persistentvolume", "pv":
		client = resources.PersistentVolume{}
	}
	fmt.Printf("resourceType : %s, resourceName : %s, resourceNamespace: %s\n", resourceType, resourceName, resourceNamespace)
	ret, err = client.GetJSON(clientset, resourceName, resourceNamespace)
	if err != nil {
		switch err {
		case clientcmd.ErrNoContext:
			// 매칭되는 이름이 없음.
			return "", clientcmd.ErrNoContext
		default:
			// 일반적인 에러처리 코드
			return "", err
			//return "", ErrUndefined
		}
	}

	return ret, err
}

// CreateResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여
func CreateResourceJSON(clientset *kubernetes.Clientset, resourceType string, resourceInfoJSON string) (bool, error) {
	var ret bool
	var err error
	var client resources.Resource

	fmt.Println("*** CreateResourceJSON()")
	switch resourceType {
	case "Deployment", "deployment", "deploy":
		client = resources.Deployment{}

	case "Service", "service", "svc":
		client = resources.Service{}

	case "PersistentVolumeClaim", "persistentvolumeclaim", "pvc":
		client = resources.PersistentVolumeClaim{}

		//case "pv":
		//	ret = GetPVYaml(clientset, resourceName)

	case "PersistentVolume", "persistentvolume", "pv":
		client = resources.PersistentVolume{}

		//case "pv":
		//	ret = GetPVYaml(clientset, resourceName)

	}
	ret, err = client.CreateResourceForJSON(clientset, resourceInfoJSON)
	if err != nil {
		switch err {
		//case clientcmd.ErrNoContext:
		//	// 매칭되는 이름이 없음.
		//	return false, clientcmd.ErrNoContext
		default:
			// 일반적인 에러처리 코드
			return false, err
		}
	}
	fmt.Println(ret)

	return ret, err
}

// InitKube : 초기화
func InitKube() *kubernetes.Clientset {

	var a string
	a = "a"
	var b []byte
	b = []byte(a)

	clientcmd.NewClientConfigFromBytes(b)

	//var kubeconfig *string
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	/*if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	*/
	fmt.Println("10.0.0.222:8001")
	config, err := clientcmd.BuildConfigFromFlags("10.0.0.222:8001", kubeconfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n----------------")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func int32Ptr(i int32) *int32 { return &i }
