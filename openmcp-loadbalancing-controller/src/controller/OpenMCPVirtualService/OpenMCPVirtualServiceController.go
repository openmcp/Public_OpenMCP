package OpenMCPVirtualService

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	"openmcp/openmcp/apis"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-analytic-engine/src/protobuf"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	"github.com/getlantern/deepcopy"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP VirtualS] Function Called NewController")
	cm = myClusterManager

	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(4).Info("Error getting delegating client for ghost cluster [", ghost.Name, "]")
			//return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		} else {
			ghostclients = append(ghostclients, ghostclient)
		}
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})

	// if err := v1alpha3.AddToScheme(live.GetScheme()); err != nil {
	// 	return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	// }

	// if err := co.WatchResourceReconcileObject(context.TODO(), live, &v1alpha3.VirtualService{}, controller.WatchOptions{}); err != nil {
	// 	return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	// }
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	if err := v1alpha3.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPVirtualService{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	// Note: At the moment, all clusters share the same scheme under the hood
	// (k8s.io/client-go/kubernetes/scheme.Scheme), yet multicluster-controller gives each cluster a scheme pointer.
	// Therefore, if we needed a custom resource in multiple clusters, we would redundantly
	// add it to each cluster's scheme, which points to the same underlying scheme.

	//for _, ghost := range ghosts {
	//	if err := co.WatchResourceReconcileController(ghost, &appsv1.Deployment{}, controller.WatchOptions{}); err != nil {
	//		return nil, fmt.Errorf("setting up PodGhost watch in ghost cluster: %v", err)
	//	}
	//}

	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var i int = 0
var syncIndex int = 0

type LocCluster struct {
	clusterName string
	region      string
	zones       []string
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {

	omcplog.V(4).Info("[OpenMCP VirtualService] Function Called Reconcile", req.Name, ", ", req.Namespace)

	i += 1

	// Fetch the OpenMCPDeployment instance
	ovs := &resourcev1alpha1.OpenMCPVirtualService{}
	err := r.live.Get(context.TODO(), req.NamespacedName, ovs)
	if err != nil && errors.IsNotFound(err) {

		omcplog.V(2).Info("[Delete Detect]")
		omcplog.V(2).Info("Delete VirtualService")
		obj := &v1alpha3.VirtualService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
		}
		err2 := r.live.Delete(context.TODO(), obj)
		if err != nil {
			return reconcile.Result{}, err2
		}

		return reconcile.Result{}, nil
	}
	omcplog.V(5).Info("Resource Get => [Name] : " + ovs.Name + " [Namespace]  : " + ovs.Namespace)

	checkVs := &v1alpha3.VirtualService{}
	err = r.live.Get(context.TODO(), req.NamespacedName, checkVs)
	if err == nil || (err != nil && errors.IsNotFound(err)) {
		vs, err2 := MakeVirtualService(ovs)
		if err == nil {
			// Update VirtualService
			vs.ResourceVersion = checkVs.ResourceVersion
			err3 := r.live.Update(context.TODO(), vs)
			if err3 != nil {
				return reconcile.Result{}, err3
			}
		} else if err != nil && errors.IsNotFound(err) {
			// Create VirtualService
			err3 := r.live.Create(context.TODO(), vs)
			if err3 != nil {
				return reconcile.Result{}, err3
			}
		} else {
			return reconcile.Result{}, err2
		}
	}

	return reconcile.Result{}, nil
}

func getLocClusters() []LocCluster {
	// locationSlice := []LocCluster{}

	// region := ""
	// zones := []string{}
	// for _, memberCluster := range cm.Cluster_list.Items {
	// 	nodeList := &corev1.NodeList{}
	// 	err := cm.Cluster_genClients[memberCluster.Name].List(context.TODO(), nodeList, "default")
	// 	if err != nil {
	// 		fmt.Println("get NodeList Error")
	// 		continue
	// 	}
	// 	for _, node := range nodeList.Items {
	// 		if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
	// 			region = node.Labels["topology.kubernetes.io/region"]
	// 		}
	// 		zones = append(zones, node.Labels["topology.kubernetes.io/zone"])
	// 	}
	// 	l := LocCluster{
	// 		clusterName: memberCluster.Name,
	// 		region:      region,
	// 		zones:       zones,
	// 	}
	// 	locationSlice = append(locationSlice, l)

	// }
	// return locationSlice

	locationSlice := []LocCluster{}

	// region := ""
	// zones := []string{}

	//ocList := &clusterv1alpha1.OpenMCPClusterList{}
	//err := cm.Host_client.List(context.TODO(), ocList, "openmcp")
	ocList, err := cm.Crd_client.OpenMCPCluster("openmcp").List(v1.ListOptions{})

	if err != nil {
		fmt.Println("OpenMCPClusterList err : ", err)
	}

	for _, oc := range ocList.Items {
		if oc.Spec.JoinStatus == "JOIN" {
			region := oc.Spec.NodeInfo.Region
			zones := strings.Split(oc.Spec.NodeInfo.Zone, ",")

			l := LocCluster{
				clusterName: oc.Name,
				region:      region,
				zones:       zones,
			}
			locationSlice = append(locationSlice, l)
		}
	}
	return locationSlice

}
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func MakeVirtualService(ovs *resourcev1alpha1.OpenMCPVirtualService) (*v1alpha3.VirtualService, error) {
	omcplog.V(4).Info("func MakeVirtualService Called")
	vs := &v1alpha3.VirtualService{}
	vs.Name = ovs.Name
	vs.Namespace = ovs.Namespace
	err := deepcopy.Copy(&vs.Labels, &ovs.Labels)
	if err != nil {
		return nil, err
	}
	err = deepcopy.Copy(&vs.Spec.Hosts, &ovs.Spec.Hosts)
	if err != nil {
		return nil, err
	}
	err = deepcopy.Copy(&vs.Spec.Gateways, &ovs.Spec.Gateways)
	if err != nil {
		return nil, err
	}

	vs.Spec.Http, err = createVsHttps(ovs)
	if err != nil {
		return nil, err
	}

	reference.SetMulticlusterControllerReference(vs, reference.NewMulticlusterOwnerReference(ovs, ovs.GroupVersionKind(), "openmcp"))

	omcplog.V(4).Info("func MakeVirtualService Ended")
	return vs, nil
}

type RegionZone struct {
	region string
	zone   string
}

func createVsHttps(ovs *resourcev1alpha1.OpenMCPVirtualService) ([]*networkingv1alpha3.HTTPRoute, error) {
	omcplog.V(4).Info("func createVsHttps Called")
	vsHttps := []*networkingv1alpha3.HTTPRoute{}

	locClusters := getLocClusters()
	usedRegionZones := []RegionZone{}

	for _, ovsHttp := range ovs.Spec.Http {

		// 디폴트 경로 생성
		exactRegion := "default"
		exactZone := "default"
		vsHttp, err := createVsHttp(ovsHttp, exactRegion, exactZone, ovs.Namespace)
		if err != nil {
			return nil, err
		}
		vsHttps = append(vsHttps, vsHttp)

		// 지역(클러스터)별 경로 생성
		for _, locCluster := range locClusters {
			for _, zone := range locCluster.zones {

				skipFlag := false
				for _, usedRegionZone := range usedRegionZones {
					if usedRegionZone.region == locCluster.region && usedRegionZone.zone == zone {
						skipFlag = true
						break
					}
				}
				if skipFlag {
					continue
				}

				usedRegionZones = append(usedRegionZones, RegionZone{locCluster.region, zone})

				exactRegion = locCluster.region
				exactZone = zone

				vsHttp, err := createVsHttp(ovsHttp, exactRegion, exactZone, ovs.Namespace)
				if err != nil {
					return nil, err
				}

				vsHttps = append(vsHttps, vsHttp)

			}

		}

	}
	omcplog.V(4).Info("func createVsHttps Ended")

	return vsHttps, nil
}
func createVsHttp(ovsHttp *networkingv1alpha3.HTTPRoute, exactRegion, exactZone, ns string) (*networkingv1alpha3.HTTPRoute, error) {
	omcplog.V(4).Info("func createVsHttp Called")
	vsHttp := &networkingv1alpha3.HTTPRoute{}
	err := deepcopy.Copy(&vsHttp, &ovsHttp)

	if err != nil {
		return nil, err
	}

	for _, match := range vsHttp.Match {
		stringMatch := &networkingv1alpha3.StringMatch{
			MatchType: &networkingv1alpha3.StringMatch_Exact{
				Exact: exactRegion,
			},
		}
		if match.Headers == nil {
			headers := make(map[string]*networkingv1alpha3.StringMatch)
			match.Headers = headers
		}

		match.Headers["client-region"] = stringMatch
		//vsHttp.Match[i].Headers["client-zone"] = stringMatch

	}
	for _, match := range vsHttp.Match {
		stringMatch := &networkingv1alpha3.StringMatch{
			MatchType: &networkingv1alpha3.StringMatch_Exact{
				Exact: exactZone,
			},
		}
		if match.Headers == nil {
			headers := make(map[string]*networkingv1alpha3.StringMatch)
			match.Headers = headers
		}

		match.Headers["client-zone"] = stringMatch

	}

	vsHttp.Route, _ = createVsHttpRoutes(ovsHttp.Route, exactRegion, exactZone, ns)

	omcplog.V(4).Info("func createVsHttp Ended")
	return vsHttp, nil
}
func createVsHttpRoutes(ovsHttpRoutes []*networkingv1alpha3.HTTPRouteDestination, exactRegion, exactZone, ns string) ([]*networkingv1alpha3.HTTPRouteDestination, error) {
	omcplog.V(4).Info("func createVsHttpRoutes Called")
	vsHttpRoutes := []*networkingv1alpha3.HTTPRouteDestination{}
	locClusters := getLocClusters()

	for _, ovsHttpRoute := range ovsHttpRoutes {

		if exactRegion == "default" {
			// default 경로생성
			vsHttpRoute, err := createVsHttpRoute(ovsHttpRoute, -1, nil, ns)
			if err != nil {
				return nil, err
			}
			if vsHttpRoute == nil {
				// 해당 Cluster에 해당 서비스가 존재하지 않는경우 생성하지 않음
				continue
			}
			vsHttpRoutes = append(vsHttpRoutes, vsHttpRoute)
		} else {
			for i, locCluster := range locClusters {
				vsHttpRoute, err := createVsHttpRoute(ovsHttpRoute, i, &locCluster, ns)
				if err != nil {
					return nil, err
				}
				if vsHttpRoute == nil {
					// 해당 Cluster에 해당 서비스가 존재하지 않는경우 생성하지 않음
					continue
				}
				vsHttpRoutes = append(vsHttpRoutes, vsHttpRoute)
			}
			// default 경로생성
			if len(vsHttpRoutes) == 0 {
				vsHttpRoute, err := createVsHttpRoute(ovsHttpRoute, -1, nil, ns)
				if err != nil {
					return nil, err
				}
				vsHttpRoutes = append(vsHttpRoutes, vsHttpRoute)
			}
		}

	}

	setWeight(vsHttpRoutes, exactRegion, exactZone, ns)
	omcplog.V(4).Info("func createVsHttpRoutes Ended")

	return vsHttpRoutes, nil
}
func createVsHttpRoute(ovsHttpRoute *networkingv1alpha3.HTTPRouteDestination, i int, locCluster *LocCluster, ns string) (*networkingv1alpha3.HTTPRouteDestination, error) {
	omcplog.V(4).Info("func createVsHttpRoute Called")
	// 서비스가 있는지 체크
	svcDomain := ovsHttpRoute.Destination.Host
	svcDomainSplit := strings.Split(svcDomain, ".")

	svcName := svcDomainSplit[0]
	svcNS := ns
	if len(svcDomainSplit) >= 2 {
		svcNS = svcDomainSplit[1]
	}

	omcplog.V(4).Info("locCluster: ", locCluster)

	if locCluster != nil {
		svc := &corev1.Service{}

		if cm.Cluster_genClients[locCluster.clusterName] == nil {
			return nil, nil
		}
		err := cm.Cluster_genClients[locCluster.clusterName].Get(context.TODO(), svc, svcNS, svcName)

		if err != nil && errors.IsNotFound(err) {
			omcplog.V(4).Info(err)
			return nil, nil
		}
		if err != nil {
			omcplog.V(0).Info(err)
			return nil, err
		}
		listOption := &client.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				svc.Spec.Selector,
			),
		}
		omcplog.V(2).Info("Find Service '", svcName, "'in '", locCluster.clusterName, "'")
		podList := &corev1.PodList{}
		err = cm.Cluster_genClients[locCluster.clusterName].List(context.TODO(), podList, svcNS, listOption)
		if len(podList.Items) == 0 {
			omcplog.V(4).Info("!!!!", locCluster.clusterName, " is Not Exist Pod about Svc: '", svcName, "(", svcNS, ")', LabelSelector: '", svc.Spec.Selector, "'")
			return nil, nil
		}
	}

	vsHttpRoute := &networkingv1alpha3.HTTPRouteDestination{}
	err := deepcopy.Copy(&vsHttpRoute, &ovsHttpRoute)
	if err != nil {
		return nil, err
	}
	if locCluster != nil {
		vsHttpRoute.Destination.Subset = locCluster.clusterName
	}
	omcplog.V(4).Info("func createVsHttpRoute Ended")

	return vsHttpRoute, nil
}

