package noderegistry

import (
	"context"
	"fmt"
	"openmcp/openmcp/apis"
	v1alpha1 "openmcp/openmcp/apis/globalcache/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nodeapi "openmcp/openmcp/openmcp-globalcache/pkg/run/dist"
)

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager
	omcplog.V(4).Info("NewController start")
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		omcplog.V(0).Info("getting delegating client for live cluster: ", err)
		return nil, err
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(0).Info("getting delegating client for ghost cluster: ", err)
			return nil, err
		}
		ghostclients = append(ghostclients, ghostclient)
	}
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		omcplog.V(0).Info("adding APIs to live cluster's scheme: ", err)
		return nil, err
	}
	if err := co.WatchResourceReconcileObject(live, &v1alpha1.NodeRegistry{}, controller.WatchOptions{}); err != nil {
		omcplog.V(0).Info("setting up Pod watch in live cluster: ", err)
		return nil, err
	}
	omcplog.V(4).Info("NewController end")
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

// Reconcile reads that state of the cluster for a GlobalRegistry object and makes changes based on the state read
// and what is in the GlobalRegistry.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	omcplog.V(3).Info("Function Called Reconcile")

	instance := &v1alpha1.NodeRegistry{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.V(0).Info("get instance error")
	}
	r.Run(instance)
	if instance.Status.Succeeded == true {
		// 이미 성공한 케이스는 로직을 안탄다.
		omcplog.V(4).Info(instance.Name + " already succeed")
		return reconcile.Result{Requeue: false}, nil
	}
	if instance.Status.Succeeded == false && instance.Status.Reason != "" {
		// 이미 실패한 케이스는 로직을 다시 안탄다.
		omcplog.V(4).Info(instance.Name + " already failed")
		return reconcile.Result{Requeue: false}, nil
	}
	return reconcile.Result{}, nil
}

func (r *reconciler) Run(instance *v1alpha1.NodeRegistry) (bool, error) {
	omcplog.V(3).Info("\n[Command]] :" + instance.Spec.Command)
	var registryManager nodeapi.RegistryManager
	imageList := instance.Spec.ImageLists
	omcplog.V(3).Info("여기", imageList)
	ImageSpec := imageList[0]
	omcplog.V(3).Info("globalcache select image : " + ImageSpec.ImageName)
	// omcplog.V(3).Info("globalcache select image : openmcp/keti-http-generatore")
	// omcplog.V(3).Info("globalcache select image : keti-preprocessor")
	// omcplog.V(3).Info("globalcache select image : openmcp/keti-iotgateway")

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
			err = registryManager.CreatePullJobForCluster(ImageSpec.ImageName, ImageSpec.TagName)
			if err != nil {
				return false, err
			}
		} else {
			err = registryManager.CreatePullJob(instance.Spec.NodeName, ImageSpec.ImageName, ImageSpec.TagName)
			if err != nil {
				return false, err
			}
		}
	case "push":
		err = registryManager.CreatePushJobForCluster(ImageSpec.ImageName, ImageSpec.TagName)
		if err != nil {
			return false, err
		}

	//case "tagList":
	default:
		return false, fmt.Errorf("Command is not valid")
	}

	return true, nil
}
