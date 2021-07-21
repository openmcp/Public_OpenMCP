package openmcpcluster

import (
	"context"
	"fmt"
	"io/ioutil"
	"openmcp/openmcp/apis"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"path/filepath"
	"strings"
	"time"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var cm *clusterManager.ClusterManager

var r = &reconciler{}

type reconciler struct {
	live           client.Client
	ghosts         map[string]client.Client
	ghostNamespace string
}

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string, myClusterManager *clusterManager.ClusterManager) (*controller.Controller, error) {
	cm = myClusterManager

	liveClient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}

	ghostClients := map[string]client.Client{}
	for _, ghost := range ghosts {
		ghostTmp, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostClients[ghost.Name] = ghostTmp
	}

	r.live = liveClient
	r.ghosts = ghostClients
	r.ghostNamespace = ghostNamespace

	co := controller.New(r, controller.Options{})

	//live.GetScheme() - apis scheme ADD
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	//omcplog.V(4).Info("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(context.TODO(), live, &clusterv1alpha1.OpenMCPCluster{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}

	return co, nil
}

func BuildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func (r *reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	//OpenMCPCluster 리소스 변화 감지
	omcplog.V(4).Info(">> Reconcile()")

	clusterInstance := &clusterv1alpha1.OpenMCPCluster{}
	err := r.live.Get(context.TODO(), request.NamespacedName, clusterInstance)

	//OpenMCPCluster 리소스가 없는 경우, 삭제
	if err != nil {
		if errors.IsNotFound(err) {

			//r.DeleteOpenMCPCluster(cm, request.Namespace, request.Name)

			return reconcile.Result{}, nil
		}
		omcplog.V(0).Info("!!! Failed to get clusterInstance")
		return reconcile.Result{}, err
	}

	//조건 추가 - STATUS 비교
	if clusterInstance.Spec.JoinStatus == "JOIN" {
		omcplog.V(2).Info(clusterInstance.Name + " [ JOIN ]")

	} else if clusterInstance.Spec.JoinStatus == "UNJOIN" {
		omcplog.V(2).Info(clusterInstance.Name + " [ UNJOIN ]")

	} else if clusterInstance.Spec.JoinStatus == "JOINING" && clusterInstance.Spec.MetalLBRange.AddressFrom != "IP_ADDRESS_FROM" && clusterInstance.Spec.MetalLBRange.AddressTo != "IP_ADDRESS_TO" {
		joinAvailable := checkClustersJoinState(clusterInstance.Name)

		if joinAvailable == "FALSE" {
			omcplog.V(2).Info("Another process is Running. Please Try again later.")
			omcplog.V(2).Info("You can Check the status of other clusters By Executing the Following command")
			omcplog.V(2).Info("=> kubectl get openmcpcluster -n openmcp")

			return reconcile.Result{}, nil
		}

		omcplog.V(2).Info(clusterInstance.Name + " [ JOINING ] Start")
		omcplog.V(4).Info("Metallb Configmap (", clusterInstance.Spec.MetalLBRange.AddressFrom, ",", clusterInstance.Spec.MetalLBRange.AddressTo, ")")
		joinCheck := MergeConfigAndJoin(*clusterInstance)

		if joinCheck == "TRUE" {
			omcplog.V(4).Info("OpenMCP Module Deploy---")
			moduleDirectory := []string{"namespace", "custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller", "configmap", "istio"}
			for i, dirname := range moduleDirectory {
				moduleDirectory[i] = "/init/" + dirname
			}

			util.CmdExec2("cp /mnt/config $HOME/.kube/config")

			util.CmdExec2("chmod 755 " + "/init/vertical-pod-autoscaler/hack/*")
			util.CmdExec2("/init/vertical-pod-autoscaler/hack/vpa-up.sh " + clusterInstance.Name)
			//fmt.Println(moduleDirectory)

			InstallInitModule(moduleDirectory, clusterInstance.Name, clusterInstance.Spec.MetalLBRange.AddressFrom, clusterInstance.Spec.MetalLBRange.AddressTo)
			omcplog.V(4).Info("--- JOIN Complete ---")

			clusterInstance.Spec.JoinStatus = "JOIN"
			err = r.live.Update(context.TODO(), clusterInstance)
			if err != nil {
				omcplog.V(0).Info("[" + clusterInstance.Name + "] Error Status Not Changed (JOINING -> JOIN): " + err.Error())
			} else {
				omcplog.V(2).Info(clusterInstance.Name + " [ JOINING ] Done.")
			}
		}
	} else if clusterInstance.Spec.JoinStatus == "UNJOINING" {
		joinAvailable := checkClustersJoinState(clusterInstance.Name)

		if joinAvailable == "FALSE" {
			omcplog.V(2).Info("Another process is Running. Please Try again later.")
			omcplog.V(2).Info("You can Check the status of other clusters By Executing the Following command")
			omcplog.V(2).Info("=> kubectl get openmcpcluster -n openmcp")

			return reconcile.Result{}, nil
		}

		omcplog.V(2).Info(clusterInstance.Name + " [ UNJOINING ] Start")

		//config 파일 확인 (클러스터 조인 유무)
		memberkc := &cobrautil.KubeConfig{}
		err := yaml.Unmarshal(clusterInstance.Spec.KubeconfigInfo, memberkc)
		memberIP := memberkc.Clusters[0].Cluster.Server

		openmcpkc := &cobrautil.KubeConfig{}
		yamlFile, err := ioutil.ReadFile("/mnt/config")
		if err != nil {
			omcplog.V(4).Info("yamlFile.Get err   #%v ", err)
		}

		err = yaml.Unmarshal(yamlFile, openmcpkc)
		if err != nil {
			omcplog.V(4).Info("Unmarshal: %v", err)
		}

		unjoinCheck := ""

		for _, cluster := range openmcpkc.Clusters {
			lower_memberIP := strings.ToLower(memberIP)
			if strings.Contains(cluster.Cluster.Server, lower_memberIP) {
				unjoinCheck = cluster.Name
				break
			}
		}

		if unjoinCheck != "" {
			omcplog.V(4).Info("OpenMCP Module Delete---")
			moduleDirectory := []string{"istio", "custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller", "configmap", "namespace"}
			for i, dirname := range moduleDirectory {
				moduleDirectory[i] = "/init/" + dirname
			}
			util.CmdExec2("cp /mnt/config $HOME/.kube/config")
			util.CmdExec2("chmod 755 " + "/init/vertical-pod-autoscaler/hack/*")
			util.CmdExec2("/init/vertical-pod-autoscaler/hack/vpa-down.sh " + clusterInstance.Name)
			util.CmdExec2("kubectl delete --context openmcp -n istio-system secret istio-remote-secret-" + clusterInstance.Name)
			UninstallInitModule(moduleDirectory, clusterInstance.Name)

			omcplog.V(4).Info("Cluster Unjoin---")
			UnjoinAndDeleteConfig(memberkc, openmcpkc)

			omcplog.V(4).Info("--- UNJOIN Complete ---")
		} else {
			omcplog.V(4).Info("Not Exists Cluster Info")
		}
		clusterInstance.Spec.JoinStatus = "UNJOIN"
		err = r.live.Update(context.TODO(), clusterInstance)
		if err != nil {
			omcplog.V(0).Info("[" + clusterInstance.Name + "] Error Status Not Changed (UNJOINING -> UNJOIN): " + err.Error())
		} else {
			omcplog.V(2).Info(clusterInstance.Name + " [ UNJOINING ] Done.")
		}

	} else {
		omcplog.V(4).Info("NOT Ready")
	}

	return reconcile.Result{}, nil
}

