package httphandler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/omcplog"
	"strings"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Live *cluster.Cluster

func JoinHandler(w http.ResponseWriter, r *http.Request) {
	nodeinfo := r.FormValue("data")

	file, _, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	c := &cobrautil.KubeConfig{}
	serverIP := c.Clusters[0].Cluster.Server[9:]
	serverIP = strings.Split(serverIP, ":")[0]
	err = yaml.Unmarshal(fileBytes, c)

	CreateClusterResource(c.Clusters[0].Name, serverIP, fileBytes, nodeinfo)

	a := []byte("OK\n")
	w.Write(a)
}

func JoinCloudHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	_ = r.ParseForm()

	c_name := r.Form.Get("clustername")
	c_type := r.Form.Get("clustertype")
	c_loc := r.Form.Get("clusterlocation")

	fmt.Println(c_name, " / ", c_type, " / ", c_loc)

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

	CreateCloudClusterResource(c_name, c_type, fileBytes, c_loc)

	a := []byte("OK\n")
	w.Write(a)
}

func CreateClusterResource(name string, serverIP string, config []byte, nodeinfo string) (string, error) {
	regionmap := map[string]string{}
	zonemap := map[string]string{}

	slice1 := strings.Split(nodeinfo, "#")

	for _, str := range slice1 {
		if len(str) > 0 {
			slice2 := strings.Split(str, ",")

			if len(slice2) == 3 {
				regionmap[slice2[1]] = ""
				zonemap[slice2[2]] = ""
			}
		}
	}

	regionstr := ""
	zonestr := ""

	for key, _ := range regionmap {
		if regionstr != "" {
			regionstr = regionstr + "," + key
		} else {
			regionstr = key
		}
	}

	for key, _ := range zonemap {
		if zonestr != "" {
			zonestr = zonestr + "," + key
		} else {
			zonestr = key
		}
	}

	ni := clusterv1alpha1.NodeInfo{
		Region: regionstr,
		Zone:   zonestr,
	}

	omcplog.V(4).Info("'", name, "' nodeInfo ", ni.Region, " / ", ni.Zone)

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
			ServerIP:               serverIP,
			ClusterPlatformType:    "On-Premise",
			ClusterNetworkLocation: "",
			JoinStatus:             "UNJOIN",
			NodeInfo:               ni,
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
		omcplog.V(4).Info("Fail to create openmcpcluster resource'", name, "'")
		fmt.Println(err)
	} else {
		omcplog.V(4).Info("Success to create openmcpcluster resource '", name, "'")
	}

	return clusterCR.Name, err
}

func CreateCloudClusterResource(c_name string, c_type string, config []byte, c_loc string) (string, error) {

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
			ClusterPlatformType:    c_type,
			ClusterNetworkLocation: "",
			JoinStatus:             "UNJOIN",
			KubeconfigInfo:         config,
		},
	}

	//[Google] GKE Region/Zone Setting
	if c_type == "GKE" && c_loc == "asia-northeast1-a" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "JP",
			Zone:   "Tokyo",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "GKE" && c_loc == "europe-west3-a" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "DE",
			Zone:   "Hesse",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "GKE" && c_loc == "australia-southeast1-a" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "AU",
			Zone:   "New-South-Wales",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	//[Amazon] EKS Region/Zone Setting
	if c_type == "EKS" && c_loc == "eu-west-2" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "GB",
			Zone:   "England",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "EKS" && c_loc == "eu-west-3" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "FR",
			Zone:   "Île-de-France",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "EKS" && c_loc == "us-east-1" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "US",
			Zone:   "Virginia",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	//[Azure] AKS Region/Zone Setting
	if c_type == "AKS" && c_loc == "norwayeast" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "NO",
			Zone:   "Oslo",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "AKS" && c_loc == "southeastasia" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "TH",
			Zone:   "Bangkok",
		}

		clusterCR.Spec.NodeInfo = ni
	}

	if c_type == "AKS" && c_loc == "westus" {
		ni := clusterv1alpha1.NodeInfo{
			Region: "US",
			Zone:   "California",
		}

		clusterCR.Spec.NodeInfo = ni
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
