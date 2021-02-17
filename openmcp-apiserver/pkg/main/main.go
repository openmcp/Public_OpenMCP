/*
Copyright 2018 The Multicluster-Controller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
	"context"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-apiserver/pkg/auth"
	"openmcp/openmcp/openmcp-apiserver/pkg/httphandler"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"gopkg.in/yaml.v2"
	cobrautil "openmcp/openmcp/omcpctl/util"
)

func CreateClusterResource(name string, config []byte) (string, error) {

	clusterCR := &clusterv1alpha1.OpenMCPCluster{
		TypeMeta: v1.TypeMeta{
			Kind:       "OpenMCPCluster",
			APIVersion: "apiextensions.k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: "openmcp",
		},
		Spec: clusterv1alpha1.OpenMCPClusterSpec{
			ClusterStatus: "STANDBY",
			ClusterInfo:   config,
		},
		//Status: clusterv1alpha1.OpenMCPClusterStatus{ClusterStatus: "STANDBY"},
	}

	fmt.Println("clusterCR : ", clusterCR.Spec.ClusterStatus)

	liveClient, _ := live.GetDelegatingClient()

	err := liveClient.Create(context.TODO(), clusterCR)
	//err_status := liveClient.Status().Update(context.TODO(), clusterCR)

	if err != nil { //|| err_status != nil {
		//fmt.Println("err : ", err)
		//fmt.Println("err_status : ", err_status)
		omcplog.V(4).Info("Fail to create openmcpcluster resource")
	} else {
		omcplog.V(4).Info("Success to create openmcpcluster resource")
	}

	//fmt.Println(string(clusterCR.Spec.ClusterInfo))

	return clusterCR.Name, err
}

func join(w http.ResponseWriter, r *http.Request) {

	file, _, err := r.FormFile("file")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)

	c := &cobrautil.KubeConfig{}
	err = yaml.Unmarshal(fileBytes, c)

	CreateClusterResource(c.Clusters[0].Name, fileBytes)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("INVALID_FILE"))
		return
	}

	a := []byte("OK\n")
	w.Write(a)
}

var live *cluster.Cluster

func main() {
	logLevel.KetiLogInit()

	for {
		cm := clusterManager.NewClusterManager()

		HTTPServer_PORT := "8080"

		httpManager := &httphandler.HttpManager{
			HTTPServer_PORT: HTTPServer_PORT,
			ClusterManager:  cm,
		}

		handler := http.NewServeMux()

		handler.HandleFunc("/token", auth.TokenHandler)
		handler.Handle("/", auth.AuthMiddleware(http.HandlerFunc(httpManager.RouteHandler)))
		handler.HandleFunc("/join", join)

		server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}

		go func() {
			omcplog.V(2).Info("Start OpenMCP API Server")
			err := server.ListenAndServe()
			//err := server.ListenAndServeTLS("ca.crt","ca.key")
			if err != nil {
				omcplog.V(0).Info(err)
			}
		}()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		live = cluster.New(host_ctx, host_cfg, cluster.Options{})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
			ghosts = append(ghosts, ghost)
		}

		reshape_cont, _ := reshape.NewController(live, ghosts, namespace)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

		m := manager.New()
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}

		/*if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("OpenMCP API Server Shutdown Failed:%+v", err)
		}*/

		omcplog.V(2).Info("OpenMCP API Server Exited Properly")
	}

}
