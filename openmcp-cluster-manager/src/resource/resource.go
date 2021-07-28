package resource

import (
	"context"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func DeleteSubResourceAll(clusterName string, cm *clusterManager.ClusterManager) error {
	_ = deleteSubResourceDeployment(clusterName, cm)
	_ = deleteSubResourceService(clusterName, cm)
	return nil
}

func deleteSubResourceDeployment(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Deployment Resource by OpenMCP")
	odepList, err := cm.Crd_client.OpenMCPDeployment(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, odep := range odepList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Deployment '" + odep.Name + "' in Cluster")
			dep := &appsv1.Deployment{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), dep, odep.Namespace, odep.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Deployment. Continue Next Deployment")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Deployment '" + odep.Name + "' in Cluster. Update BlockSubResource = true")
				odep.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).UpdateStatus(&odep)
				if err2 != nil {
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + odep.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), dep, odep.Namespace, odep.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Deployment '" + odep.Name + "' in Cluster")

				changed_odep, err5 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).Get(odep.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				changed_odep.Status.ClusterMaps[clusterName] -= *dep.Spec.Replicas
				changed_odep.Status.SchedulingNeed = true
				changed_odep.Status.SchedulingComplete = false
				changed_odep.Status.BlockSubResource = false
				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPDeployment(changed_odep.Namespace).UpdateStatus(changed_odep)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPDeployment Stauts Changed")

			} else {
				return err3
			}
		}

		// changed_odep, err5 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).Get(odep.Name, metav1.GetOptions{})
		// if err5 != nil {
		// 	omcplog.V(0).Info("[Resource Clean] Error: ", err5)
		// }
		// changed_odep.Status.BlockSubResource = false

	}
	return nil
}

func deleteSubResourceService(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Service Resource by OpenMCP")
	osvcList, err := cm.Crd_client.OpenMCPService(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, osvc := range osvcList.Items {

		omcplog.V(2).Info("[Resource Clean] '" + osvc.Name + "' Update BlockSubResource = true")
		osvc.Status.BlockSubResource = true
		_, err2 := cm.Crd_client.OpenMCPService(osvc.Namespace).UpdateStatus(&osvc)
		if err2 != nil {
			return err2
		}

		if _, ok := cm.Cluster_genClients[clusterName]; ok {
			omcplog.V(2).Info("[Resource Clean] Find Service '" + osvc.Name + "' in Cluster")
			svc := &corev1.Service{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), svc, osvc.Namespace, osvc.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Service. Continue Next Service")

			} else if err3 == nil {
				omcplog.V(2).Info("[Resource Clean] Find Service '" + osvc.Name + "' in Cluster. Delete Start")
				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), svc, osvc.Namespace, osvc.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleteed Service '" + osvc.Name + "' in Cluster")

				osvc.Status.ClusterMaps[clusterName] -= 1

			} else {
				return err3
			}
		}

		omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")

		changed_osvc, err5 := cm.Crd_client.OpenMCPService(osvc.Namespace).Get(osvc.Name, metav1.GetOptions{})
		if err5 != nil {
			omcplog.V(0).Info("[Resource Clean] Error: ", err5)
		}
		changed_osvc.Status.BlockSubResource = false
		_, err6 := cm.Crd_client.OpenMCPService(osvc.Namespace).UpdateStatus(changed_osvc)
		if err6 != nil {
			omcplog.V(0).Info(err6)
		}
		omcplog.V(2).Info("[Resource Clean] Done. OpenMCPService Stauts Changed")

	}
	return nil
}
