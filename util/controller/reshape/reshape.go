package reshape

import (
	"context"
	"fmt"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"os"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"github.com/jinzhu/copier"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fedapis "sigs.k8s.io/kubefed/pkg/apis"
	"sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var c chan string

var prev_length int
var clusterStatus map[string]bool
var prev_clusterStatus map[string]bool

var cm *clusterManager.ClusterManager

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	omcplog.V(2).Info("Start Reshape Controller")
	c = make(chan string)
	cm = myClusterManager

	clusterStatus = make(map[string]bool)
	prev_clusterStatus = make(map[string]bool)

	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	/*ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}*/

	r := &reconciler{live: liveclient} //, ghosts: ghostclients, ghostNamespace: ghostNamespace}

	co := controller.New(r, controller.Options{})
	if err := fedapis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	if err := co.WatchResourceReconcileObject(context.TODO(), live, &v1beta1.KubeFedCluster{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	//r.initGlobal()

	return co, nil
}
func (r *reconciler) initGlobal() {

	// Fetch the instance
	kubeFedClusterList := &v1beta1.KubeFedClusterList{}
	err := r.live.List(context.TODO(), kubeFedClusterList, &client.ListOptions{})
	if err != nil {
		klog.V(0).Info(err)
	} else {
		for _, cluster := range kubeFedClusterList.Items {
			prev_clusterStatus[cluster.Name] = true
			clusterStatus[cluster.Name] = true
			for _, cond := range cluster.Status.Conditions {
				if cond.Type == "Offline" {
					prev_clusterStatus[cluster.Name] = false
					clusterStatus[cluster.Name] = false
					break
				}
			}

		}
	}

}

type reconciler struct {
	live client.Client
	//ghosts         []client.Client
	//ghostNamespace string
}

//var i int = 0

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	//i += 1

	// Fetch the instance
	kubeFedClusterList := &v1beta1.KubeFedClusterList{}
	err := r.live.List(context.TODO(), kubeFedClusterList, &client.ListOptions{})
	if err != nil {
		klog.V(0).Info(err)
	}

	//fmt.Println("Reshape Cluster [len(kubeFedClusterList.Items)] -> ", len(kubeFedClusterList.Items))
	//fmt.Println("Reshape Cluster [prev_length] -> ", prev_length)

	ReshapeFlag := false
	if len(kubeFedClusterList.Items) != prev_length {
		ReshapeFlag = true
	}

	for _, cluster := range kubeFedClusterList.Items {
		for _, cond := range cluster.Status.Conditions {
			clusterStatus[cluster.Name] = true
			if cond.Type == "Offline" {
				clusterStatus[cluster.Name] = false
				break
			}
		}
		if prev_clusterStatus[cluster.Name] != clusterStatus[cluster.Name] {
			prev_clusterStatus[cluster.Name] = clusterStatus[cluster.Name]
			ReshapeFlag = true

		}

	}

	if ReshapeFlag {

		omcplog.V(2).Info("Reshape Cluster ...")
		//i = 0

		newCm := clusterManager.NewClusterManager()
		//cm.Mutex.Lock()
		copier.Copy(cm, newCm)
		//cm.Mutex.Unlock()

		prev_length = len(kubeFedClusterList.Items)
		//c <- "reshape"
	}

	return reconcile.Result{}, nil // err
}

func SetupSignalHandler() (stopCh <-chan struct{}) {

	stop := make(chan struct{})

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
