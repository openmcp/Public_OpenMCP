package cache

import (
	"context"
	"encoding/json"
	"openmcp/openmcp/apis"
	v1alpha1 "openmcp/openmcp/apis/cache/v1alpha1"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"sort"
	"strconv"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	nodeapi "openmcp/openmcp/openmcp-cache/pkg/run/dist"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	//imagelist := make(map[string]int)
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			omcplog.V(0).Info("getting delegating client for ghost cluster: ", err)
			return nil, err
		}

		// listClient := *cm.Cluster_kubeClients[ghost.Name]
		// pods, _ := listClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		// for i, item := range pods.Items {
		// 	omcplog.V(0).Info(item.Spec.NodeName)
		// 	omcplog.V(0).Info("--------- ", i)
		// 	for j, container := range item.Spec.Containers {
		// 		imagename := strings.Split(container.Image, ":")
		// 		omcplog.V(0).Info(j, "-----", imagename)
		// 		_, ok := imagelist[imagename[0]]
		// 		if ok {
		// 			imagelist[imagename[0]] += 1
		// 		} else {
		// 			imagelist[imagename[0]] = 1
		// 		}
		// 	}
		// }
		ghostclients = append(ghostclients, ghostclient)

	}
	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})

	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		omcplog.V(0).Info("adding APIs to live cluster's scheme: ", err)
		return nil, err
	}
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &v1alpha1.Cache{}, controller.WatchOptions{}); err != nil {
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

type imageInfo struct {
	clusterName  string
	nodeName     string
	imageName    string
	imageVersion string
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

	instance := &v1alpha1.Cache{}
	err := r.live.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		omcplog.V(0).Info("get instance error")
		r.MakeStatus(instance, false, "", err)
	}
	if instance.Status.Succeeded == true {
		omcplog.V(4).Info(instance.Name + " running")
		r.MakeStatus(instance, true, "", err)
	}
	if instance.Status.Succeeded == false && instance.Status.Reason != "" {
		// 이미 실패한 케이스는 로직을 다시 안탄다.
		omcplog.V(4).Info(instance.Name + " already failed")
		return reconcile.Result{Requeue: false}, nil
	}

	r.Run(instance)
	omcplog.V(4).Info("end")

	ti, err := strconv.ParseInt(instance.Spec.Timer, 10, 64)
	if err != nil {
		omcplog.V(4).Info("Timer error")
		return reconcile.Result{}, err
	}

	time.Sleep(time.Duration(ti) * time.Minute)

	return reconcile.Result{Requeue: true}, nil
}

type kv struct {
	Key   string
	Value int
}

func sortValue(m map[string]int) []kv {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	return ss
}

