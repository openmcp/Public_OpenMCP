package globalregistry

import (
	"fmt"
	v1alpha1 "openmcp/openmcp/apis/globalcache/v1alpha1"
	globalapi "openmcp/openmcp/openmcp-globalcache/pkg/run/registry"
)

// Run : 실제 로직단
func (r *reconciler) Run(instance *v1alpha1.GlobalRegistry) (bool, error) {

	fmt.Println("\n[Command]] :" + instance.Spec.Command)
	var registryManager globalapi.RegistryManager
	//delete - tagName null 일 경우 전체 삭체, list, tagList
	switch instance.Spec.Command {
	case "delete":
		if instance.Spec.TagName == "" {
			_, err := registryManager.DeleteGlobalRegistryImageAllTag(instance.Spec.ImageName)
			if err != nil {
				return false, err
			}
		} else {
			_, err := registryManager.DeleteGlobalRegistryImage(instance.Spec.ImageName, instance.Spec.TagName)
			if err != nil {
				return false, err
			}
		}
		// tags, _ := registryManager.ListGlobalRegistryImageTag(instance.Spec.ImageName)
		// fmt.Println("1111 : ", tags)
		// if len(tags) == 0 {

		// }
	//case "list":
	//case "tagList":
	default:
		return false, fmt.Errorf("Command is not valid")
	}

	return true, nil
}
