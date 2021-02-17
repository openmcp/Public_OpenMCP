package openmcpcluster

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"openmcp/openmcp/apis"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util"
	"openmcp/openmcp/util/clusterManager"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	"strings"
	"time"
)

var cm *clusterManager.ClusterManager
var log = logf.Log.WithName("controller_openmcpcluster")
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
	if err := co.WatchResourceReconcileObject(live, &clusterv1alpha1.OpenMCPCluster{}, controller.WatchOptions{}); err != nil {
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
	if clusterInstance.Spec.ClusterStatus == "STANDBY" {
		omcplog.V(4).Info(clusterInstance.Name + " [ STANDBY ]")

	} else if clusterInstance.Spec.ClusterStatus == "JOIN" {
		omcplog.V(4).Info(clusterInstance.Name + " [ JOIN ]")
		joinCheck := MergeConfigAndJoin(*clusterInstance)

		if joinCheck == "TRUE" {
			omcplog.V(4).Info("OpenMCP Module Deploy---")
			moduleDirectory := []string{"namespace", "custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller", "configmap"}
			for i, dirname := range moduleDirectory {
				moduleDirectory[i] = "/init/" + dirname
			}
			util.CmdExec2("cp /mnt/config $HOME/.kube/config")
			InstallInitModule(moduleDirectory, clusterInstance.Name)
			omcplog.V(4).Info("--- JOIN Complete ---")
		}

	} else if clusterInstance.Spec.ClusterStatus == "UNJOIN" {
		omcplog.V(4).Info(clusterInstance.Name + " [ UNJOIN ]")

		//config 파일 확인 (클러스터 조인 유무)
		memberkc := &cobrautil.KubeConfig{}
		err := yaml.Unmarshal(clusterInstance.Spec.ClusterInfo, memberkc)
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
			if strings.Contains(cluster.Cluster.Server, memberIP) {
				unjoinCheck = cluster.Name
				break
			}
		}

		if unjoinCheck != "" {
			omcplog.V(4).Info("OpenMCP Module Delete---")
			moduleDirectory := []string{"custom-metrics-apiserver", "metallb", "metric-collector", "metrics-server", "nginx-ingress-controller", "configmap", "namespace"}
			for i, dirname := range moduleDirectory {
				moduleDirectory[i] = "/init/" + dirname
			}
			util.CmdExec2("cp /mnt/config $HOME/.kube/config")
			UninstallInitModule(moduleDirectory, clusterInstance.Name)
			UnjoinAndDeleteConfig(*clusterInstance, memberkc, openmcpkc)

			omcplog.V(4).Info("--- UNJOIN Complete ---")
		} else {
			omcplog.V(4).Info("Not Exists Cluster Info")
		}
	}

	return reconcile.Result{}, nil
}

func InstallInitModule(directory []string, clustername string) {

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
					InstallInitModule([]string{dirname + "/" + f.Name()}, clustername)
				} else {
					if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
						util.CmdExec2("/usr/local/bin/kubectl apply -f " + dirname + "/" + f.Name() + " --context " + clustername)
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
					InstallInitModule([]string{dirname + "/" + f.Name()}, clustername)
				} else {
					if filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".yml" {
						util.CmdExec2("/usr/local/bin/kubectl delete -f " + dirname + "/" + f.Name() + " --context " + clustername)
					}
				}
			}
		}
	}
}

