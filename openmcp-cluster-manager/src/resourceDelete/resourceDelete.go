package resourceDelete

import (
	"context"
	"fmt"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-cluster-manager/src/syncResource"
	"openmcp/openmcp/util/clusterManager"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func DeleteSubResourceAll(clusterName string, cm *clusterManager.ClusterManager) error {
	_ = deleteSubResourceDeployment(clusterName, cm)
	_ = deleteSubResourceService(clusterName, cm)
	_ = deleteSubResourceIngress(clusterName, cm)
	_ = deleteSubResourceConfigMap(clusterName, cm)
	_ = deleteSubResourceJob(clusterName, cm)
	_ = deleteSubResourceSecret(clusterName, cm)
	_ = deleteSubResourceNamespace(clusterName, cm)
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

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Service '" + osvc.Name + "' in Cluster")
			svc := &corev1.Service{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), svc, osvc.Namespace, osvc.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Service. Continue Next Service")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Service '" + osvc.Name + "' in Cluster. Update BlockSubResource = true")
				osvc.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPService(osvc.Namespace).UpdateStatus(&osvc)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + osvc.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), svc, osvc.Namespace, osvc.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Service '" + osvc.Name + "' in Cluster")

				changed_osvc, err5 := cm.Crd_client.OpenMCPService(osvc.Namespace).Get(osvc.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_osvc.Status.ClusterMaps[clusterName] -= 1
				changed_osvc.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_osvc.Status.BlockSubResource = false
				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPService(changed_osvc.Namespace).UpdateStatus(changed_osvc)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPService Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceIngress(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Ingress Resource by OpenMCP")
	oingList, err := cm.Crd_client.OpenMCPIngress(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oing := range oingList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Ingress '" + oing.Name + "' in Cluster")
			ing := &extv1b1.Ingress{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ing, oing.Namespace, oing.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Ingress. Continue Next Ingress")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Ingress '" + oing.Name + "' in Cluster. Update BlockSubResource = true")
				oing.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPIngress(oing.Namespace).UpdateStatus(&oing)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + oing.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ing, oing.Namespace, oing.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Ingress '" + oing.Name + "' in Cluster")

				changed_oing, err5 := cm.Crd_client.OpenMCPIngress(oing.Namespace).Get(oing.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_oing.Status.ClusterMaps[clusterName] -= 1
				changed_oing.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oing.Status.BlockSubResource = false

				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPIngress(changed_oing.Namespace).UpdateStatus(changed_oing)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPIngress Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceConfigMap(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] ConfigMap Resource by OpenMCP")
	oConfigmapList, err := cm.Crd_client.OpenMCPConfigMap(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oConfigmap := range oConfigmapList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find ConfigMap '" + oConfigmap.Name + "' in Cluster")
			configmap := &corev1.ConfigMap{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), configmap, oConfigmap.Namespace, oConfigmap.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found ConfigMap. Continue Next ConfigMap")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found ConfigMap '" + oConfigmap.Name + "' in Cluster. Update BlockSubResource = true")
				oConfigmap.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPConfigMap(oConfigmap.Namespace).UpdateStatus(&oConfigmap)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + oConfigmap.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), configmap, oConfigmap.Namespace, oConfigmap.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted ConfigMap '" + oConfigmap.Name + "' in Cluster")

				changed_oConfigmap, err5 := cm.Crd_client.OpenMCPConfigMap(oConfigmap.Namespace).Get(oConfigmap.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_oConfigmap.Status.ClusterMaps[clusterName] -= 1
				changed_oConfigmap.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oConfigmap.Status.BlockSubResource = false

				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPConfigMap(changed_oConfigmap.Namespace).UpdateStatus(changed_oConfigmap)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPConfigMap Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceJob(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Job Resource by OpenMCP")
	oJobList, err := cm.Crd_client.OpenMCPJob(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oJob := range oJobList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Job '" + oJob.Name + "' in Cluster")
			job := &batchv1.Job{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), job, oJob.Namespace, oJob.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Job. Continue Next Job")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Job '" + oJob.Name + "' in Cluster. Update BlockSubResource = true")
				oJob.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPJob(oJob.Namespace).UpdateStatus(&oJob)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + oJob.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), job, oJob.Namespace, oJob.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Job '" + oJob.Name + "' in Cluster")

				changed_oJob, err5 := cm.Crd_client.OpenMCPJob(oJob.Namespace).Get(oJob.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_oJob.Status.ClusterMaps[clusterName] -= 1
				changed_oJob.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oJob.Status.BlockSubResource = false

				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPJob(changed_oJob.Namespace).UpdateStatus(changed_oJob)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPJob Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceNamespace(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Namespace Resource by OpenMCP")
	onsList, err := cm.Crd_client.OpenMCPNamespace(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ons := range onsList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Namespace '" + ons.Name + "' in Cluster")
			ns := &corev1.Namespace{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ns, ons.Namespace, ons.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Namespace. Continue Next Namespace")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Namespace '" + ons.Name + "' in Cluster. Update BlockSubResource = true")
				ons.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPNamespace(ons.Namespace).UpdateStatus(&ons)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + ons.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ns, ons.Namespace, ons.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Job '" + ons.Name + "' in Cluster")

				changed_ons, err5 := cm.Crd_client.OpenMCPNamespace(ons.Namespace).Get(ons.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_ons.Status.ClusterMaps[clusterName] -= 1
				changed_ons.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_ons.Status.BlockSubResource = false

				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPNamespace(changed_ons.Namespace).UpdateStatus(changed_ons)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPNamespace Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceSecret(clusterName string, cm *clusterManager.ClusterManager) error {
	omcplog.V(2).Info("[Resource Clean]'" + clusterName + "' Start")
	omcplog.V(2).Info("[Resource Clean] Secret Resource by OpenMCP")
	osecList, err := cm.Crd_client.OpenMCPSecret(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, osec := range osecList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Resource Clean] Try Find Secret '" + osec.Name + "' in Cluster")
			sec := &corev1.Secret{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), sec, osec.Namespace, osec.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Resource Clean] is Not Found Secret. Continue Next Secret")

			} else if err3 == nil {

				omcplog.V(2).Info("[Resource Clean] Found Secret '" + osec.Name + "' in Cluster. Update BlockSubResource = true")
				osec.Status.BlockSubResource = true
				_, err2 := cm.Crd_client.OpenMCPSecret(osec.Namespace).UpdateStatus(&osec)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Resource Clean] '" + osec.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), sec, osec.Namespace, osec.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Resource Clean] Deleted Job '" + osec.Name + "' in Cluster")

				changed_osec, err5 := cm.Crd_client.OpenMCPSecret(osec.Namespace).Get(osec.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Resource Clean] Error: ", err5)
				}

				// changed_osec.Status.ClusterMaps[clusterName] -= 1
				changed_osec.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_osec.Status.BlockSubResource = false

				omcplog.V(2).Info("[Resource Clean] ReSchduling And BlockSubResource = false")
				_, err6 := cm.Crd_client.OpenMCPSecret(changed_osec.Namespace).UpdateStatus(changed_osec)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Resource Clean] Done. OpenMCPSecret Stauts Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
