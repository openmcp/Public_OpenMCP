package crd

import (
	"fmt"
	"path/filepath"

	"openmcp/openmcp/openmcp-snapshot/pkg/run/etcd/resources"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	//"context"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// DynamicGetResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여 작성
func DynamicGetResourceJSON(clientset dynamic.Interface, resourceType string, resourceName string, resourceNamespace string) (string, error) {
	var ret string
	var err error
	var client = resources.Dynamic{}

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
	fmt.Println(ret)

	return ret, err
}

// DynamicCreateResourceJSON : https://mingrammer.com/gobyexample/json/ 를 참조하여
func DynamicCreateResourceJSON(clientset dynamic.Interface, resourceType string, resourceInfoJSON string) (bool, error) {
	var ret bool
	var err error
	var client = resources.Dynamic{}

	fmt.Println("*** CreateResourceJSON()")

	ret, err = client.CreateResourceForJSON(clientset, resourceInfoJSON)
	//err1 := &errors.StatusError{ErrStatus: *err}
	if err != nil {
		return false, err
	}
	fmt.Println(ret)

	return ret, err
}

// DynamicInitKube : 초기화
func DynamicInitKube() dynamic.Interface {
	//TODO - etcd 에서 가져오는 것으로 바꿔야함.
	//var kubeconfig *string
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	/*if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	*/
	config, err := clientcmd.BuildConfigFromFlags("10.0.0.222:8001", kubeconfig)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n----------------")
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}