func MergeConfigAndJoin(clusterInstance clusterv1alpha1.OpenMCPCluster) string {
	//config파일에 해당 정보가 저장되어 있는지 확인
	memberkc := &cobrautil.KubeConfig{}
	err := yaml.Unmarshal(clusterInstance.Spec.ClusterInfo, memberkc)
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
		mem_cluster := memberkc.Clusters[0]
		mem_user := memberkc.Users[0]

		openmcpkc.Clusters = append(openmcpkc.Clusters, mem_cluster)
		openmcpkc.Contexts = append(openmcpkc.Contexts, mem_context)
		openmcpkc.Users = append(openmcpkc.Users, mem_user)

		cobrautil.WriteKubeConfig(openmcpkc, "/mnt/config")

		omcplog.V(4).Info("Ready to Join.")
		omcplog.V(4).Info("Join Start---")

		cluster_config, err_config := BuildConfigFromFlags(mem_cluster.Name, "/mnt/config")
		openmcp_config, err_oconfig := BuildConfigFromFlags("openmcp", "/mnt/config")

		if err_config != nil || err_oconfig != nil {
			omcplog.V(4).Info("err - ", err_config)
			omcplog.V(4).Info("err - ", err_oconfig)
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

			ns, err_ns := cluster_client.CoreV1().Namespaces().Create(&Namespace)

			if err_ns != nil {
				omcplog.V(4).Info("Fail to Create Namespace Resource in " + mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_ns)
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

			sa, err_sa := cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Create(&ServiceAccount)

			if err_sa != nil {
				omcplog.V(4).Info("Fail to Create ServiceAccount Resource in " + mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_sa)
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

			cr, err_cr := cluster_client.RbacV1().ClusterRoles().Create(&ClusterRole)

			if err_cr != nil {
				omcplog.V(4).Info("Fail to Create ClusterRole Resource in ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_cr)
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

			crb, err_crb := cluster_client.RbacV1().ClusterRoleBindings().Create(&ClusterRoleBinding)

			if err_crb != nil {
				omcplog.V(4).Info("Fail to Create ClusterRoleBinding Resource in ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_crb)
			} else {
				omcplog.V(4).Info("[Step 4] Create ClusterRoleBinding Resource [" + crb.Name + "] in " + mem_cluster.Name)
			}

			time.Sleep(1 * time.Second)

			//5. Get & CREATE secret (in openmcp)
			cluster_sa, err_sa1 := cluster_client.CoreV1().ServiceAccounts("kube-federation-system").Get(sa.Name, metav1.GetOptions{})

			cluster_secret, err_sc := cluster_client.CoreV1().Secrets("kube-federation-system").Get(cluster_sa.Secrets[0].Name, metav1.GetOptions{})

			if err_sc != nil || err_sa1 != nil {
				omcplog.V(4).Info("Fail to Get Secret Resource From ", mem_cluster.Name)
				omcplog.V(4).Info("err: ", err_sc)
				omcplog.V(4).Info("err: ", err_sa1)
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

			secret_instance, err_secret := openmcp_client.CoreV1().Secrets("kube-federation-system").Create(Secret)

			if err_secret != nil {
				omcplog.V(4).Info("Fail to Create secret Resource in openmcp")
				omcplog.V(4).Info("err: ", err_secret)
			} else {
				omcplog.V(4).Info("[Step 5-2] Create Secret Resource [" + mem_cluster.Name + "] in openmcp")
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

func UnjoinAndDeleteConfig(clusterInstance clusterv1alpha1.OpenMCPCluster, memberkc *cobrautil.KubeConfig, openmcpkc *cobrautil.KubeConfig) {
	memberIP := memberkc.Clusters[0].Cluster.Server

	target_name := ""
	target_user := ""

	var target_name_index int
	var target_context_index int
	var target_user_index int

	for i, cluster := range openmcpkc.Clusters {
		if strings.Contains(cluster.Cluster.Server, memberIP) {
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

	mem_cluster := memberkc.Clusters[0]
	cluster_config, _ := BuildConfigFromFlags(mem_cluster.Name, "/mnt/config")
	cluster_client := kubernetes.NewForConfigOrDie(cluster_config)

	//1. DELETE cluster role binding / cluster role / namespace
	err_deletecrb := cluster_client.RbacV1().ClusterRoleBindings().Delete("kubefed-controller-manager:"+target_name+"-openmcp", &metav1.DeleteOptions{})
	err_deletecr := cluster_client.RbacV1().ClusterRoles().Delete("kubefed-controller-manager:"+target_name+"-openmcp", &metav1.DeleteOptions{})
	err_deletens := cluster_client.CoreV1().Namespaces().Delete("kube-federation-system", &metav1.DeleteOptions{})

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

	omcplog.V(4).Info("Complete to Delete" + target_name + " Info")

}
