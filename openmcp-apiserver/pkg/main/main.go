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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"

	//	corev1 "k8s.io/api/core/v1"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/remotecommand"
	"strings"
	"time"

	//"github.com/auth0/go-jwt-middleware"
	//"github.com/dgrijalva/jwt-go"
	//"google.golang.org/grpc"
	//"io"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"log"
	//"net"
	"net/http"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
	"sigs.k8s.io/kubefed/pkg/controller/util"
	//"time"

	//"openmcp/openmcp/openmcp-apiserver/pkg/protobuf"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes/scheme"
)

const (
	APP_KEY = "openmcp-apiserver"
)

type ClusterManager struct {
	Fed_namespace   string
	Host_config     *rest.Config
	Host_client     genericclient.Client
	Cluster_list    *fedv1b1.KubeFedClusterList
	Cluster_configs map[string]*rest.Config
	Cluster_clients map[string]genericclient.Client
}

func ListKubeFedClusters(client genericclient.Client, namespace string) *fedv1b1.KubeFedClusterList {
	clusterList := &fedv1b1.KubeFedClusterList{}
	err := client.List(context.TODO(), clusterList, namespace)
	if err != nil {
		fmt.Println("Error retrieving list of federated clusters: %+v", err)
	}
	if len(clusterList.Items) == 0 {
		fmt.Println("No federated clusters found")
	}
	return clusterList
}

func KubeFedClusterConfigs(clusterList *fedv1b1.KubeFedClusterList, client genericclient.Client, fedNamespace string) map[string]*rest.Config {
	clusterConfigs := make(map[string]*rest.Config)
	for _, cluster := range clusterList.Items {
		config, _ := util.BuildClusterConfig(&cluster, client, fedNamespace)
		clusterConfigs[cluster.Name] = config
	}
	return clusterConfigs
}
func KubeFedClusterClients(clusterList *fedv1b1.KubeFedClusterList, cluster_configs map[string]*rest.Config) map[string]genericclient.Client {

	cluster_clients := make(map[string]genericclient.Client)
	for _, cluster := range clusterList.Items {
		clusterName := cluster.Name
		cluster_config := cluster_configs[clusterName]
		cluster_client := genericclient.NewForConfigOrDie(cluster_config)
		cluster_clients[clusterName] = cluster_client
	}
	return cluster_clients
}

func NewClusterManager() *ClusterManager {
	fed_namespace := "kube-federation-system"
	host_config, _ := rest.InClusterConfig()
	host_client := genericclient.NewForConfigOrDie(host_config)
	cluster_list := ListKubeFedClusters(host_client, fed_namespace)
	cluster_configs := KubeFedClusterConfigs(cluster_list, host_client, fed_namespace)
	cluster_clients := KubeFedClusterClients(cluster_list, cluster_configs)

	cm := &ClusterManager{
		Fed_namespace:   fed_namespace,
		Host_config:     host_config,
		Host_client:     host_client,
		Cluster_list:    cluster_list,
		Cluster_configs: cluster_configs,
		Cluster_clients: cluster_clients,
	}
	return cm
}

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
	fmt.Println(restclient)
	fmt.Println(config)
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

		fmt.Println("host ", r.Host)
		fmt.Println("url ", r.URL)
		fmt.Println("url/host ", r.URL.Host)
		fmt.Println("body ", r.Body)
		fmt.Println(r.URL.Query()["containername"])
		fmt.Println(r.URL.Query()["clustername"])
		fmt.Println(r.URL.Query()["podname"])
		fmt.Println(r.URL.Query()["podnamespace"])
		fmt.Println(r.URL.Query()["stdin"])
		fmt.Println(r.URL.Query()["tty"])
		fmt.Println(r.URL.Query()["stdout"])
		fmt.Println(r.URL.Query()["stderr"])

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

		fmt.Println(bodyString)

		stdin_io := bytes.NewBufferString(stdin)
		stdout_io := bytes.NewBufferString(stdout)
		stderr_io := bytes.NewBufferString(stderr)
		//tty_io := bytes.NewBufferString(tty)

		var a []string

		a = strings.Split(bodyString, ",")

		fmt.Println(a)
		fmt.Println(a[0])
		//fmt.Println(a[1])

		restClient, _ := restclient.RESTClientFor(h.ClusterManager.Host_config)

		ExecCmdExample(restClient,h.ClusterManager.Host_config ,podname, podnamespace, a, "", stdin_io, stdout_io, stderr_io)

	}else {

		clusterNames, ok := r.URL.Query()["clustername"]

		fmt.Println(clusterNames, ok)
		if !ok || len(clusterNames[0]) < 1 {
			w.Write([]byte("Url Param 'clustername' is missing"))
			return
		}
		fmt.Println()

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
			fmt.Println("Check1", err)
			// handle err
		}

		fmt.Println(r.Header.Get("Content-Type"))
		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
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

		//fmt.Printf("%s\n", prettyJSON.Bytes())
		w.Write(body)
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
	ClusterManager  *ClusterManager
}
const (
	GRPC_PORT = "8080"

)


func main() {

	//HTTPServer_IP := "10.0.3.20"
	HTTPServer_PORT := "8080"

	cm := NewClusterManager()

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

	fmt.Println("Run OpenMCP API Server")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}

/*	l, err := net.Listen("tcp", ":"+GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	st := &HttpManager{
		ClusterManager: NewClusterManager(),
	}

	protobuf.RegisterRequestAPIServerServer(grpcServer, st)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}*/
}

// GET http://10.0.3.20:31635/token?username=openmcp&password=keti
// Get the Token
// Add Header
// --> Key : Authorization
// --> Value : Bearer {TOKEN}
// GET http://10.0.3.20:31635/api?clustername=openmcp
