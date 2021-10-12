package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
)

func YamlApply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	var jsonErrs []jsonErr

	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), 1000)

	for {
		// /apis/openmcp.k8s.io/v1alpha1/namespaces/default/migrations?clustername=openmcp"
		var urlString string
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			break
		}

		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		if err != nil {
			log.Fatal("1  ", err)
		}
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		namespace := "default"

		if unstructuredObj.GetNamespace() != "" {
			namespace = unstructuredObj.GetNamespace()
		} else {
			namespace = "default"
		}
		// portal:
		// openmcpurl: 192.168.0.152
		// port:       31635
		// kubeconfig: config
		// openmcpURL2 := "192.168.0.142"

		// fmt.Println(obj)
		// fmt.Println(gvk.Kind, gvk.Version, gvk.Group, unstructuredObj.GetNamespace())
		if gvk.Group == "" {
			urlString = "https://" + openmcpURL + "/apis/" + gvk.Version + "/namespaces/" + namespace + "/" + gvk.Kind + "s?clustername=openmcp"
		} else {
			urlString = "https://" + openmcpURL + "/apis/" + gvk.Group + "/" + gvk.Version + "/namespaces/" + namespace + "/" + gvk.Kind + "s?clustername=openmcp"
		}
		urlString = strings.ToLower(urlString)
		pBody := bytes.NewBuffer(rawObj.Raw)
		// fmt.Println("urlString:    ", urlString)
		// fmt.Println("pBody:     ", pBody)
		resp, err := PostYaml(urlString, pBody)
		var msg jsonErr
		// fmt.Println("resp:     ", resp)

		if err != nil {
			msg = jsonErr{503, "failed", "request fail | " + gvk.Kind + " | " + namespace + " | " + unstructuredObj.GetName()}
			// json.NewEncoder(w).Encode(msg)
		}

		var data map[string]interface{}
		json.Unmarshal([]byte(resp), &data)
		if data != nil {
			if data["kind"].(string) == "Status" {
				msg = jsonErr{501, "failed", data["message"].(string) + " | " + gvk.Kind + " / " + namespace + " / " + unstructuredObj.GetName()}
				// json.NewEncoder(w).Encode(msg)
			} else {
				msg = jsonErr{200, "success", "Resource Created" + " | " + gvk.Kind + " / " + namespace + " / " + unstructuredObj.GetName()}
				// json.NewEncoder(w).Encode(msg)
			}
		}

		jsonErrs = append(jsonErrs, msg)
		// fmt.Println(jsonErrs)
	}
	json.NewEncoder(w).Encode(jsonErrs)
}
