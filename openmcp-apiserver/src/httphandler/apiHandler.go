package httphandler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func (h *HttpManager) ApiHandler(w http.ResponseWriter, r *http.Request) {
	// Migration
	// POST http://10.0.3.20:31635/apis/openmcp.k8s.io/v1alpha1/namespaces/openmcp/migrations?clustername=openmcp (yaml)

	if strings.Contains(r.URL.Path, "exec") {
		clusterNames, ok := r.URL.Query()["clustername"]

		if !ok || len(clusterNames[0]) < 1 {
			w.Write([]byte("Url Param 'clustername' is missing"))
			return
		}

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

		podname := r.URL.Query()["podname"][0]
		podnamespace := r.URL.Query()["podnamespace"][0]
		stdin = r.URL.Query()["stdin"][0]
		stdout = r.URL.Query()["stdout"][0]
		stderr = r.URL.Query()["stderr"][0]

		body, _ := ioutil.ReadAll(r.Body)
		bodyString := string(body)

		omcplog.V(5).Info(bodyString)

		stdin_io := bytes.NewBufferString(stdin)
		stdout_io := bytes.NewBufferString(stdout)
		stderr_io := bytes.NewBufferString(stderr)

		var a []string

		a = strings.Split(bodyString, ",")

		omcplog.V(5).Info(a)
		omcplog.V(5).Info(a[0])

		restClient, _ := restclient.RESTClientFor(h.ClusterManager.Host_config)

		ExecCmdExample(restClient, h.ClusterManager.Host_config, podname, podnamespace, a, "", stdin_io, stdout_io, stderr_io)

	} else {

		clusterNames, ok := r.URL.Query()["clustername"]

		omcplog.V(5).Info(clusterNames, ok)
		if !ok || len(clusterNames[0]) < 1 {
			w.Write([]byte("Url Param 'clustername' is missing"))
			return
		}

		APISERVER := ""
		TOKEN := ""
		clusterName := clusterNames[0]
		//RESTART:
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

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}

		var req *http.Request
		var err error

		if r.Method == "GET" || r.Method == "DELETE" {
			req, err = http.NewRequest(r.Method, APISERVER+r.URL.Path, nil)
		} else if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			req, err = http.NewRequest(r.Method, APISERVER+r.URL.Path, r.Body)
		}

		if err != nil {
			// handle err
			omcplog.V(0).Info(err)
		}

		omcplog.V(5).Info("Content-Type : ", r.Header.Get("Content-Type"))
		omcplog.V(5).Info("Authorization", "Bearer "+TOKEN)

		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		//req.Header.Set("Content-Type", "application/yaml")
		req.Header.Set("Authorization", "Bearer "+TOKEN)

		resp, err := client.Do(req)
		omcplog.V(3).Info("Request Done!")
		if err != nil {
			// handle err
			omcplog.V(0).Info(err)
			h.ClusterManager = clusterManager.NewClusterManager()
			return
			//goto RESTART
		}
		defer resp.Body.Close()

		if resp.Status >= "400" {
			omcplog.V(0).Info("clusterName:", clusterName, " / resp.Status:", resp.Status)
			return

		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			omcplog.V(0).Info(err)
			//panic(err.Error())
			return
		}

		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, body, "", "\t")
		if err != nil {
			omcplog.V(0).Info(err)
			//panic(err.Error())
			return

		}

		//omcplog.V(5).Info(string(prettyJSON.Bytes()))

		w.Write(body)

	}

}

func ExecCmdExample(restclient *restclient.RESTClient, config *restclient.Config, podName string, podNamespace string,
	command []string, containerName string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	omcplog.V(5).Info(restclient)
	omcplog.V(5).Info(config)
	req := restclient.Post().Resource("pods").Name(podName).
		Namespace(podNamespace).SubResource("exec")
	option := &corev1.PodExecOptions{
		Container: containerName,
		Command:   command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
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
