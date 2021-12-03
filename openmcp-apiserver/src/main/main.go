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
	"log"
	"net/http"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-apiserver/src/auth"
	"openmcp/openmcp/openmcp-apiserver/src/httphandler"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
)

var live *cluster.Cluster

func main() {
	logLevel.KetiLogInit()

	cm := clusterManager.NewClusterManager()

	host_ctx := "openmcp"
	namespace := "openmcp"

	host_cfg := cm.Host_config
	live = cluster.New(host_ctx, host_cfg, cluster.Options{})

	httphandler.Live = live

	HTTPServer_PORT := "8080"

	httpManager := &httphandler.HttpManager{
		HTTPServer_PORT: HTTPServer_PORT,
		ClusterManager:  cm,
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/token", auth.TokenHandler)
	handler.Handle("/", auth.AuthMiddleware(http.HandlerFunc(httpManager.RouteHandler)))
	handler.HandleFunc("/join", httphandler.JoinHandler)
	handler.HandleFunc("/joinCloud", httphandler.JoinCloudHandler)

	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}

	go func() {
		omcplog.V(2).Info("Start OpenMCP API Server")
		//err := server.ListenAndServe()
		//err := server.ListenAndServeTLS("/tmp/cert/10.0.3.40/server.crt", "/tmp/cert/10.0.3.40/server.key")
		err := server.ListenAndServeTLS("/tmp/cert/server.crt", "/tmp/cert/server.key")
		if err != nil {
			omcplog.V(0).Info(err)
		}
	}()

	ghosts := []*cluster.Cluster{}

	for _, ghost_cluster := range cm.Cluster_list.Items {
		ghost_ctx := ghost_cluster.Name
		ghost_cfg := cm.Cluster_configs[ghost_ctx]

		ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
		ghosts = append(ghosts, ghost)
	}

	reshape_cont, err_reshape := reshape.NewController(live, ghosts, namespace, cm)
	if err_reshape != nil {
		omcplog.V(2).Info("err_reshape : ", err_reshape)
		return
	}

	loglevel_cont, err_log := logLevel.NewController(live, ghosts, namespace)
	if err_log != nil {
		omcplog.V(2).Info("err_log : ", err_log)
		return
	}

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
