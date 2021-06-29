package api

import (
	"openmcp/openmcp/openmcp-globalcache/pkg/run/dist"
	"openmcp/openmcp/util/clusterManager"
)

// 외부 통신용 api. 인터페이스 정의서에 정의된 내용은 모두 여기에 포함된다.

//DistributeRepositoryAgent -
//case : join
var cm *clusterManager.ClusterManager

func DistributeRepositoryAgent(clusterName string) error {
	var manager dist.RegistryManager
	err := manager.Init(clusterName, cm)
	if err != nil {
		return err
	}
	return nil
}

//DeleteRepositoryAgent -
//case : unjoin
func DeleteRepositoryAgent(clusterName string) error {
	var manager dist.RegistryManager
	err := manager.Init(clusterName, cm)
	if err != nil {
		return err
	}
	err = manager.DeleteRegistryAgent("nanumdev3")
	if err != nil {
		return err
	}
	return nil
}