func checkClustersJoinState(name string) string {
	clusterInstance := &clusterv1alpha1.OpenMCPClusterList{}
	err := r.live.List(context.TODO(), clusterInstance)

	if err != nil {
		for _, clusters := range clusterInstance.Items {
			if clusters.Name != name && (clusters.Spec.JoinStatus == "JOINING" || clusters.Spec.JoinStatus == "UNJOINING") {
				return "FALSE"
			}
		}
	}

	return "TRUE"

}

func InstallInitModule(directory []string, clustername string, ipaddressfrom string, ipaddressto string) {

	for i := 0; i < len(directory); i++ {
		dirname, _ := filepath.Abs(directory[i])
		//fmt.Println(dirname)
		fi, err := os.Stat(dirname)
		if err != nil {
			fmt.Println(err)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			files, err1 := ioutil.ReadDir(dirname)

			if err1 != nil {
				fmt.Println(err1)
			}
			for _, f := range files {
				fi, err2 := os.Stat(dirname + "/" + f.Name())
				if err2 != nil {
					fmt.Println(err2)
				}

				if fi.Mode().IsDir() {
					InstallInitModule([]string{dirname + "/" + f.Name()}, clustername, ipaddressfrom, ipaddressto)
				} else {
					if strings.Contains(f.Name(), "istio_install.sh") {
						util.CmdExec2("cp " + dirname + "/istio_install.sh " + dirname + "/istio_install-" + clustername + ".sh")
						util.CmdExec2("sed -i 's|REPLACE_DIRECTORY|" + dirname + "|g' " + dirname + "/istio_install-" + clustername + ".sh")
						util.CmdExec2("sed -i 's|REPLACE_CLUSTERNAME|" + clustername + "|g' " + dirname + "/istio_install-" + clustername + ".sh")
						util.CmdExec2("chmod 755 " + dirname + "/gen-eastwest-gateway.sh")
						util.CmdExec2("chmod 755 " + dirname + "/istio_install-" + clustername + ".sh")
						util.CmdExec2(dirname + "/istio_install-" + clustername + ".sh")
						util.CmdExec2("rm " + dirname + "/istio_install-" + clustername + ".sh")
						fmt.Println("*** ", dirname+" created")
					}
					if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
						if strings.Contains(dirname, "metric-collector/operator") {
							//fmt.Println("*** ", dirname+"/"+f.Name())
							util.CmdExec2("cp " + dirname + "/operator.yaml " + dirname + "/operator_" + clustername + ".yaml")
							util.CmdExec2("sed -i 's|REPLACE_CLUSTER_NAME|\"" + clustername + "\"|g' " + dirname + "/operator_" + clustername + ".yaml")
							util.CmdExec2("/usr/local/bin/kubectl apply -f " + dirname + "/operator_" + clustername + ".yaml --context " + clustername)
							util.CmdExec2("rm " + dirname + "/operator_" + clustername + ".yaml")
							fmt.Println("*** ", dirname+"/operator_"+clustername+" created")
						} else if strings.Contains(dirname, "metallb/configmap") {
							//fmt.Println("*** ", dirname+"/"+f.Name())
							util.CmdExec2("cp " + dirname + "/metallb_configmap.yaml " + dirname + "/metallb_configmap_" + clustername + ".yaml")
							util.CmdExec2("sed -i 's|CLUSTER_ADDRESS_FROM|" + ipaddressfrom + "|g' " + dirname + "/metallb_configmap_" + clustername + ".yaml")
							util.CmdExec2("sed -i 's|CLUSTER_ADDRESS_TO|" + ipaddressto + "|g' " + dirname + "/metallb_configmap_" + clustername + ".yaml")
							util.CmdExec2("/usr/local/bin/kubectl apply -f " + dirname + "/metallb_configmap_" + clustername + ".yaml --context " + clustername)
							util.CmdExec2("rm " + dirname + "/metallb_configmap_" + clustername + ".yaml")
							fmt.Println("*** ", dirname+"/metallb_configmap_"+clustername+" created")
						} else {
							if strings.Contains(dirname, "istio") {
							} else {
								util.CmdExec2("/usr/local/bin/kubectl apply -f " + dirname + "/" + f.Name() + " --context " + clustername)
							}
						}
					}
				}
			}
		}

	}

}

