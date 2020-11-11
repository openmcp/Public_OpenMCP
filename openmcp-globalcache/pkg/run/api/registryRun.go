package api

import (
	"openmcp/openmcp/openmcp-globalcache/pkg/run/registry"
)

//ListGlobalRegistryImage 는 글로벌 Registry 의 이미지 리스트를 출력하는 함수.
func ListGlobalRegistryImage() ([]string, error) {
	var r registry.RegistryManager
	result, err := r.ListGlobalRegistryImage()
	if err != nil {
		return nil, err
	}
	return result, nil
}

//ListGlobalRegistryImageTag 는 글로벌 Registry image의 tag 리스트를 출력하는 함수.
func ListGlobalRegistryImageTag(imageName string) ([]string, error) {
	var r registry.RegistryManager
	result, err := r.ListGlobalRegistryImageTag(imageName)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//DeleteGlobalRegistryImage 는 글로벌 Registry 의 image 를 삭제하는 함수 (태그를 지정해야 한다.)
// param : imagename -
func DeleteGlobalRegistryImage(imageName string, tag string) error {
	var r registry.RegistryManager
	_, err := r.DeleteGlobalRegistryImage(imageName, tag)
	if err != nil {
		return err
	}
	return nil
}

//DeleteGlobalRegistryImageAllTag 는 글로벌 Registry 의 image 를 삭제하는 함수. (모든 태그 삭제.)
// 파르르..
func DeleteGlobalRegistryImageAllTag(imageName string) error {
	var r registry.RegistryManager
	_, err := r.DeleteGlobalRegistryImageAllTag(imageName)
	if err != nil {
		return err
	}
	return nil
}

//
func PushGlobalRegistryImage(clusterName string, nodeName string, imageName string, Tag string) error {
	//TODo
	return nil
}

//
func PullGlobalRegistryImage(clusterName string, nodeName string, imageName string, Tag string) error {
	//TODO
	return nil
}
