package DestinationRule

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-loadbalancing-controller2/pkg/DestinationRuleWeight"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"github.com/gogo/protobuf/types"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var cm *clusterManager.ClusterManager

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(4).Info("Start DestinationRuleController")
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

	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}
	if err := v1alpha3.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &resourcev1alpha1.OpenMCPService{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	return co, nil
}

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(4).Info("DestinationRulesController Reconcile Called")
	omcplog.V(4).Info("[Resource Info] Namespace : ", req.Namespace, ", Name : ", req.Name)

	// OpenMCPService 존재하지 않으면 DR 삭제
	check_osvc := &resourcev1alpha1.OpenMCPService{}
	err_osvc := r.live.Get(context.TODO(), req.NamespacedName, check_osvc)
	if err_osvc != nil && errors.IsNotFound(err_osvc) {
		obj := &v1alpha3.DestinationRule{
			ObjectMeta: v1.ObjectMeta{
				Name:      req.Name,
				Namespace: req.Namespace,
			},
		}
		err_osvc_dr_delete := r.live.Delete(context.TODO(), obj)
		if err_osvc_dr_delete != nil {
			omcplog.V(2).Info(err_osvc_dr_delete)
			//return reconcile.Result{}, err_osvc_dr_delete
		}

		omcplog.V(2).Info(">>> Delete DestinationRules Resource")

		return reconcile.Result{}, nil
	}

	// OpenMCPService 존재 시, DestinationRule 존재 유무 확인
	check_dr := &v1alpha3.DestinationRule{}
	err_dr := r.live.Get(context.TODO(), req.NamespacedName, check_dr)

	// 없으면 Create
	if err_dr == nil || (err_dr != nil && errors.IsNotFound(err_dr)) {
		if err_dr != nil && errors.IsNotFound(err_dr) {

			//원본----------------------------------------------------
			drweight := map[string][]DestinationRuleWeight.DRWeight{}

			drweight = DestinationRuleWeight.DistributeList[req.NamespacedName].DRScore

			distributeTarget := map[string]map[string]uint32{}

			for key, value := range drweight {
				tmp := map[string]uint32{}
				for _, w := range value {
					tmp[w.ToRegionZone] = uint32(w.ConvertToWeight)
				}
				distributeTarget[key] = tmp
			}

			//-------------------------------------------------------

			//임시----------------------------------------------------
			/*
				distributeTarget := map[string]map[string]uint32{}

				tmp := map[string]uint32{}
				tmp["kr/seoul/*"] = 90
				tmp["kr/gyeonggi-do/*"] = 10
				distributeTarget["kr/seoul/*"] = tmp

				tmp2 := map[string]uint32{}
				tmp2["kr/seoul/*"] = 10
				tmp2["kr/gyeonggi-do/*"] = 90
				distributeTarget["kr/gyeonggi-do/*"] = tmp2

				fmt.Println(distributeTarget)
			*/
			//----------------------------------------------------

			distribute := []*networkingv1alpha3.LocalityLoadBalancerSetting_Distribute{}

			for k, v := range distributeTarget {
				tmp_dis := &networkingv1alpha3.LocalityLoadBalancerSetting_Distribute{
					From: k,
					To:   v,
				}
				distribute = append(distribute, tmp_dis)
			}

			dr := &v1alpha3.DestinationRule{}
			dr.Name = req.Name
			dr.Namespace = req.Namespace
			dr.Spec.Host = req.Name + "." + req.Namespace + ".svc.cluster.local"

			tp := &networkingv1alpha3.TrafficPolicy{
				LoadBalancer: &networkingv1alpha3.LoadBalancerSettings{
					LbPolicy: nil,
					LocalityLbSetting: &networkingv1alpha3.LocalityLoadBalancerSetting{
						Distribute: distribute,
						Failover:   nil,
						Enabled: &types.BoolValue{
							Value:                true,
							XXX_NoUnkeyedLiteral: struct{}{},
							XXX_unrecognized:     nil,
							XXX_sizecache:        0,
						},
						XXX_NoUnkeyedLiteral: struct{}{},
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
				OutlierDetection: &networkingv1alpha3.OutlierDetection{
					ConsecutiveErrors:              0,
					SplitExternalLocalOriginErrors: false,
					ConsecutiveLocalOriginFailures: nil,
					ConsecutiveGatewayErrors:       nil,
					Consecutive_5XxErrors: &types.UInt32Value{
						Value:                100,
						XXX_NoUnkeyedLiteral: struct{}{},
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
					Interval: &types.Duration{
						Seconds:              1,
						Nanos:                0,
						XXX_NoUnkeyedLiteral: struct{}{},
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
					BaseEjectionTime: &types.Duration{
						Seconds:              60,
						Nanos:                0,
						XXX_NoUnkeyedLiteral: struct{}{},
						XXX_unrecognized:     nil,
						XXX_sizecache:        0,
					},
					MaxEjectionPercent:   0,
					MinHealthPercent:     0,
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				},
			}

			dr.Spec.TrafficPolicy = tp

			subsets := []*networkingv1alpha3.Subset{}

			for _, memberCluster := range cm.Cluster_list.Items {
				clusterSubset := make(map[string]string)
				clusterSubset["cluster"] = memberCluster.Name

				tmp := &networkingv1alpha3.Subset{
					Name:   memberCluster.Name,
					Labels: clusterSubset,
				}

				subsets = append(subsets, tmp)
			}

			dr.Spec.Subsets = subsets

			err_dr_create := r.live.Create(context.TODO(), dr)

			if err_dr_create != nil {
				omcplog.V(2).Info(err_dr_create)
				//return reconcile.Result{}, err_dr_create
			}

			omcplog.V(2).Info(">>> Create DestinationRules Resource")

			return reconcile.Result{}, nil
		}

	} else {
		omcplog.V(2).Info(err_dr)
	}

	return reconcile.Result{}, nil
}
