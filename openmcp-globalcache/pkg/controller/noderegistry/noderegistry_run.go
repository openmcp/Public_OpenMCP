package noderegistry

import (
	"fmt"

	v1alpha1 "openmcp/openmcp/apis/globalcache/v1alpha1"
	nodeapi "openmcp/openmcp/openmcp-globalcache/pkg/run/dist"
)

// Run : 실제 로직단
func (r *reconciler) Run(instance *v1alpha1.NodeRegistry) (bool, error) {

	fmt.Println("\n[Command]] :" + instance.Spec.Command)
	var registryManager nodeapi.RegistryManager
	err := registryManager.Init(instance.Spec.ClusterName, cm)
	if err != nil {
		return false, err
	}
	err = registryManager.SetNodeLabelSync()
	if err != nil {
		return false, err
	}

	// push, pull - nodeName 이 없을 경우 Cluster 단위 명령
	switch instance.Spec.Command {
	case "pull":
		if instance.Spec.NodeName == "" {
			err = registryManager.CreatePullJobForCluster(instance.Spec.ImageName, instance.Spec.TagName)
			if err != nil {
				return false, err
			}
		} else {
			err = registryManager.CreatePullJob(instance.Spec.NodeName, instance.Spec.ImageName, instance.Spec.TagName)
			if err != nil {
				return false, err
			}
		}
	case "push":
		err = registryManager.CreatePushJobForCluster(instance.Spec.ImageName, instance.Spec.TagName)
		if err != nil {
			return false, err
		}

	//case "tagList":
	default:
		return false, fmt.Errorf("Command is not valid")
	}

	return true, nil
}
