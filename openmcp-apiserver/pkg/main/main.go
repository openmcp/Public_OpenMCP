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
	"bytes"
	"context"

	"crypto/tls"
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	restclient "k8s.io/client-go/rest"

	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"openmcp/openmcp/omcplog"

	"openmcp/openmcp/util/clusterManager"

	"strings"
	"time"

	"log"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"

	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
)

//func APIServer(cm *clusterManager.ClusterManager) {
//	//HTTPServer_IP := "10.0.3.20"
//	HTTPServer_PORT := "8080"
//
//	httpManager := &HttpManager{
//		//HTTPServer_IP: HTTPServer_IP,
//		HTTPServer_PORT: HTTPServer_PORT,
//		ClusterManager:  cm,
//	}
//
//	handler := http.NewServeMux()
//
//	//handler.HandleFunc("/token", TokenHandler)
//	//handler.Handle("/", AuthMiddleware(http.HandlerFunc(httpManager.ExampleHandler)))
//	handler.HandleFunc("/", httpManager.ExampleHandler)
//
//	//handler.HandleFunc("/omcpexec", httpManager.ExampleHandler2)
//
//	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}
//
//	omcplog.V(2).Info("Start OpenMCP API Server")
//	err := server.ListenAndServe()
//	if err != nil {
//		omcplog.V(0).Info(err)
//	}
//}
func main() {
	logLevel.KetiLogInit()


	for {
		cm := clusterManager.NewClusterManager()

		//HTTPServer_IP := "10.0.3.20"
		HTTPServer_PORT := "8080"

		httpManager := &HttpManager{
			//HTTPServer_IP: HTTPServer_IP,
			HTTPServer_PORT: HTTPServer_PORT,
			ClusterManager:  cm,
		}

		handler := http.NewServeMux()

		//handler.HandleFunc("/token", TokenHandler)
		//handler.Handle("/", AuthMiddleware(http.HandlerFunc(httpManager.ExampleHandler)))
		handler.HandleFunc("/", httpManager.ExampleHandler)

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



const (
	APP_KEY = "openmcp-apiserver"
)

/*func(h *HttpManager) SendOpenMCPAPIServer(ctx context.Context, r *protobuf.RequestInfo) (*protobuf.ResponseInfo, error) {

	clusterName := r.ClusterName
	fmt.Println(clusterName)
	if len(clusterName) < 1 {
		message := "URL Param 'clustername' is missing"
		return &protobuf.ResponseInfo{Message: message}, nil
	}

	APISERVER := ""
	TOKEN := ""

	if clusterName == "openmcp" {
		APISERVER = h.ClusterManager.Host_config.Host
		TOKEN = h.ClusterManager.Host_config.BearerToken
	} else {
		for _, cluster := range h.ClusterManager.Cluster_list.Items {
			if cluster.Name == clusterName {
				APISERVER = cluster.Spec.APIEndpoint
				TOKEN = h.ClusterManager.Cluster_configs[cluster.Name].BearerToken
			}
		}

	}

	//To k8s api server
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	var req *http.Request
	var err error

	if r.Method == "GET" || r.Method == "DELETE"{
		req, err = http.NewRequest(r.Method, APISERVER+r.Path, nil)
	}else if r.Method == "POST" || r.Method == "PUT"{
		req, err = http.NewRequest(r.Method, APISERVER+r.Path, bytes.NewBufferString(r.Body))
	}

	if err != nil {
		fmt.Println("Check1", err)
		// handle err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Check2", err)
		// handle err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Check3", err)
		panic(err.Error())
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		fmt.Println("Check4", err)
		panic(err.Error())
	}

	fmt.Printf("%s\n", prettyJSON.Bytes())
	//w.Write(prettyJSON.Bytes())

	return &protobuf.ResponseInfo{Message: string(prettyJSON.Bytes()), Status: "OK"}, nil

}*/

/*func (h *HttpManager) ExampleHandler2(w http.ResponseWriter, r *http.Request) {

	fmt.Println("exec ", r.URL)
	fmt.Println("exec ", r.URL.Path)
}*/

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

func ExecCmdExample(restclient *restclient.RESTClient, config *restclient.Config, podName string, podNamespace string,
	command []string, containerName string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	omcplog.V(5).Info(restclient)
	omcplog.V(5).Info(config)
	req := restclient.Post().Resource("pods").Name(podName).
		Namespace(podNamespace).SubResource("exec")
	option := &corev1.PodExecOptions{
		Container: containerName,
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}
	if stdin == nil {
		option.Stdin = false
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *HttpManager) ExampleHandler(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.Path, "exec"){
		clusterNames, ok := r.URL.Query()["clustername"]

		fmt.Println(clusterNames, ok)
		if !ok || len(clusterNames[0]) < 1 {
			w.Write([]byte("Url Param 'clustername' is missing"))
			return
		}
		fmt.Println()

	/*	APISERVER := ""
		TOKEN := ""
		clusterName := clusterNames[0]

		if clusterName == "openmcp" {
			APISERVER = h.ClusterManager.Host_config.Host
			TOKEN = h.ClusterManager.Host_config.BearerToken
		} else {
			for _, cluster := range h.ClusterManager.Cluster_list.Items {
				if cluster.Name == clusterName {
					APISERVER = cluster.Spec.APIEndpoint
					TOKEN = h.ClusterManager.Cluster_configs[cluster.Name].BearerToken
				}
			}
		}*/

		omcplog.V(3).Info("host ", r.Host)
		omcplog.V(3).Info("url ", r.URL)
		omcplog.V(3).Info("url/host ", r.URL.Host)
		omcplog.V(3).Info("body ", r.Body)
		omcplog.V(3).Info(r.URL.Query()["containername"])
		omcplog.V(3).Info(r.URL.Query()["clustername"])
		omcplog.V(3).Info(r.URL.Query()["podname"])
		omcplog.V(3).Info(r.URL.Query()["podnamespace"])
		omcplog.V(3).Info(r.URL.Query()["stdin"])
		omcplog.V(3).Info(r.URL.Query()["tty"])
		omcplog.V(3).Info(r.URL.Query()["stdout"])
		omcplog.V(3).Info(r.URL.Query()["stderr"])

		stdin := "false"
		stdout := "true"
		stderr := "true"
		//tty := "false"

		//containername := r.URL.Query()["containername"][0]
		//clustername := r.URL.Query()["clustername"][0]
		podname := r.URL.Query()["podname"][0]
		podnamespace := r.URL.Query()["podnamespace"][0]
		stdin = r.URL.Query()["stdin"][0]
		//tty = r.URL.Query()["tty"][0]
		stdout = r.URL.Query()["stdout"][0]
		stderr = r.URL.Query()["stderr"][0]

		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)

		omcplog.V(5).Info(bodyString)

		stdin_io := bytes.NewBufferString(stdin)
		stdout_io := bytes.NewBufferString(stdout)
		stderr_io := bytes.NewBufferString(stderr)
		//tty_io := bytes.NewBufferString(tty)

		var a []string

		a = strings.Split(bodyString, ",")

		omcplog.V(5).Info(a)
		omcplog.V(5).Info(a[0])
		//fmt.Println(a[1])

		restClient, _ := restclient.RESTClientFor(h.ClusterManager.Host_config)

		ExecCmdExample(restClient,h.ClusterManager.Host_config ,podname, podnamespace, a, "", stdin_io, stdout_io, stderr_io)

	}else {

		clusterNames, ok := r.URL.Query()["clustername"]

		omcplog.V(5).Info(clusterNames, ok)
		if !ok || len(clusterNames[0]) < 1 {
			w.Write([]byte("Url Param 'clustername' is missing"))
			return
		}

		APISERVER := ""
		TOKEN := ""
		clusterName := clusterNames[0]

		if clusterName == "openmcp" {
			APISERVER = h.ClusterManager.Host_config.Host
			TOKEN = h.ClusterManager.Host_config.BearerToken
		} else {
			for _, cluster := range h.ClusterManager.Cluster_list.Items {
				if cluster.Name == clusterName {
					APISERVER = cluster.Spec.APIEndpoint
					TOKEN = h.ClusterManager.Cluster_configs[cluster.Name].BearerToken
				}
			}
		}
		omcplog.V(3).Info("APISERVER : " + APISERVER)
		omcplog.V(3).Info("clusterName : " + clusterName)
		omcplog.V(3).Info("Method : " + r.Method)
		omcplog.V(3).Info("URLPath : " + r.URL.Path)

		// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
		// TODO: This is insecure; use only in dev environments.

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		var req *http.Request
		var err error

		if r.Method == "GET" || r.Method == "DELETE" {
			req, err = http.NewRequest(r.Method, APISERVER+r.URL.Path, nil)
		} else if r.Method == "POST" || r.Method == "PUT" {
			req, err = http.NewRequest(r.Method, APISERVER+r.URL.Path, r.Body)
		}

		if err != nil {
			omcplog.V(0).Info(err)
			// handle err
		}

		omcplog.V(5).Info("Content-Type : ", r.Header.Get("Content-Type"))
		omcplog.V(5).Info("Authorization", "Bearer "+TOKEN)

		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		req.Header.Set("Authorization", "Bearer "+TOKEN)


		resp, err := client.Do(req)
		omcplog.V(3).Info("Request Done!")
		if err != nil {
			omcplog.V(0).Info(err)
			// handle err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			omcplog.V(0).Info(err)
			panic(err.Error())
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			omcplog.V(0).Info(err)
			panic(err.Error())
		}

		omcplog.V(5).Info(string(prettyJSON.Bytes()))

		w.Write(body)
		//omcplog.V(5).Info(string(body))

		//w.Write(prettyJSON.Bytes())
	}

}

//// TokenHandler is our handler to take a username and password and,
//// if it's valid, return a token used for future requests.
func TokenHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()

	// Check the credentials provided - if you store these in a database then
	// this is where your query would go to check.
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username != "openmcp" || password != "keti" {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error":"invalid_credentials"}`)
		return
	}

	// We are happy with the credentials, so build a token. We've given it
	// an expiry of 1 hour.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": username,
		"exp":  time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat":  time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(APP_KEY))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"token_generation_failed"}`)
		return
	}
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

