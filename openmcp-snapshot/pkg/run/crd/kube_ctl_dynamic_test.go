package crd

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/api/errors"
)

// TestGetYaml : getYaml 테스트
func TestDynamicGetResourceJson(t *testing.T) {
	clientset := DynamicInitKube()

	// 케이스1
	resourceType := "deployment"
	resourceName := "test"
	resourceNamespace := "default"
	val, err := DynamicGetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)

	if err != nil {
		t.Error("Error : ", err) // 에러 발생
	}

	isSuccess, etcdErr := InsertEtcd(resourceType+"-"+resourceName, val)
	if etcdErr != nil {
		t.Error("Error : ", etcdErr) // 에러 발생
	}
	if !isSuccess {
		t.Error("Insert Etcd Fail") // 에러 발생
	}
}

// TestGetYaml : getYaml 테스트
func TestDynamicCreateResourceJSON(t *testing.T) {

	// 성공 케이스1
	resourceType := "deployment"
	resourceName := "test"

	B(t, resourceType, resourceName)
}

func B(t *testing.T, resourceType string, resourceName string) {
	t.Log(resourceType)

	// ETCD에서 json 가져오기
	jsonStr, etcdErr := GetEtcd(resourceType + "-" + resourceName)
	if etcdErr != nil {
		t.Errorf("Error %s", etcdErr)
	}

	clientset := DynamicInitKube()
	// defer 는 fatal 전에
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Create Resource ERROR", r)
		}
	}()

	var err error
	var isSuccess bool
	isSuccess, err = DynamicCreateResourceJSON(clientset, resourceType, jsonStr)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			t.Error("Error - IsAlreadyExists : ", err) // 에러 발생
		} else {
			t.Error("Error : ", err) // 에러 발생
		}
	} else if !isSuccess {
		t.Error("false") // 에러 발생
	}

}
