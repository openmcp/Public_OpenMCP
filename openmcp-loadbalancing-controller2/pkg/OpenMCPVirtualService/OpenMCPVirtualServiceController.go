package OpenMCPVirtualService

import (
	"context"
	"fmt"

	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	"openmcp/openmcp/apis"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"admiralty.io/multicluster-controller/pkg/reference"
	"github.com/getlantern/deepcopy"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("[OpenMCP Deployment] Function Called NewController")
	cm = myClusterManager

	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
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
	zone        string
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
			ObjectMeta: v1.ObjectMeta{
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
		vs, err2 := makeVirtualService(ovs)
		if err2 != nil {
			return reconcile.Result{}, err2
		}
		if err2 == nil {
			// Update VirtualService
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
		}
	}

	return reconcile.Result{}, nil
}

func getLocClusters() []LocCluster {
	locationSlice := []LocCluster{}
	for _, memberCluster := range cm.Cluster_list.Items {
		nodeList := &corev1.NodeList{}
		err := cm.Cluster_genClients[memberCluster.Name].List(context.TODO(), nodeList, "default")
		if err != nil {
			fmt.Println("get NodeList Error")
			continue
		}
		for _, node := range nodeList.Items {
			if _, ok := node.Labels["node-role.kubernetes.io/master"]; ok {

				l := LocCluster{
					clusterName: memberCluster.Name,
					region:      node.Labels["topology.kubernetes.io/region"],
					zone:        node.Labels["topology.kubernetes.io/zone"],
				}
				locationSlice = append(locationSlice, l)

				break
			}

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
func makeVirtualService(ovs *resourcev1alpha1.OpenMCPVirtualService) (*v1alpha3.VirtualService, error) {
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
	vs.Spec.Http = createVsHttps(ovs)

	reference.SetMulticlusterControllerReference(vs, reference.NewMulticlusterOwnerReference(ovs, ovs.GroupVersionKind(), "openmcp"))

	return vs, nil
}
func createVsHttps(ovs *resourcev1alpha1.OpenMCPVirtualService) (vsHttps []*networkingv1alpha3.HTTPRoute) {
	vshttps := []*networkingv1alpha3.HTTPRoute{}

	locClusters := getLocClusters()
	usedZones := []string{}

	for _, http := range ovs.Spec.Http {

		// 디폴트 경로 생성
		exactZone := "default"
		vsHttp, _ := createVsHttp(http, exactZone)
		vshttps = append(vshttps, vsHttp)

		// 지역(클러스터)별 경로 생성
		for _, locCluster := range locClusters {

			if contains(usedZones, locCluster.zone) {
				continue
			}

			exactZone := locCluster.zone
			usedZones = append(usedZones, locCluster.zone)

			vsHttp, _ := createVsHttp(http, exactZone)

			vshttps = append(vshttps, vsHttp)

		}

	}
	return vshttps
}
func createVsHttp(http *networkingv1alpha3.HTTPRoute, exactZone string) (*networkingv1alpha3.HTTPRoute, error) {
	vsHttp := &networkingv1alpha3.HTTPRoute{}
	err := deepcopy.Copy(&vsHttp, &http)
	if err != nil {
		return nil, err
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
		//vsHttp.Match[i].Headers["client-zone"] = stringMatch

	}

	vsHttp.Route = []*networkingv1alpha3.HTTPRouteDestination{}
	locClusters := getLocClusters()

	for _, hr := range http.Route {

		if exactZone == "default " {

		}
		for i, locCluster2 := range locClusters {

			vsRoute := &networkingv1alpha3.HTTPRouteDestination{}
			err = deepcopy.Copy(&vsRoute, &hr)
			if err != nil {
				return nil, err
			}

			// TODO 서비스가 있는지 체크
			// TODO Weight 계산
			vsRoute.Destination.Subset = locCluster2.clusterName
			vsRoute.Weight = 1
			if i == len(locClusters)-1 {
				vsRoute.Weight = 100 - int32(len(locClusters)) + 1
			}

			vsHttp.Route = append(vsHttp.Route, vsRoute)
		}
	}
	return vsHttp, nil
}
