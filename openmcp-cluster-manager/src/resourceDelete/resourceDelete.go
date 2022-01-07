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

	v1 "k8s.io/api/core/v1"
)

func DeleteSubResourceAll(clusterName string, cm *clusterManager.ClusterManager) error {
	_ = deleteSubResourceDeployment(clusterName, cm)
	//	_ = deleteSubResourceService(clusterName, cm)
	//	_ = deleteSubResourceIngress(clusterName, cm)
	_ = deleteSubResourceConfigMap(clusterName, cm)
	_ = deleteSubResourceJob(clusterName, cm)
	_ = deleteSubResourceSecret(clusterName, cm)
	_ = deleteSubResourceNamespace(clusterName, cm)

	_ = deleteSubResourcePV(clusterName, cm)
	_ = deleteSubResourcePVC(clusterName, cm)
	_ = deleteSubResourceStatefulSet(clusterName, cm)
	_ = deleteSubResourceDaemonSet(clusterName, cm)

	return nil
}

func deleteSubResourceDeployment(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Deployment Resource by OpenMCP")
	odepList, err := cm.Crd_client.OpenMCPDeployment(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, odep := range odepList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Deployment '" + odep.Name + "' in Cluster")
			dep := &appsv1.Deployment{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), dep, odep.Namespace, odep.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Deployment. Continue Next Deployment")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Deployment '" + odep.Name + "' in Cluster. Update CheckSubResource = false")
				odep.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).UpdateStatus(&odep)
				if err2 != nil {
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + odep.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), dep, odep.Namespace, odep.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Deployment '" + odep.Name + "' in Cluster")

				changed_odep, err5 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).Get(odep.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				changed_odep.Status.ClusterMaps[clusterName] -= *dep.Spec.Replicas
				changed_odep.Status.SchedulingNeed = true
				changed_odep.Status.SchedulingComplete = false
				changed_odep.Status.CheckSubResource = true
				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPDeployment(changed_odep.Namespace).UpdateStatus(changed_odep)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPDeployment Status Changed")

			} else {
				return err3
			}
		}

		// changed_odep, err5 := cm.Crd_client.OpenMCPDeployment(odep.Namespace).Get(odep.Name, metav1.GetOptions{})
		// if err5 != nil {
		// 	omcplog.V(0).Info("[Delete Resource] Error: ", err5)
		// }
		// changed_odep.Status.CheckSubResource = true

	}
	return nil
}