// AuthMiddleware is our middleware to check our token is valid. Returning
// a 401 status to the client if it is not valid.
func AuthMiddleware(next http.Handler) http.Handler {
	if len(APP_KEY) == 0 {
		log.Fatal("HTTP server unable to start, expected an APP_KEY for JWT auth")
	}
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(APP_KEY), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

type HttpManager struct {
	HTTPServer_IP   string
	HTTPServer_PORT string
	ClusterManager  *clusterManager.ClusterManager
}
const (
	GRPC_PORT = "8080"

)

//func main() {
//	//HTTPServer_IP := "10.0.3.20"
//	HTTPServer_PORT := "8080"
//
//	cm := NewClusterManager()
//
//	httpManager := &HttpManager{
//		//HTTPServer_IP: HTTPServer_IP,
//		HTTPServer_PORT: HTTPServer_PORT,
//		ClusterManager:  cm,
//	}
//
//	handler := http.NewServeMux()
//
//	//handler.HandleFunc("/token", TokenHandler)
//	//handler.Handle("/", AuthMiddleware(http.HandlerFunc(httpManager.`ExampleHandler`)))
//	handler.HandleFunc("/", httpManager.ExampleHandler)
//
//	//handler.HandleFunc("/omcpexec", httpManager.ExampleHandler2)
//
//	server := &http.Server{Addr: ":" + HTTPServer_PORT, Handler: handler}
//
//	fmt.Println("Run OpenMCP API Server")
//	err := server.ListenAndServe()
//	if err != nil {
//		fmt.Println(err)
//	}
//
///*	l, err := net.Listen("tcp", ":"+GRPC_PORT)
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	grpcServer := grpc.NewServer()
//	st := &HttpManager{
//		ClusterManager: NewClusterManager(),
//	}
//
//	protobuf.RegisterRequestAPIServerServer(grpcServer, st)
//	if err := grpcServer.Serve(l); err != nil {
//		log.Fatalf("fail to serve: %v", err)
//	}*/
//}

// GET http://10.0.3.20:31635/token?username=openmcp&password=keti
// Get the Token
// Add Header
// --> Key : Authorization
// --> Value : Bearer {TOKEN}
// GET http://10.0.3.20:31635/api?clustername=openmcp
