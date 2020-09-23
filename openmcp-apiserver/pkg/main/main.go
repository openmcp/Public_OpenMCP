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
	"context"
	"openmcp/openmcp/openmcp-apiserver/pkg/auth"
	"openmcp/openmcp/openmcp-apiserver/pkg/httphandler"

	"net/http"
	"openmcp/openmcp/omcplog"

	"openmcp/openmcp/util/clusterManager"

	"log"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"

	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
)


func main() {
	logLevel.KetiLogInit()


	for {
		cm := clusterManager.NewClusterManager()

		//HTTPServer_IP := "10.0.3.20"
		HTTPServer_PORT := "8080"

		httpManager := &httphandler.HttpManager{
			//HTTPServer_IP: HTTPServer_IP,
			HTTPServer_PORT: HTTPServer_PORT,
			ClusterManager:  cm,
		}

		handler := http.NewServeMux()

		//handler.HandleFunc("/token", TokenHandler)
		//handler.Handle("/", AuthMiddleware(http.HandlerFunc(httpManager.ExampleHandler)))
		handler.HandleFunc("/token", auth.TokenHandler)
		handler.Handle("/", auth.AuthMiddleware(http.HandlerFunc(httpManager.RouteHandler)))
		//handler.HandleFunc("/metrics/{name:[a-z]+}", httpManager.MetricsHandler)

		//handler.HandleFunc("/omcpexec", httpManager.ExampleHandler2)

		server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}

		go func() {
			omcplog.V(2).Info("Start OpenMCP API Server")
			err := server.ListenAndServe()
			if err != nil {
				omcplog.V(0).Info(err)
			}
		}()


		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config
		//live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
		live := cluster.New(host_ctx, host_cfg, cluster.Options{})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			//ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
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

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("OpenMCP API Server Shutdown Failed:%+v", err)
		}

		log.Print("OpenMCP API Server Exited Properly")
	}

}

//fmt.Println("Connect Etcd Main")
//fmt.Println("-----------------------------")
//fmt.Println("Host : ", r.Host)
//fmt.Println("URL : ", r.URL)
//fmt.Println("URL.Host : ", r.URL.Host)
//fmt.Println("URL.Path : ", r.URL.Path)
//fmt.Println("URL.ForceQuery : ", r.URL.ForceQuery)
//fmt.Println("URL.Fragment : ", r.URL.Fragment)
//fmt.Println("URL.Opaque : ", r.URL.Opaque)
//fmt.Println("URL.RawPath : ", r.URL.RawPath)
//fmt.Println("URL.RawQuery : ", r.URL.RawQuery)
//fmt.Println("URL.Scheme : ", r.URL.Scheme)
//fmt.Println("URL.User : ", r.URL.User)
//fmt.Println("RequestURI : ", r.RequestURI)
//fmt.Println("Method : ", r.Method)
//fmt.Println("RemoteAddr : ", r.RemoteAddr)
//fmt.Println("Proto : ", r.Proto)
//fmt.Println("Header : ", r.Header)
//client kubernetes.Interface




// GET http://10.0.3.20:31635/token?username=openmcp&password=keti
// Get the Token
// Add Header
// --> Key : Authorization
// --> Value : Bearer {TOKEN}
// GET http://10.0.3.20:31635/api?clustername=openmcp