func deleteSubResourceService(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Service Resource by OpenMCP")
	osvcList, err := cm.Crd_client.OpenMCPService(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, osvc := range osvcList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Service '" + osvc.Name + "' in Cluster")
			svc := &corev1.Service{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), svc, osvc.Namespace, osvc.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Service. Continue Next Service")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Service '" + osvc.Name + "' in Cluster. Update CheckSubResource = false")
				osvc.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPService(osvc.Namespace).UpdateStatus(&osvc)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + osvc.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), svc, osvc.Namespace, osvc.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Service '" + osvc.Name + "' in Cluster")

				changed_osvc, err5 := cm.Crd_client.OpenMCPService(osvc.Namespace).Get(osvc.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_osvc.Status.ClusterMaps[clusterName] -= 1
				changed_osvc.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_osvc.Status.CheckSubResource = true
				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPService(changed_osvc.Namespace).UpdateStatus(changed_osvc)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPService Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceIngress(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Ingress Resource by OpenMCP")
	oingList, err := cm.Crd_client.OpenMCPIngress(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oing := range oingList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Ingress '" + oing.Name + "' in Cluster")
			ing := &extv1b1.Ingress{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ing, oing.Namespace, oing.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Ingress. Continue Next Ingress")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Ingress '" + oing.Name + "' in Cluster. Update CheckSubResource = false")
				oing.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPIngress(oing.Namespace).UpdateStatus(&oing)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + oing.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ing, oing.Namespace, oing.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Ingress '" + oing.Name + "' in Cluster")

				changed_oing, err5 := cm.Crd_client.OpenMCPIngress(oing.Namespace).Get(oing.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_oing.Status.ClusterMaps[clusterName] -= 1
				changed_oing.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oing.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPIngress(changed_oing.Namespace).UpdateStatus(changed_oing)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPIngress Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceConfigMap(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] ConfigMap Resource by OpenMCP")
	oConfigmapList, err := cm.Crd_client.OpenMCPConfigMap(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oConfigmap := range oConfigmapList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find ConfigMap '" + oConfigmap.Name + "' in Cluster")
			configmap := &corev1.ConfigMap{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), configmap, oConfigmap.Namespace, oConfigmap.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found ConfigMap. Continue Next ConfigMap")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found ConfigMap '" + oConfigmap.Name + "' in Cluster. Update CheckSubResource = false")
				oConfigmap.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPConfigMap(oConfigmap.Namespace).UpdateStatus(&oConfigmap)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + oConfigmap.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), configmap, oConfigmap.Namespace, oConfigmap.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted ConfigMap '" + oConfigmap.Name + "' in Cluster")

				changed_oConfigmap, err5 := cm.Crd_client.OpenMCPConfigMap(oConfigmap.Namespace).Get(oConfigmap.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_oConfigmap.Status.ClusterMaps[clusterName] -= 1
				changed_oConfigmap.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oConfigmap.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPConfigMap(changed_oConfigmap.Namespace).UpdateStatus(changed_oConfigmap)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPConfigMap Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceJob(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Job Resource by OpenMCP")
	oJobList, err := cm.Crd_client.OpenMCPJob(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oJob := range oJobList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Job '" + oJob.Name + "' in Cluster")
			job := &batchv1.Job{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), job, oJob.Namespace, oJob.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Job. Continue Next Job")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Job '" + oJob.Name + "' in Cluster. Update CheckSubResource = false")
				oJob.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPJob(oJob.Namespace).UpdateStatus(&oJob)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + oJob.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), job, oJob.Namespace, oJob.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Job '" + oJob.Name + "' in Cluster")

				changed_oJob, err5 := cm.Crd_client.OpenMCPJob(oJob.Namespace).Get(oJob.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_oJob.Status.ClusterMaps[clusterName] -= 1
				changed_oJob.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oJob.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPJob(changed_oJob.Namespace).UpdateStatus(changed_oJob)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPJob Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceNamespace(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Namespace Resource by OpenMCP")
	onsList, err := cm.Crd_client.OpenMCPNamespace(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ons := range onsList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Namespace '" + ons.Name + "' in Cluster")
			ns := &corev1.Namespace{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ns, ons.Namespace, ons.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Namespace. Continue Next Namespace")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Namespace '" + ons.Name + "' in Cluster. Update CheckSubResource = false")
				ons.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPNamespace(ons.Namespace).UpdateStatus(&ons)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + ons.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ns, ons.Namespace, ons.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Namespace '" + ons.Name + "' in Cluster")

				changed_ons, err5 := cm.Crd_client.OpenMCPNamespace(ons.Namespace).Get(ons.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_ons.Status.ClusterMaps[clusterName] -= 1
				changed_ons.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_ons.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPNamespace(changed_ons.Namespace).UpdateStatus(changed_ons)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPNamespace Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
func deleteSubResourceSecret(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Secret Resource by OpenMCP")
	osecList, err := cm.Crd_client.OpenMCPSecret(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, osec := range osecList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find Secret '" + osec.Name + "' in Cluster")
			sec := &corev1.Secret{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), sec, osec.Namespace, osec.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found Secret. Continue Next Secret")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found Secret '" + osec.Name + "' in Cluster. Update CheckSubResource = false")
				osec.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPSecret(osec.Namespace).UpdateStatus(&osec)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + osec.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), sec, osec.Namespace, osec.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted Secret '" + osec.Name + "' in Cluster")

				changed_osec, err5 := cm.Crd_client.OpenMCPSecret(osec.Namespace).Get(osec.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_osec.Status.ClusterMaps[clusterName] -= 1
				changed_osec.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_osec.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPSecret(changed_osec.Namespace).UpdateStatus(changed_osec)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPSecret Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}

func deleteSubResourcePV(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] PV Resource by OpenMCP")
	opvList, err := cm.Crd_client.OpenMCPPersistentVolume(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, opv := range opvList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find PV '" + opv.Name + "' in Cluster")
			pv := &v1.PersistentVolume{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), pv, opv.Namespace, opv.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found PV. Continue Next PV")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found PV '" + opv.Name + "' in Cluster. Update CheckSubResource = false")
				opv.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPPersistentVolume(opv.Namespace).UpdateStatus(&opv)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + opv.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), pv, opv.Namespace, opv.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted PV '" + opv.Name + "' in Cluster")

				changed_opv, err5 := cm.Crd_client.OpenMCPPersistentVolume(opv.Namespace).Get(opv.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_opv.Status.ClusterMaps[clusterName] -= 1
				changed_opv.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_opv.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPPersistentVolume(changed_opv.Namespace).UpdateStatus(changed_opv)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPPersistentVolume Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}

func deleteSubResourcePVC(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] PVC Resource by OpenMCP")
	opvcList, err := cm.Crd_client.OpenMCPPersistentVolumeClaim(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, opvc := range opvcList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find PVC '" + opvc.Name + "' in Cluster")
			pvc := &v1.PersistentVolumeClaim{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), pvc, opvc.Namespace, opvc.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found PVC. Continue Next PVC")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found PV '" + opvc.Name + "' in Cluster. Update CheckSubResource = false")
				opvc.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPPersistentVolumeClaim(opvc.Namespace).UpdateStatus(&opvc)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + opvc.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), pvc, opvc.Namespace, opvc.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted PVC '" + opvc.Name + "' in Cluster")

				changed_opvc, err5 := cm.Crd_client.OpenMCPPersistentVolumeClaim(opvc.Namespace).Get(opvc.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_opvc.Status.ClusterMaps[clusterName] -= 1
				changed_opvc.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_opvc.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPPersistentVolumeClaim(changed_opvc.Namespace).UpdateStatus(changed_opvc)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPPersistentVolumeClaim Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}

func deleteSubResourceStatefulSet(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] Statefulset Resource by OpenMCP")
	ossList, err := cm.Crd_client.OpenMCPStatefulSet(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, oss := range ossList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find StatefulSet '" + oss.Name + "' in Cluster")
			ss := &appsv1.StatefulSet{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ss, oss.Namespace, oss.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found StatefulSet. Continue Next StatefulSet")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found StatefulSet '" + oss.Name + "' in Cluster. Update CheckSubResource = false")
				oss.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPStatefulSet(oss.Namespace).UpdateStatus(&oss)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + oss.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ss, oss.Namespace, oss.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted StatefulSet '" + oss.Name + "' in Cluster")

				changed_oss, err5 := cm.Crd_client.OpenMCPStatefulSet(oss.Namespace).Get(oss.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_oss.Status.ClusterMaps[clusterName] -= 1
				changed_oss.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_oss.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPStatefulSet(changed_oss.Namespace).UpdateStatus(changed_oss)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPStatefulSet Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}

func deleteSubResourceDaemonSet(clusterName string, cm *clusterManager.ClusterManager) error {
	//omcplog.V(2).Info("[Delete Resource]'" + clusterName + "' Start")
	//omcplog.V(2).Info("[Delete Resource] DaemonSet Resource by OpenMCP")
	odsList, err := cm.Crd_client.OpenMCPDaemonSet(corev1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ods := range odsList.Items {

		if _, ok := cm.Cluster_genClients[clusterName]; ok {

			omcplog.V(2).Info("[Delete Resource] Try Find DaemonSet '" + ods.Name + "' in Cluster")
			ds := &appsv1.DaemonSet{}
			err3 := cm.Cluster_genClients[clusterName].Get(context.TODO(), ds, ods.Namespace, ods.Name)
			if err3 != nil && errors.IsNotFound(err3) {
				omcplog.V(2).Info("[Delete Resource] is Not Found DaemonSet. Continue Next DaemonSet")

			} else if err3 == nil {

				omcplog.V(2).Info("[Delete Resource] Found DaemonSet '" + ods.Name + "' in Cluster. Update CheckSubResource = false")
				ods.Status.CheckSubResource = false
				_, err2 := cm.Crd_client.OpenMCPDaemonSet(ods.Namespace).UpdateStatus(&ods)
				if err2 != nil {
					fmt.Println(err2)
					return err2
				}
				omcplog.V(2).Info("[Delete Resource] '" + ods.Name + "' in Cluster Delete Start")

				err4 := cm.Cluster_genClients[clusterName].Delete(context.TODO(), ds, ods.Namespace, ods.Name)
				if err4 != nil {
					return err4
				}
				omcplog.V(2).Info("[Delete Resource] Deleted DaemonSet '" + ods.Name + "' in Cluster")

				changed_ods, err5 := cm.Crd_client.OpenMCPDaemonSet(ods.Namespace).Get(ods.Name, metav1.GetOptions{})
				if err5 != nil {
					omcplog.V(0).Info("[Delete Resource] Error: ", err5)
				}

				// changed_ods.Status.ClusterMaps[clusterName] -= 1
				changed_ods.Status.ClusterMaps = syncResource.SyncResource(cm)
				changed_ods.Status.CheckSubResource = true

				omcplog.V(2).Info("[Delete Resource] ReScheduling And CheckSubResource = true")
				_, err6 := cm.Crd_client.OpenMCPDaemonSet(changed_ods.Namespace).UpdateStatus(changed_ods)
				if err6 != nil {
					omcplog.V(0).Info(err6)
				}
				omcplog.V(2).Info("[Delete Resource] Done. OpenMCPDaemonSet Status Changed")

			} else {
				return err3
			}
		}

	}
	return nil
}
