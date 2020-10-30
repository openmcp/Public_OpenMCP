package etcd

import (
	"fmt"
	"testing"
)

// TestGetYaml : getYaml 테스트
func TestGetResourceJson(t *testing.T) {
	clientset := InitKube()

	var resourceType string
	var resourceName string
	var resourceNamespace string

	var val string
	var err error
	var isSuccess bool
	var etcdErr error

	/*
		// 실패 케이스1
		resourceType := "deployment"
		resourceName := "test2"
		resourceNamespace := "default"
		val, err := GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)

		if err == nil {
			t.Error("Error : ", err) // 에러 발생
		}

		isSuccess, etcdErr := InsertEtcd(resourceType+"-"+resourceName, val)
		if etcdErr != nil {
			t.Error("Error : ", etcdErr) // 에러 발생
		}
		if !isSuccess {
			t.Error("Insert Etcd Fail") // 에러 발생
		}

		// 성공 케이스1
		resourceType = "deployment"
		resourceName = "test-dp"
		resourceNamespace = "default"

		val, err = GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)
		t.Log(resourceType)
		if err != nil {
			t.Error("Error : ", err) // 에러 발생
		}

		isSuccess, etcdErr = InsertEtcd(resourceType+"-"+resourceName, val)
		if etcdErr != nil {
			t.Error("Error : ", etcdErr) // 에러 발생
		}
		if !isSuccess {
			t.Error("Insert Etcd Fail") // 에러 발생
		}

		// 성공 케이스1
		resourceType = "svc"
		resourceName = "test-svc"
		resourceNamespace = "default"

		val, err = GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)
		t.Log(resourceType)
		if err != nil {
			t.Error("Error : ", err) // 에러 발생
		}

		isSuccess, etcdErr = InsertEtcd(resourceType+"-"+resourceName, val)
		if etcdErr != nil {
			t.Error("Error : ", etcdErr) // 에러 발생
		}
		if !isSuccess {
			t.Error("Insert Etcd Fail") // 에러 발생
		}

		// 성공 케이스1
		resourceType = "pvc"
		resourceName = "test-pvc"
		resourceNamespace = "default"

		val, err = GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)
		t.Log(resourceType)
		if err != nil {
			t.Error("Error : ", err) // 에러 발생
		}

		isSuccess, etcdErr = InsertEtcd(resourceType+"-"+resourceName, val)
		if etcdErr != nil {
			t.Error("Error : ", etcdErr) // 에러 발생
		}
		if !isSuccess {
			t.Error("Insert Etcd Fail") // 에러 발생
		}
	*/
	// 성공 케이스1
	resourceType = "pv"
	resourceName = "test-pv-snapshot"
	resourceNamespace = "default"

	val, err = GetResourceJSON(clientset, resourceType, resourceName, resourceNamespace)
	t.Log(resourceType)
	if err != nil {
		t.Error("Error : ", err) // 에러 발생
	}

	isSuccess, etcdErr = InsertEtcd(resourceType+"-"+resourceName, val)
	if etcdErr != nil {
		t.Error("Error : ", etcdErr) // 에러 발생
	}
	if !isSuccess {
		t.Error("Insert Etcd Fail") // 에러 발생
	}

}

// TestGetYaml : getYaml 테스트
func TestCreateResourceJSON(t *testing.T) {

	// 성공 케이스1
	resourceType := "pv"
	resourceName := "test-pv"

	A(t, resourceType, resourceName)

	// 성공 케이스1
	resourceType = "pvc"
	resourceName = "test-pvc"

	//A(t, resourceType, resourceName)

	// 성공 케이스1
	resourceType = "svc"
	resourceName = "test-svc"

	//A(t, resourceType, resourceName)

	// 성공 케이스1
	resourceType = "deployment"
	resourceName = "test-dp"

	//A(t, resourceType, resourceName)
}

func A(t *testing.T, resourceType string, resourceName string) {
	t.Log(resourceType)

	// ETCD에서 json 가져오기
	jsonStr, etcdErr := GetEtcd(resourceType + "-" + resourceName)
	if etcdErr != nil {
		t.Errorf("Error %s", etcdErr)
	}

	clientset := InitKube()
	// defer 는 fatal 전에
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Create Resource ERROR", r)
		}
	}()
	isSuccess, err := CreateResourceJSON(clientset, resourceType, jsonStr)
	if err != nil {

		t.Error("Error : ", err) // 에러 발생
	} else if !isSuccess {
		t.Error("false") // 에러 발생
	}

}