type ClusterPrimeNumber struct {
	clusterName    string
	primeNumber    float64
	allocateWeight int
}

var SERVER_IP = os.Getenv("GRPC_SERVER")
var SERVER_PORT = os.Getenv("GRPC_PORT")
var grpcClient = protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

func setWeight(vsHttpRoutes []*networkingv1alpha3.HTTPRouteDestination, fromRegion, fromZone, ns string) {
	omcplog.V(4).Info("func setWeight Called")
	clusterWeight := make(map[string]float64)
	var clusterWeightSum float64 = 0

	for _, vsHttpRoute := range vsHttpRoutes {

		svcDomain := vsHttpRoute.Destination.Host
		svcDomainSplit := strings.Split(svcDomain, ".")

		svcName := svcDomainSplit[0]
		svcNS := ns

		if len(svcDomainSplit) >= 2 {
			svcNS = svcDomainSplit[1]
		}

		clusterName := vsHttpRoute.Destination.Subset
		if clusterName != "" && cm.Cluster_genClients[clusterName] != nil {

			nodeList := &corev1.NodeList{}

			err := cm.Cluster_genClients[clusterName].List(context.TODO(), nodeList, "default")
			if err != nil && errors.IsNotFound(err) {
				return
			}

			svc := &corev1.Service{}
			err = cm.Cluster_genClients[clusterName].Get(context.TODO(), svc, svcNS, svcName)
			if err != nil && errors.IsNotFound(err) {
				fmt.Println("!!!! [Score Calcuation] Cluster: ", clusterName, " is Not Exist Svc : '", svcName, "(", svcNS, ")'")
				break
			}

			podList := &corev1.PodList{}
			listOption := &client.ListOptions{
				LabelSelector: labels.SelectorFromSet(
					svc.Spec.Selector,
				),
			}
			err = cm.Cluster_genClients[clusterName].List(context.TODO(), podList, svcNS, listOption)
			if err != nil && errors.IsNotFound(err) {
				fmt.Println("!!!! [Score Calcuation] Cluster: ", clusterName, " is Not Exist Pod about Svc: '", svcName, "(", svcNS, ")', LabelSelector: '", svc.Spec.Selector, "'")
				break
			}

			var cluster_X_TotalWeight float64 = 0

			//podNodeMatchFind := false
			for _, pod := range podList.Items {
			Loop1:
				for _, node := range nodeList.Items {
					// if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {
					// 	continue
					// }
					if pod.Spec.NodeName == node.Name {
						//podNodeMatchFind = true
						ocList, err := cm.Crd_client.OpenMCPCluster("openmcp").List(v1.ListOptions{})

						if err != nil {
							omcplog.V(0).Info(err)
						}

						for _, oc := range ocList.Items {
							if oc.Name == clusterName {
								toRegion := oc.Spec.NodeInfo.Region
								toZone := oc.Spec.NodeInfo.Zone

								//toRegion := node.Labels["topology.kubernetes.io/region"]
								//toZone := node.Labels["topology.kubernetes.io/zone"]

								regionZoneInfo := protobuf.RegionZoneInfo{
									FromRegion:    fromRegion,
									FromZone:      fromZone,
									ToRegion:      toRegion,
									ToZone:        toZone,
									ToClusterName: clusterName,
									ToNamespace:   svcNS,
									ToPodName:     pod.Name,
								}
								grpcResponse, gRPC_err := grpcClient.SendRegionZoneInfo(context.TODO(), &regionZoneInfo)
								if gRPC_err != nil {
									omcplog.V(0).Info(gRPC_err)
									continue
								}

								cluster_X_TotalWeight = cluster_X_TotalWeight + float64(grpcResponse.Weight)

								fmt.Println("*** [Score Calcuation]", fromRegion, "/", fromZone, " -> ", toRegion, "/", toZone, "(", clusterName, "/", svcNS, "/", pod.Name, "): ", grpcResponse.Weight)
								break Loop1
							}
						}
					}
				}
				// if podNodeMatchFind {
				// 	break
				// }

			}
			if cluster_X_TotalWeight != 0 {
				cluster_X_AVG := cluster_X_TotalWeight / float64(len(podList.Items))

				clusterWeight[clusterName] = cluster_X_AVG
				clusterWeightSum += cluster_X_AVG

				fmt.Println("*** ==> Cluster Score AVG:", cluster_X_AVG)
				//fmt.Println("----------------------------------")
			}

		}
	}
	var totalWeight int32 = 0
	var orgClusterPrimeNumbers []ClusterPrimeNumber

	omcplog.V(4).Info("Start Weight Calculation")
	for _, vsHttpRoute := range vsHttpRoutes {

		clusterName := vsHttpRoute.Destination.Subset

		if clusterName == "" {
			return
		}
		omcplog.V(4).Info("clusterName: ", clusterName)
		orgWeight := clusterWeight[clusterName] / clusterWeightSum * 100
		omcplog.V(4).Info("clusterWeight[clusterName]: ", clusterWeight[clusterName])
		omcplog.V(4).Info("clusterWeightSum: ", clusterWeightSum)
		omcplog.V(4).Info("orgWeight: ", orgWeight)
		var primeNumber float64 = orgWeight - float64(int(orgWeight))

		orgClusterPrimeNumbers = append(orgClusterPrimeNumbers, ClusterPrimeNumber{clusterName, primeNumber, 0})

		weight := int32(math.Floor(orgWeight))
		vsHttpRoute.Weight = weight
		omcplog.V(0).Info("vsHttpRoute.Weight : ", weight)

		totalWeight += weight

	}
	// weight 스케일링된 값을 100으로 맞춰주는 알고리즘
	if totalWeight != 100 {
		vsHttpRoutes[0].Weight += 100 - totalWeight
	}

	// omcplog.V(4).Info("totalWeight : ", totalWeight)

	// // weight 스케일링된 값을 100으로 맞춰주는 알고리즘
	// sort.Slice(orgClusterPrimeNumbers, func(i, j int) bool {
	// 	return orgClusterPrimeNumbers[i].primeNumber > orgClusterPrimeNumbers[j].primeNumber
	// })

	// if totalWeight != 100 {
	// 	restWeight := 100 - totalWeight

	// 	for restWeight > 0 {

	// 		for i, _ := range orgClusterPrimeNumbers {
	// 			orgClusterPrimeNumbers[i].allocateWeight += 1
	// 			restWeight -= 1

	// 			if restWeight == 0 {
	// 				break
	// 			}
	// 		}

	// 	}

	// 	for i, vsHttpRoute := range vsHttpRoutes {
	// 		for j, cpn := range orgClusterPrimeNumbers {
	// 			if vsHttpRoute.Destination.Subset == cpn.clusterName {
	// 				vsHttpRoutes[i].Weight = vsHttpRoutes[i].Weight + (int32(cpn.allocateWeight))
	// 				orgClusterPrimeNumbers[j].allocateWeight = 0
	// 				break
	// 			}
	// 		}

	// 	}

	// }

}
