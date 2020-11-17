package etcd

import "testing"

func TestInsertEtcd(t *testing.T) {
	InsertEtcd("key_5", "val_5")
}

// Get : 키 검색
func TestGetEtcd(t *testing.T) {
	// 검색됨
	ret, err := GetEtcd("key_5")
	if err != nil {
		t.Errorf("Error %s", err)
	}
	if ret != "val_5" {
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
