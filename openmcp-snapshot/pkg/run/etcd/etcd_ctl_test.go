package etcd

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestInsertEtcd(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	projectPath := strings.Split(path.Dir(filename), "/")
	pathLen := len(projectPath) - 2
	shellPath := path.Join(strings.Join(projectPath[:pathLen], "/"), "util", "template", "fileName")
	fmt.Println(shellPath)
}

// Get : 키 검색
func TestGetEtcd(t *testing.T) {
	// 검색됨
	ret, err := GetEtcd("1605230624-cluster1-PersistentVolume-snapshot-test-pv")
	if err != nil {
		t.Errorf("Error %s", err)
	}
	if ret != "1605230624-cluster1-PersistentVolume-snapshot-test-pv" {
		t.Errorf("Error value %s", ret)
	}

	// 존재하지 않는 키
	ret, err = GetEtcd("key_1")
	if err != nil {

	}
	if ret == "val_5" {
		t.Errorf("Error value %s", ret)
	}
}