func UninstallInitModule(directory []string, clustername string) {
	for i := 0; i < len(directory); i++ {
		dirname, _ := filepath.Abs(directory[i])

		fi, err := os.Stat(dirname)
		if err != nil {
			fmt.Println(err)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			files, err1 := ioutil.ReadDir(dirname)

			if err1 != nil {
				fmt.Println(err1)
			}
			for _, f := range files {
				fi, err2 := os.Stat(dirname + "/" + f.Name())
				if err2 != nil {
					fmt.Println(err2)
				}

				if fi.Mode().IsDir() {
					UninstallInitModule([]string{dirname + "/" + f.Name()}, clustername)
				} else {
					if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
						if strings.Contains(dirname, "istio") {

						} else {
							util.CmdExec2("/usr/local/bin/kubectl delete -f " + dirname + "/" + f.Name() + " --context " + clustername)
						}
					}
				}
			}
		}
	}
}

func MergeConfigAndJoin(clusterInstance clusterv1alpha1.OpenMCPCluster) string {
	//config파일에 해당 정보가 저장되어 있는지 확인
	memberkc := &cobrautil.KubeConfig{}
	err := yaml.Unmarshal(clusterInstance.Spec.KubeconfigInfo, memberkc)
	memberIP := memberkc.Clusters[0].Cluster.Server
	//memberName := memberkc.Clusters[0].Name

	openmcpkc := &cobrautil.KubeConfig{}
	yamlFile, err := ioutil.ReadFile("/mnt/config")
	if err != nil {
		omcplog.V(4).Info("yamlFile.Get err   #%v ", err)
		return "FALSE"
	}
	err = yaml.Unmarshal(yamlFile, openmcpkc)
	if err != nil {
		omcplog.V(4).Info("Unmarshal: %v", err)
		return "FALSE"
	}

	openmcpkc_org := &cobrautil.KubeConfig{}
	copier.Copy(openmcpkc_org, openmcpkc)

	clusterName := ""
	for _, cluster := range openmcpkc.Clusters {
		if strings.Contains(cluster.Cluster.Server, memberIP) {
			clusterName = cluster.Name
			break
		}
	}

	if clusterName != "" {
		omcplog.V(4).Info("Already Join")
		return "FALSE"
	} else {
		//없으면 추가
		mem_context := memberkc.Contexts[0]
		mem_context.Name = clusterInstance.Name
		mem_context.Context.Cluster = clusterInstance.Name
		mem_context.Context.User = clusterInstance.Name

		mem_cluster := memberkc.Clusters[0]
		mem_cluster.Name = clusterInstance.Name
		mem_cluster.Cluster.Server = strings.ToLower(memberIP)

		mem_user := memberkc.Users[0]
		mem_user.Name = clusterInstance.Name

		/*if clusterInstance.Spec.ClusterPlatformType == "GKE" {
			mem_user.User.AuthProvider.Config.AccessToken = clusterInstance.Spec.GkeAccessToken
		}*/

		openmcpkc.Clusters = append(openmcpkc.Clusters, mem_cluster)
		openmcpkc.Contexts = append(openmcpkc.Contexts, mem_context)
		openmcpkc.Users = append(openmcpkc.Users, mem_user)

		cobrautil.WriteKubeConfig(openmcpkc, "/mnt/config")
		util.CmdExec2("cp /mnt/config $HOME/.kube/config")

		omcplog.V(4).Info("Ready to Join.")
		omcplog.V(4).Info("Join Start---")

		cluster_config, err_config := BuildConfigFromFlags(mem_cluster.Name, "/mnt/config")
		openmcp_config, err_oconfig := BuildConfigFromFlags("openmcp", "/mnt/config")

		if err_config != nil || err_oconfig != nil {
			omcplog.V(4).Info("err - ", err_config)
			omcplog.V(4).Info("err - ", err_oconfig)
			omcplog.V(2).Info("rollback KubeConfig")
			cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
			return "FALSE"
		} else {
			cluster_client := kubernetes.NewForConfigOrDie(cluster_config)
			openmcp_client := kubernetes.NewForConfigOrDie(openmcp_config)

			//1. CREATE namespace "kube-federation-system"
			Namespace := corev1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kube-federation-system",
				},
			}

			ns, err_ns := cluster_client.CoreV1().Namespaces().Create(context.TODO(), &Namespace, metav1.CreateOptions{})

			if err_ns != nil {
				omcplog.V(4).Info("Fail to Create Namespace Resource in " + mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_ns)
				var err_ns_get error
				ns, err_ns_get = cluster_client.CoreV1().Namespaces().Get(context.TODO(), "kube-federation-system", metav1.GetOptions{})

				if err_ns_get != nil {
					omcplog.V(4).Info("err_ns_get: ", ns)
					omcplog.V(2).Info("rollback KubeConfig")
					cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
					return "FALSE"
				} else {
					omcplog.V(4).Info("Get Namespace Resource [" + ns.Name + "] in " + mem_cluster.Name)
				}

			} else {
				omcplog.V(4).Info("[Step 1] Create Namespace Resource [" + ns.Name + "] in " + mem_cluster.Name)
			}

			//2. CREATE service account
			ServiceAccount := corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServiceAccount",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      mem_cluster.Name + "-openmcp",
					Namespace: "kube-federation-system",
				},
			}

			sa, err_sa := cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Create(context.TODO(), &ServiceAccount, metav1.CreateOptions{})

			if err_sa != nil {
				omcplog.V(4).Info("Fail to Create ServiceAccount Resource in " + mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_sa)
				var err_sa_get error
				sa, err_sa_get = cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Get(context.TODO(), mem_cluster.Name+"-openmcp", metav1.GetOptions{})

				if err_sa_get != nil {
					omcplog.V(4).Info("err_sa_get: ", err_sa_get)
					omcplog.V(2).Info("rollback KubeConfig")
					cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
					return "FALSE"
				} else {
					omcplog.V(4).Info("Get ServiceAccount Resource [" + sa.Name + "] in " + mem_cluster.Name)
				}

			} else {
				omcplog.V(4).Info("[Step 2] Create ServiceAccount Resource [" + sa.Name + "] in " + mem_cluster.Name)
			}

			//3. CREATE cluster role
			ClusterRole := rbacv1.ClusterRole{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kubefed-controller-manager:" + ServiceAccount.Name,
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{rbacv1.APIGroupAll},
						Verbs:     []string{rbacv1.VerbAll},
						Resources: []string{rbacv1.ResourceAll},
					},
					{
						NonResourceURLs: []string{rbacv1.NonResourceAll},
						Verbs:           []string{"get"},
					},
				},
			}

			cr, err_cr := cluster_client.RbacV1().ClusterRoles().Create(context.TODO(), &ClusterRole, metav1.CreateOptions{})

			if err_cr != nil {
				omcplog.V(4).Info("Fail to Create ClusterRole Resource in ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_cr)

				var err_cr_get error
				cr, err_cr_get = cluster_client.RbacV1().ClusterRoles().Get(context.TODO(), "kubefed-controller-manager:"+ServiceAccount.Name, metav1.GetOptions{})

				if err_cr_get != nil {
					omcplog.V(4).Info("err_cr_get: ", err_cr_get)
					omcplog.V(2).Info("rollback KubeConfig")
					cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
					return "FALSE"
				} else {
					omcplog.V(4).Info("Get ClusterRole Resource [" + cr.Name + "] in " + mem_cluster.Name)
				}

			} else {
				omcplog.V(4).Info("[Step 3] Create ClusterRole Resource [" + cr.Name + "] in " + mem_cluster.Name)
			}

			//4. CREATE cluster role binding
			ClusterRoleBinding := rbacv1.ClusterRoleBinding{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "kubefed-controller-manager:" + ServiceAccount.Name,
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "ClusterRole",
					Name:     ClusterRole.Name,
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      "ServiceAccount",
						Name:      ServiceAccount.Name,
						Namespace: ServiceAccount.Namespace,
					},
				},
			}

			crb, err_crb := cluster_client.RbacV1().ClusterRoleBindings().Create(context.TODO(), &ClusterRoleBinding, metav1.CreateOptions{})

			if err_crb != nil {
				omcplog.V(4).Info("Fail to Create ClusterRoleBinding Resource in ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_crb)

				var err_crb_get error
				crb, err_crb_get = cluster_client.RbacV1().ClusterRoleBindings().Get(context.TODO(), "kubefed-controller-manager:"+ServiceAccount.Name, metav1.GetOptions{})

				if err_crb_get != nil {
					omcplog.V(4).Info("err_crb_get: ", err_crb_get)
					omcplog.V(2).Info("rollback KubeConfig")
					cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
					return "FALSE"
				} else {
					omcplog.V(4).Info("Get ClusterRoleBinding Resource [" + crb.Name + "] in " + mem_cluster.Name)
				}

			} else {
				omcplog.V(4).Info("[Step 4] Create ClusterRoleBinding Resource [" + crb.Name + "] in " + mem_cluster.Name)
			}

			time.Sleep(1 * time.Second)

			//5. Get & CREATE secret (in openmcp)
			cluster_sa, err_sa1 := cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Get(context.TODO(), sa.Name, metav1.GetOptions{})
			if err_sa1 != nil {
				omcplog.V(4).Info("Fail to Get Secret Resource From ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_sa1)
				omcplog.V(2).Info("rollback KubeConfig")
				cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
				return "FALSE"
			}

			cluster_secret, err_sc := cluster_client.CoreV1().Secrets("kube-federation-system").Get(context.TODO(), cluster_sa.Secrets[0].Name, metav1.GetOptions{})
			if err_sc != nil {
				omcplog.V(4).Info("Fail to Get Secret Resource From ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_sc)
				omcplog.V(2).Info("rollback KubeConfig")
				cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
				return "FALSE"
			} else {
				omcplog.V(4).Info("[Step 5-1] Get Secret Resource [" + cluster_secret.Name + "] From " + mem_cluster.Name)
			}

			Secret := &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: mem_cluster.Name + "-",
					Namespace:    "kube-federation-system",
				},
				Data: map[string][]byte{
					"token": cluster_secret.Data["token"],
				},
			}

			secret_instance, err_secret := openmcp_client.CoreV1().Secrets("kube-federation-system").Create(context.TODO(), Secret, metav1.CreateOptions{})

			if err_secret != nil {
				omcplog.V(4).Info("Fail to Create secret Resource in openmcp")
				omcplog.V(4).Info("err: ", err_secret)

				var err_secret_get error
				secret_instance, err_secret_get = openmcp_client.CoreV1().Secrets("kube-federation-system").Get(context.TODO(), Secret.Name, metav1.GetOptions{})

				if err_secret_get != nil {
					omcplog.V(4).Info("err_secret_get: ", err_secret_get)
					omcplog.V(2).Info("rollback KubeConfig")
					cobrautil.WriteKubeConfig(openmcpkc_org, "/mnt/config")
					return "FALSE"
				} else {
					omcplog.V(4).Info("Get Secret Resource [" + secret_instance.Name + "] in openmcp")
				}

			} else {
				omcplog.V(4).Info("[Step 5-2] Create Secret Resource [" + secret_instance.Name + "] in openmcp")
			}

			//6. CREATE kubefedcluster (in openmcp)
			var disabledTLSValidations []fedv1b1.TLSValidation

			if cm.Host_config.TLSClientConfig.Insecure {
				disabledTLSValidations = append(disabledTLSValidations, fedv1b1.TLSAll)
			}

			KubefedCluster := &fedv1b1.KubeFedCluster{
				TypeMeta: metav1.TypeMeta{
					Kind:       "KubeFedCluster",
					APIVersion: "core.kubefed.io/v1beta1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      mem_cluster.Name,
					Namespace: "kube-federation-system",
				},
				Spec: fedv1b1.KubeFedClusterSpec{
					APIEndpoint: mem_cluster.Cluster.Server,
					CABundle:    cluster_secret.Data["ca.crt"],
					SecretRef: fedv1b1.LocalSecretReference{
						Name: secret_instance.Name,
					},
					DisabledTLSValidations: disabledTLSValidations,
				},
			}

			err_kubefed := r.live.Create(context.TODO(), KubefedCluster)

			if err_kubefed != nil {
				omcplog.V(4).Info("Fail to Create KubefedCluster Resource in openmcp")
				omcplog.V(4).Info("err: ", err_kubefed)

			} else {
				omcplog.V(4).Info("[Step 6] Create KubefedCluster Resource [" + KubefedCluster.Name + "] in openmcp")
			}
		}
		return "TRUE"
	}
}

