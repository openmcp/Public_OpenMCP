package resourceCreate

import (
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-cluster-manager/src/syncResource"
	"openmcp/openmcp/util/clusterManager"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func CreateSubResourceAll(clusterName string, cm *clusterManager.ClusterManager) error {
	_ = createSubResourceDeployment(clusterName, cm)
	_ = createSubResourceService(clusterName, cm)
	_ = createSubResourceIngress(clusterName, cm)
	_ = createSubResourceConfigMap(clusterName, cm)
	_ = createSubResourceJob(clusterName, cm)
	_ = createSubResourceSecret(clusterName, cm)
	_ = createSubResourceNamespace(clusterName, cm)
	return nil
}

func createSubResourceDeployment(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	odepList, err := cm.Crd_client.OpenMCPDeployment(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, odep := range odepList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + odep.Name + "' in Cluster. Update Status Scheduling Need='true' & Complete='false'")
		odep.Status.SchedulingNeed = true
		odep.Status.SchedulingComplete = false
		_, err2 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).UpdateStatus(&odep)
		if err2 != nil {
			return err2
		}

	}
	return nil
}

func createSubResourceService(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	osvcList, err := cm.Crd_client.OpenMCPService(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, osvc := range osvcList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + osvc.Name + "' in Cluster. Update Status ClusterMap")
		//osvc.Status.ClusterMaps[clusterName] = 1
		osvc.Status.ClusterMaps = syncResource.SyncResource(cm)
		_, err2 := cm.Crd_client.OpenMCPService(osvc.Namespace).UpdateStatus(&osvc)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
func createSubResourceIngress(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	oingList, err := cm.Crd_client.OpenMCPIngress(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, oing := range oingList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + oing.Name + "' in Cluster. Update Status ClusterMap")
		//oing.Status.ClusterMaps[clusterName] = 1
		oing.Status.ClusterMaps = syncResource.SyncResource(cm)
		_, err2 := cm.Crd_client.OpenMCPIngress(oing.Namespace).UpdateStatus(&oing)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
func createSubResourceConfigMap(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	ocmList, err := cm.Crd_client.OpenMCPConfigMap(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ocm := range ocmList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + ocm.Name + "' in Cluster. Update Status ClusterMap")
		//ocm.Status.ClusterMaps[clusterName] = 1
		ocm.Status.ClusterMaps = syncResource.SyncResource(cm)
		_, err2 := cm.Crd_client.OpenMCPConfigMap(ocm.Namespace).UpdateStatus(&ocm)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
func createSubResourceJob(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	ojobList, err := cm.Crd_client.OpenMCPJob(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ojob := range ojobList.Items {
		//ojob.Status.ClusterMaps[clusterName] = 1
		ojob.Status.ClusterMaps = syncResource.SyncResource(cm)
		omcplog.V(2).Info("[Resource Create] Found Deployment '" + ojob.Name + "' in Cluster. Update Status ClusterMap")
		_, err2 := cm.Crd_client.OpenMCPJob(ojob.Namespace).UpdateStatus(&ojob)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
func createSubResourceNamespace(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	onsList, err := cm.Crd_client.OpenMCPNamespace(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ons := range onsList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + ons.Name + "' in Cluster. Update Status ClusterMap")
		//ons.Status.ClusterMaps[clusterName] = 1
		ons.Status.ClusterMaps = syncResource.SyncResource(cm)
		_, err2 := cm.Crd_client.OpenMCPNamespace(ons.Namespace).UpdateStatus(&ons)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
func createSubResourceSecret(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Create]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Create] Deployment Resource by OpenMCP")
	osecList, err := cm.Crd_client.OpenMCPSecret(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, osec := range osecList.Items {

		omcplog.V(2).Info("[Resource Create] Found Deployment '" + osec.Name + "' in Cluster. Update Status ClusterMap")
		//osec.Status.ClusterMaps[clusterName] = 1
		osec.Status.ClusterMaps = syncResource.SyncResource(cm)
		_, err2 := cm.Crd_client.OpenMCPSecret(osec.Namespace).UpdateStatus(&osec)
		if err2 != nil {
			return err2
		}

	}
	return nil
}
