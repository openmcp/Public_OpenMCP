package utils

import (
	"encoding/base64"
	"encoding/json"

	"github.com/storageos/go-api/types"
)

//GetClientset kubernetes의 clientset 생성
// func GetClientset(clusterInfo string) (*kubernetes.Clientset, error) {
// 	var clientset *kubernetes.Clientset
// 	con, err := clientcmd.NewClientConfigFromBytes([]byte(clusterInfo))
// 	if err != nil {
// 		return nil, err
// 	}
// 	clientconf, err := con.ClientConfig()
// 	if err != nil {
// 		return nil, err
// 	}
// 	clientset, err = kubernetes.NewForConfig(clientconf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return clientset, nil
// }

// MakeDockerAuth : Docker Auth 를 생성하는 함수.
func MakeDockerAuth() string {
	authConfig := types.AuthConfig{
		Username:      GlobalRepo.Username,
		Password:      GlobalRepo.Password,
		ServerAddress: GlobalRepo.LoginURL,
	}

	buf, _ := json.Marshal(authConfig)
	regauth := base64.URLEncoding.EncodeToString(buf)
	return regauth
}

func SetGlobalRegistryCommand(command string) string {

	retCommand := "echo `hostname`; "
	retCommand += "mkdir -p /etc/docker/certs.d/" + DockerHubRepo.URI + "; "
	retCommand += "echo -e \"" + DockerHubRepo.Cert + "\" > /etc/docker/certs.d/" + DockerHubRepo.URI + "/server.crt; "
	retCommand += "echo " + DockerHubRepo.Password + " | docker login -u " + DockerHubRepo.Username + " --password-stdin " + DockerHubRepo.URI + ";"
	retCommand += command
	return retCommand
}
func SetDockerHublRegistryCommand(command string) string {

	retCommand := "echo `hostname`; "
	retCommand += "mkdir -p /etc/docker/certs.d/" + GlobalRepo.URI + "; "
	retCommand += "echo -e \"" + GlobalRepo.Cert + "\" > /etc/docker/certs.d/" + GlobalRepo.URI + "/server.crt; "
	retCommand += "echo " + GlobalRepo.Password + " | docker login -u " + GlobalRepo.Username + " --password-stdin " + GlobalRepo.URI + ";"
	retCommand += command
	return retCommand
}
