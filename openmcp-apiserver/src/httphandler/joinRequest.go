package httphandler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/omcplog"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Live *cluster.Cluster

func JoinHandler(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	c := &cobrautil.KubeConfig{}
	err = yaml.Unmarshal(fileBytes, c)

	CreateClusterResource(c.Clusters[0].Name, fileBytes)

	a := []byte("OK\n")
	w.Write(a)
}

func JoinCloudHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	_ = r.ParseForm()

	c_name := r.Form.Get("clustername")
	c_type := r.Form.Get("clustertype")
	access_token := r.Form.Get("accesstoken")

	fmt.Println(c_name, " / ", c_type)

	c_file, _, err := r.FormFile("file")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}
	defer c_file.Close()

	fileBytes, err := ioutil.ReadAll(c_file)

	c := &cobrautil.KubeConfig{}
	err = yaml.Unmarshal(fileBytes, c)

	CreateCloudClusterResource(c_name, c_type, access_token, fileBytes)

	a := []byte("OK\n")
	w.Write(a)
}

func CreateClusterResource(name string, config []byte) (string, error) {

	clusterCR := &clusterv1alpha1.OpenMCPCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "OpenMCPCluster",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "openmcp",
		},
		Spec: clusterv1alpha1.OpenMCPClusterSpec{
			ClusterPlatformType: "On-Premise",
			JoinStatus:          "UNJOIN",
			MetalLBRange: clusterv1alpha1.MetalLBRange{
				AddressFrom: "",
				AddressTo:   "",
			},
			KubeconfigInfo: config,
		},
	}

	liveClient, _ := Live.GetDelegatingClient()

	err := liveClient.Create(context.TODO(), clusterCR)

	if err != nil {
		omcplog.V(4).Info("Fail to create openmcpcluster resource")
		fmt.Println(err)
	} else {
		omcplog.V(4).Info("Success to create openmcpcluster resource")
	}

	return clusterCR.Name, err
}

func CreateCloudClusterResource(c_name string, c_type string, access_token string, config []byte) (string, error) {

	clusterCR := &clusterv1alpha1.OpenMCPCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "OpenMCPCluster",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      c_name,
			Namespace: "openmcp",
		},
		Spec: clusterv1alpha1.OpenMCPClusterSpec{
			ClusterPlatformType: c_type,
			JoinStatus:          "UNJOIN",
			KubeconfigInfo:      config,
		},
	}

	liveClient, _ := Live.GetDelegatingClient()

	err := liveClient.Create(context.TODO(), clusterCR)

	if err != nil {
		omcplog.V(4).Info("Fail to create openmcpcluster resource")
		fmt.Println(err)
	} else {
		omcplog.V(4).Info("Success to create openmcpcluster resource")
	}

	return clusterCR.Name, err
}
