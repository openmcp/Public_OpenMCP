package syncWeight

import (
	"fmt"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"time"

	"openmcp/openmcp/openmcp-loadbalancing-controller/src/controller/OpenMCPVirtualService"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func SyncVSWeight(myClusterManager *clusterManager.ClusterManager, quit, quitok chan bool) {
	cm := myClusterManager
	for {
		select {
		case <-quit:
			omcplog.V(2).Info("SyncWeight Quit")
			quitok <- true
			return
		default:
			fmt.Println("SyncWeight Called")

			//vsList, err := cm.Crd_istio_client.VirtualService(corev1.NamespaceAll).List(metav1.ListOptions{})

			ovsList, err := cm.Crd_client.OpenMCPVirtualService(corev1.NamespaceAll).List(metav1.ListOptions{})
			if err != nil {
				omcplog.V(0).Info("Error:", err)
				time.Sleep(time.Second * 5)
				continue
			}
			if len(ovsList.Items) == 0 {
				omcplog.V(0).Info("Not Exist OpenMCP VirtualService List")
				time.Sleep(time.Second * 5)
				continue
			}
			for _, ovs := range ovsList.Items {

				checkVs, err := cm.Crd_istio_client.VirtualService(ovs.Namespace).Get(ovs.Name, metav1.GetOptions{})

				if err == nil || (err != nil && errors.IsNotFound(err)) {
					vs, err2 := OpenMCPVirtualService.MakeVirtualService(&ovs)
					if err2 != nil {
						omcplog.V(0).Info("Error:", err2)

					}
					if err == nil {
						// Update VirtualService
						vs.ResourceVersion = checkVs.ResourceVersion
						_, err3 := cm.Crd_istio_client.VirtualService(vs.Namespace).Update(vs)
						if err3 != nil {
							omcplog.V(0).Info("Error:", err3)
						}
					} else if err != nil && errors.IsNotFound(err) {
						// Create VirtualService
						_, err3 := cm.Crd_istio_client.VirtualService(vs.Namespace).Create(vs)
						if err3 != nil {
							omcplog.V(0).Info("Error:", err3)
						}
					} else {
						omcplog.V(0).Info("Error:", err)
					}

				}
			}

			//vsList := &v1alpha3.VirtualServiceList{}
			//err := cm.Host_client.List(context.TODO(), vsList, corev1.NamespaceAll)

			/*
				if len(vsList.Items) == 0 {
					fmt.Println(err)
					time.Sleep(time.Second * 5)

				} else if err != nil {
					fmt.Println(err)
					time.Sleep(time.Second * 5)
					continue

				}

				for _, vs := range vsList.Items {
					for k, http := range vs.Spec.Http {
						if len(http.Match) == 0 {
							continue
						}
						if _, ok := http.Match[0].Headers["client-region"]; !ok {
							continue
						}
						if _, ok := http.Match[0].Headers["client-zone"]; !ok {
							continue
						}

						omcplog.V(4).Info("SyncWeight setWeight : ", vs.Name, vs.Namespace)

						exactRegion := http.Match[0].Headers["client-region"].GetExact()
						exactZone := http.Match[0].Headers["client-zone"].GetExact()
						//setWeight(vs.Spec.Http[k].Route, exactRegion, exactZone, vs.Namespace)

						createVsHttpRoutes(vs.Spec.Http[k].Route, exactRegion, exactZone, vs.Namespace)

					}

					_, err := cm.Crd_istio_client.VirtualService(vs.Namespace).Update(&vs)
					if err != nil {
						fmt.Println(err)
					}

				}
			*/

			time.Sleep(time.Second * 5)
		}
	}

}
