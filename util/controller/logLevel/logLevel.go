package logLevel

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/controller"
	"admiralty.io/multicluster-controller/pkg/reconcile"
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	//"k8s.io/klog"

	//"k8s.io/klog"
	"openmcp/openmcp/omcplog"

	"openmcp/openmcp/openmcp-resource-controller/apis"
	ketiv1alpha1 "openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
)

func NewController(live *cluster.Cluster, ghosts []*cluster.Cluster, ghostNamespace string) (*controller.Controller, error) {
	liveclient, err := live.GetDelegatingClient()
	if err != nil {
		return nil, fmt.Errorf("getting delegating client for live cluster: %v", err)
	}
	ghostclients := []client.Client{}
	for _, ghost := range ghosts {
		ghostclient, err := ghost.GetDelegatingClient()
		if err != nil {
			return nil, fmt.Errorf("getting delegating client for ghost cluster: %v", err)
		}
		ghostclients = append(ghostclients, ghostclient)
	}

	co := controller.New(&reconciler{live: liveclient, ghosts: ghostclients, ghostNamespace: ghostNamespace}, controller.Options{})
	if err := apis.AddToScheme(live.GetScheme()); err != nil {
		return nil, fmt.Errorf("adding APIs to live cluster's scheme: %v", err)
	}

	fmt.Printf("%T, %s\n", live, live.GetClusterName())
	if err := co.WatchResourceReconcileObject(live, &ketiv1alpha1.OpenMCPPolicy{}, controller.WatchOptions{}); err != nil {
		return nil, fmt.Errorf("setting up Pod watch in live cluster: %v", err)
	}
	return co, nil
}

type reconciler struct {
	live           client.Client
	ghosts         []client.Client
	ghostNamespace string
}

var logLevel = "0"

func (r *reconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {

	if req.Namespace == "openmcp" && req.Name == "log-version" {
		prevLogLevel := logLevel
		logLevel = r.getLogLevel()
		if prevLogLevel != logLevel {
			omcplog.Info("LogLevel Changed, Used LogLevel (" + prevLogLevel + " -> " + logLevel + ")")
			flag.Set("omcpv", logLevel)
			flag.Parse()
		}
	}
	return reconcile.Result{}, nil // err
}

func (r *reconciler) getLogLevel() string {
	instance := &ketiv1alpha1.OpenMCPPolicy{}
	nsn := types.NamespacedName{
		Namespace: "openmcp",
		Name:      "log-version",
	}
	err := r.live.Get(context.TODO(), nsn, instance)
	if err != nil && errors.IsNotFound(err) {
		omcplog.Info("Not Exist Policy 'log-version', Use Default LogLevel (0)")
		return "0"
	} else if err != nil {
		omcplog.Info("FatalError ! ", err)
	}
	if instance.Spec.PolicyStatus == "Enabled" {
		for _, policy := range instance.Spec.Template.Spec.Policies {
			//if policy.Type == "Version" && len(policy.Value) == 1 {
			//	matched, _ := regexp.MatchString("[0-9]", policy.Value[0])
			//	if matched && len(policy.Value[0]) == 1 {
			//		return policy.Value[0]
			//	}
			//	klog.Info("Policy 'log-version' Value must be [0-9], Use Default LogLevel (0)")
			//	return "0"
			//}
			if policy.Type == "Version" && len(policy.Value) == 1 {
				logLevelString := policy.Value[0]
				logLevelInt, err := strconv.Atoi(logLevelString)
				if err == nil && logLevelInt >= -1 && logLevelInt <= 9 {
					return logLevelString
				}
				omcplog.Info("Policy 'log-version' Value must be [-1~9], Use Default LogLevel (0)")
				return "0"
			}
		}
		omcplog.Info("Policy 'log-version' Format Error, Use Default LogLevel (0)")
		return "0"

	} else {
		omcplog.Info("Policy 'log-version' Disabled, Use Default LogLevel (0)")
		return "0"
	}

}

func KetiLogInit() {
	omcplog.InitFlags(nil)
	flag.Set("omcpv", "0")
	flag.Parse()
}