func (r *reconciler) Run(instance *v1alpha1.Cache) (bool, error) {
	clusternamelist := []string{}
	clusterInstanceList := &clusterv1alpha1.OpenMCPClusterList{}
	err := r.live.List(context.TODO(), clusterInstanceList, &client.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			//r.DeleteOpenMCPCluster(cm, request.Namespace, request.Name)
			return true, nil
		}
		return false, err
	}
	for _, item := range clusterInstanceList.Items {
		if item.Spec.JoinStatus == "JOIN" {
			clusternamelist = append(clusternamelist, item.Name)
		}
	}
	imageList := make(map[string]imageInfo)
	imageCount := make(map[string]int)
	for _, clientName := range clusternamelist {
		listClient := *cm.Cluster_kubeClients[clientName]
		pods, _ := listClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		for _, item := range pods.Items {
			// omcplog.V(0).Info(item.Spec.NodeName)
			for _, container := range item.Spec.Containers {
				imagefullname := strings.Split(container.Image, ":")
				imagename := imagefullname[0]
				if imagename == "docker" {
					continue
				}
				imagename = strings.Replace(imagename, "docker.io/", "", -1)
				omcplog.V(3).Info("image name : ", imagename)
				imageversion := "1"
				imageinfo := imageInfo{}
				imageinfo.clusterName = clientName
				imageinfo.nodeName = item.Spec.NodeName
				imageinfo.imageName = imagename
				imageinfo.imageVersion = imageversion
				_, ok := imageCount[imagename]
				if ok {
					imageCount[imagename] += 1
				} else {
					imageCount[imagename] = 1
					imageList[imagename] = imageinfo
				}
			}
		}
	}
	var registryManager nodeapi.RegistryManager
	cachecount := instance.Spec.CacheCount
	sortedImageList := sortValue(imageCount)
	for _, item := range sortedImageList[:cachecount] {
		omcplog.V(3).Info("imagecache select image : ", item.Key)
		imagename := item.Key
		imageInfo := imageList[imagename]
		err = registryManager.Init(imageInfo.clusterName, cm)
		if err != nil {
			r.MakeStatusWithSource(instance, false, instance.Spec, err, err)
			return false, err
		}
		err := registryManager.CreatePushJob(imageInfo.nodeName, imageInfo.imageName)
		if err != nil {
			r.MakeStatusWithSource(instance, false, instance.Spec, err, err)
			return false, err
		}
	}

	omcplog.V(3).Info("image push job complete!")
	for _, pullList := range clusternamelist {
		listClient := *cm.Cluster_kubeClients[pullList]
		err = registryManager.Init(pullList, cm)
		if err != nil {
			r.MakeStatusWithSource(instance, false, instance.Spec, err, err)
			return false, err
		}
		nodelist, _ := listClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		for _, item := range nodelist.Items {
			omcplog.V(3).Info("nodename : ", item.Name)
			for _, image := range sortedImageList[:cachecount] {
				err = registryManager.CreatePullJob(item.Name, image.Key)
				if err != nil {
					r.MakeStatusWithSource(instance, false, instance.Spec, err, err)
					return false, err
				}
			}
		}
		err := registryManager.DeleteJob()
		if err != nil {
			r.MakeStatusWithSource(instance, false, instance.Spec, err, err)
			return false, err
		}
	}
	omcplog.V(3).Info("image pull job complete!")

	data := v1alpha1.Data{}
	for _, item := range sortedImageList[:cachecount] {
		imageinfo := v1alpha1.ImageInfo{
			ImageName:  item.Key,
			ImageCount: int64(item.Value),
		}
		data.ImageList = append(data.ImageList, imageinfo)
	}
	data.Timestamp = time.Now().Local().Format("2006-01-02 15:04:05")
	instance.Status.History = append(instance.Status.History, data)
	r.MakeStatusWithSource(instance, true, instance.Spec, nil, nil)
	// push, pull - nodeName 이 없을 경우 Cluster 단위 명령
	return true, nil
}

func (r *reconciler) MakeStatusWithSource(instance *v1alpha1.Cache, CacheStatus bool, CacheSpec v1alpha1.CacheSpec, err error, detailErr error) {
	r.makeStatusRun(instance, CacheStatus, CacheSpec, "", err, detailErr)
}

func (r *reconciler) MakeStatus(instance *v1alpha1.Cache, CacheStatus bool, elapsed string, err error) {
	r.makeStatusRun(instance, CacheStatus, v1alpha1.CacheSpec{}, elapsed, err, nil)
}

func (r *reconciler) makeStatusRun(instance *v1alpha1.Cache, CacheStatus bool, CacheSpec v1alpha1.CacheSpec, elapsedTime string, err error, detailErr error) {
	instance.Status.Succeeded = CacheStatus

	if elapsedTime == "" {
		elapsedTime = "0"
	}

	omcplog.V(3).Info("CacheStatus : ", CacheStatus)

	if !CacheStatus {
		omcplog.V(3).Info("err : ", err.Error())
		tmp := make(map[string]interface{})

		//tmp["ResourceType"] = snapshotSource.ResourceType
		//tmp["ResourceName"] = snapshotSource.ResourceName
		//tmp["VolumeSnapshotClassName"] = snapshotSource.VolumeDataSource.VolumeSnapshotClassName
		//tmp["VolumeSnapshotSourceKind"] = snapshotSource.VolumeDataSource.VolumeSnapshotSourceKind
		//tmp["VolumeSnapshotSourceName"] = snapshotSource.VolumeDataSource.VolumeSnapshotSourceName
		//tmp["VolumeSnapshotKey"] = instance.Status.VolumeDataSource.VolumeSnapshotKey
		tmp["Reason"] = err.Error()
		//tmp["ReasonDetail"] = detailErr.Error()

		jsonTmp, err := json.Marshal(tmp)
		if err != nil {
			omcplog.V(3).Info(err, "-----------")
		}
		instance.Status.Reason = string(jsonTmp)
		if detailErr != nil {
			instance.Status.Reason = detailErr.Error()
		}
	}
	omcplog.V(3).Info("history log : ", instance.Status.History)
	//r.live.Update(context.TODO(), instance)
	//r.live.Status().Patch(context.TODO(), instance)
	//r.live.Status().Update(context.TODO(), instance)
	//err = r.live.Status().Update(context.TODO(), instance)
	omcplog.V(3).Info("live update")
	err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	err = r.live.Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}

	omcplog.V(3).Info("live update end")
}