func UnjoinAndDeleteConfig(memberkc *cobrautil.KubeConfig, openmcpkc *cobrautil.KubeConfig) {
	memberIP := memberkc.Clusters[0].Cluster.Server
	target_name := ""
	target_user := ""

	var target_name_index int
	var target_context_index int
	var target_user_index int

	for i, cluster := range openmcpkc.Clusters {
		lower_memberIP := strings.ToLower(memberIP)
		if strings.Contains(cluster.Cluster.Server, lower_memberIP) {
			target_name = cluster.Name
			target_name_index = i
			break
		}
	}
	for j, context := range openmcpkc.Contexts {
		if target_name == context.Context.Cluster {
			target_user = context.Context.User
			target_context_index = j
			break
		}
	}
	for k, user := range openmcpkc.Users {
		if target_user == user.Name {
			target_user_index = k
			break
		}
	}

	if target_name != "" {

		cluster_config, _ := BuildConfigFromFlags(target_name, "/mnt/config")
		cluster_client := kubernetes.NewForConfigOrDie(cluster_config)

		//1. DELETE cluster role binding / cluster role / namespace
		err_deletecrb := cluster_client.RbacV1().ClusterRoleBindings().Delete(context.TODO(), "kubefed-controller-manager:"+target_name+"-openmcp", metav1.DeleteOptions{})
		err_deletecr := cluster_client.RbacV1().ClusterRoles().Delete(context.TODO(), "kubefed-controller-manager:"+target_name+"-openmcp", metav1.DeleteOptions{})
		err_deletens := cluster_client.CoreV1().Namespaces().Delete(context.TODO(), "kube-federation-system", metav1.DeleteOptions{})

		if err_deletecrb == nil && err_deletecr == nil && err_deletens == nil {
			omcplog.V(4).Info("[Step 1] DELETE CR/CRB/NS Resource in ", target_name)
		} else {
			omcplog.V(4).Info("Fail to DELETE CR/CRB/NS Resource in ", target_name)
			omcplog.V(4).Info("err_deletecrb: ", err_deletecrb)
			omcplog.V(4).Info("err_deletecr: ", err_deletecr)
			omcplog.V(4).Info("err_deletens: ", err_deletens)
		}

		kfc_instance := &fedv1b1.KubeFedCluster{}
		err := r.live.Get(context.TODO(), types.NamespacedName{Name: target_name, Namespace: "kube-federation-system"}, kfc_instance)

		if err == nil {
			//2. DELETE secret (in openmcp)
			sec_instance := &corev1.Secret{}
			err_isec := r.live.Get(context.TODO(), types.NamespacedName{Name: kfc_instance.Spec.SecretRef.Name, Namespace: "kube-federation-system"}, sec_instance)

			if err_isec == nil {
				err_deletesec := r.live.Delete(context.TODO(), sec_instance)

				if err_deletesec != nil {
					omcplog.V(4).Info("Fail to DELETE Secret Resource in openmcp")
					omcplog.V(4).Info("err: ", err_deletesec)
				} else {
					omcplog.V(4).Info("[Step 2] DELETE Secret Resource [" + sec_instance.Name + "] in openmcp")
				}
			}

			//3. DELETE kubefedcluster (in openmcp)
			err_kubefed := r.live.Delete(context.TODO(), kfc_instance)

			if err_kubefed != nil {
				omcplog.V(4).Info("Fail to DELETE KubefedCluster Resource in openmcp")
				omcplog.V(4).Info("err: ", err_kubefed)
			} else {
				omcplog.V(4).Info("[Step 3] DELETE KubefedCluster Resource [" + kfc_instance.Name + "] in openmcp")
			}
		} else {
			omcplog.V(4).Info(err)
		}

		openmcpkc.Clusters = append(openmcpkc.Clusters[:target_name_index], openmcpkc.Clusters[target_name_index+1:]...)
		openmcpkc.Contexts = append(openmcpkc.Contexts[:target_context_index], openmcpkc.Contexts[target_context_index+1:]...)
		openmcpkc.Users = append(openmcpkc.Users[:target_user_index], openmcpkc.Users[target_user_index+1:]...)

		cobrautil.WriteKubeConfig(openmcpkc, "/mnt/config")

		omcplog.V(4).Info("Complete to Delete " + target_name + " Info")
	} else {
		omcplog.V(4).Info("Fail to Delete " + target_name + " Info")
	}

}
