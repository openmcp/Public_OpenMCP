package registry

import (
	"strings"

	"openmcp/openmcp/openmcp-globalcache/pkg/utils"
)

//RegistryManager dd
type RegistryManager struct {
}

//getListURL dd
func (r *RegistryManager) getListURL() string {
	return utils.GlobalRepo.LoginURL + "/v2/_catalog"
}

//getRequestDeleteHeaderUrl dd
func (r *RegistryManager) getRequestDeleteHeaderURL(imageName string, imageTag string) string {
	newImageName := imageName
	if strings.Contains(imageName, utils.GlobalRepo.URI+"/") {
		newImageName = strings.ReplaceAll(imageName, utils.GlobalRepo.URI+"/", "")
	}
	return utils.GlobalRepo.LoginURL + "/v2/" + newImageName + "/manifests/" + imageTag
}

//getDeleteURL dd
func (r *RegistryManager) getDeleteURL(imageName string, imageDigest string) string {
	newImageName := imageName
	if strings.Contains(imageName, utils.GlobalRepo.URI+"/") {
		newImageName = strings.ReplaceAll(imageName, utils.GlobalRepo.URI+"/", "")
	}
	return utils.GlobalRepo.LoginURL + "/v2/" + newImageName + "/manifests/" + imageDigest
}

//getImageTagListURL dd
func (r *RegistryManager) getImageTagListURL(imageName string) string {
	newImageName := imageName
	if strings.Contains(imageName, utils.GlobalRepo.URI+"/") {
		newImageName = strings.ReplaceAll(imageName, utils.GlobalRepo.URI+"/", "")
	}
	return utils.GlobalRepo.LoginURL + "/v2/" + newImageName + "/tags/list"
}

/*
//getAuthKey dd
func (r *RegistryManager) getAuthKey() string {
	return util.GlobalRepo.RegistryAuthKey
}
*/
