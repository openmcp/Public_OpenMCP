package resource

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"

	"k8s.io/client-go/tools/clientcmd"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"strings"
)

func ClusterInfo(kfc *v1beta1.KubeFedCluster) []string{

	ClusterNamespace := kfc.Namespace
	ClusterName := kfc.Name
	ClusterStatus := string(kfc.Status.Conditions[0].Status)

	ClusterRegion := *kfc.Status.Region
	ClusterZone := strings.Join(kfc.Status.Zones, ",")
	ClusterAPIEndpoint := kfc.Spec.APIEndpoint

	//ClusterClusterName := kubecontext.Context.Cluster
	age := cobrautil.GetAge(kfc.CreationTimestamp.Time)

	data := []string{ClusterNamespace, ClusterName, ClusterStatus, ClusterRegion, ClusterZone, ClusterAPIEndpoint, age}

	return data
}
func PrintKubeFedCluster(body []byte) {
	no := v1beta1.KubeFedCluster{}
	err := yaml.Unmarshal(body, &no)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}


	data := ClusterInfo(&no)
	datas = append(datas, data)

	DrawClusterTable(datas)

}
func PrintKubeFedClusterList(body []byte) {
	resourceStruct := v1beta1.KubeFedClusterList{}
	err := yaml.Unmarshal(body, &resourceStruct)
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}
	datas := [][]string{}

	for _, no := range resourceStruct.Items {
		data := ClusterInfo(&no)
		datas = append(datas, data)
	}

	if len(resourceStruct.Items) == 0 {
		ns := "default"
		if cobrautil.Option_namespace != "" {
			ns = cobrautil.Option_namespace
		}
		fmt.Println("No resources found in "+ ns +" Cluster.")
		return
	}

	DrawClusterTable(datas)

}

func DrawClusterTable(datas [][]string){
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"NS", "ClusterName", "Status", "Region", "Zones", "APIEndpoint", "AGE"})
	table.SetBorder(false)
	table.AppendBulk(datas)
	table.Render()
}


//func GetCluster(clusterName string){
//	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
//	genClient := genericclient.NewForConfigOrDie(kubeconfig)
//
//	kubeFedCluster := &v1beta1.KubeFedCluster{}
//	genClient.Get(context.Background(), kubeFedCluster,"kube-federation-system", clusterName)
//
//	if cobrautil.Option_filetype == ""{
//		datas := [][]string{}
//		data := []string{kubeFedCluster.Name, string(kubeFedCluster.Status.Conditions[0].Status), *kubeFedCluster.Status.Region, strings.Join(kubeFedCluster.Status.Zones, ","), kubeFedCluster.Spec.APIEndpoint, kubeFedCluster.GenerateName }
//
//		datas = append(datas, data)
//		table := tablewriter.NewWriter(os.Stdout)
//		table.SetHeader([]string{"ClusterName", "Status", "Region", "Zones", "apiEndpoint", "AGE"})
//		table.SetBorder(false)
//		table.AppendBulk(datas)
//		table.Render()
//
//	} else if cobrautil.Option_filetype == "yaml"{
//		res, _ := util.CmdExec("kubectl get kubefedclusters " + clusterName + " -n kube-federation-system -o yaml")
//		fmt.Println(res)
//	}
//}
//
func GetClusterList(){
	kubeconfig, _ := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	genClient := genericclient.NewForConfigOrDie(kubeconfig)

	clusterList := clusterManager.ListKubeFedClusters(genClient, "kube-federation-system")

	if cobrautil.Option_filetype == ""{
		datas := [][]string{}
		for _, Cluster := range clusterList.Items {
			data := []string{Cluster.Name, string(Cluster.Status.Conditions[0].Status), *Cluster.Status.Region, strings.Join(Cluster.Status.Zones, ","), Cluster.Spec.APIEndpoint, Cluster.GenerateName }
			datas = append(datas, data)

		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ClusterName", "Status", "Region", "Zones", "apiEndpoint", "AGE"})
		table.SetBorder(false)
		table.AppendBulk(datas)
		table.Render()

	} else if cobrautil.Option_filetype == "yaml"{
		res, _ := util.CmdExec("kubectl get kubefedclusters -n kube-federation-system -o yaml")
		fmt.Println(res)
	}
}

