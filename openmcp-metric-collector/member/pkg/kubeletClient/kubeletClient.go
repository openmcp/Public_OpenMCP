package kubeletClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"net"
	"net/http"
	"net/url"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/stats"
	"strconv"
)

type ErrNotFound struct {
	endpoint string
}

func (err *ErrNotFound) Error() string {
	return fmt.Sprintf("%q not found", err.endpoint)
}

type KubeletClient struct {
	port            int
	deprecatedNoTLS bool
	client          *http.Client
}

func NewKubeletClient() (*KubeletClient, error) {
	klog.V(0).Info( "Func NewKubeletClient Called")
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	port := 10250
	deprecatedNoTLS := true

	c := &http.Client{
		Transport: transport,
	}
	return &KubeletClient{
		port:            port,
		client:          c,
		deprecatedNoTLS: deprecatedNoTLS,
	}, nil
}

func (kc *KubeletClient) GetSummary(host, token string) (*stats.Summary, error) {
	klog.V(0).Info( "GetSummary HTTP New Request, only_cpu_and_memory=false")

	scheme := "https"
	port := kc.port
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/stats/summary",
		//RawQuery: "only_cpu_and_memory=true",
	}
	klog.V(0).Info( "GetSummary URL: ", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		klog.V(0).Info( "check2")
		return nil, err
	}
	summary := &stats.Summary{}
	client := kc.client
	if client == nil {
		client = http.DefaultClient
	}
	err = kc.makeRequestAndGetValue(client, req.WithContext(context.TODO()), token, summary)

	summary.IP = host

	//req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer "+ token)
	//
	//resp, err := kc.client.Do(req)
	//if err != nil {
	//	klog.V(0).Info( "Check2", err)
	//	// handle err
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	klog.V(0).Info( "Check3", err)
	//	panic(err.Error())
	//}
	//
	//var prettyJSON bytes.Buffer
	//err = json.Indent(&prettyJSON, body, "", "\t")
	//if err != nil {
	//	klog.V(0).Info( "Check4", err)
	//	panic(err.Error())
	//}
	//fmt.Printf("%s\n", prettyJSON.Bytes())

	return summary, nil
}

func (kc *KubeletClient) makeRequestAndGetValue(client *http.Client, req *http.Request, token string, value interface{}) error {
	//klog.V(0).Info( "Get Metric Using Kubelet API")
	//klog.V(0).Info( "Func makeRequestAndGetValue Called")
	klog.V(0).Info( "makeRequestAndGetValue HTTP GET")
	// TODO(directxman12): support validating certs by hostname

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	//klog.V(0).Info( "makeRequestAndGetValue1")


	klog.V(0).Info( "Request Host:", req.Host)
	response, err := client.Do(req)
	klog.V(0).Info( "Status: ", response.Status)
	if err != nil {
		klog.V(0).Info( "check3")
		return err
	}
	//klog.V(0).Info( "makeRequestAndGetValue2")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		klog.V(0).Info( "check4")
		return fmt.Errorf("failed to read response body - %v", err)
	}
	//klog.V(0).Info( "makeRequestAndGetValue3")
	if response.StatusCode == http.StatusNotFound {
		klog.V(0).Info( "check55")
		return &ErrNotFound{req.URL.String()}
	} else if response.StatusCode != http.StatusOK {
		klog.V(0).Info( "check5")
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	}
	//klog.V(0).Info( "makeRequestAndGetValue4")
	kubeletAddr := "[unknown]"
	if req.URL != nil {
		kubeletAddr = req.URL.Host
	}
	//klog.V(0).Info( "makeRequestAndGetValue5")
	klog.V(10).Infof("Raw response from Kubelet at %s: %s", kubeletAddr, string(body))

	//////////////////////////////////////////
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		klog.V(0).Info( "Check11111111111111", err)
		panic(err.Error())
	}
	klog.V(0).Infof("%s%s\n", prettyJSON.Bytes()[:400],"...................")
	//klog.V(0).Infof(1,"%s\n", prettyJSON.Bytes())
	//////////////////////////////////////////
	//klog.V(0).Info( "makeRequestAndGetValue6")

	err = json.Unmarshal(body, value)
	if err != nil {
		klog.V(0).Info("check6")
		return fmt.Errorf("failed to parse output. Response: %q. Error: %v", string(body), err)
	}
	//klog.V(0).Info( "makeRequestAndGetValue7")
	return nil
}
