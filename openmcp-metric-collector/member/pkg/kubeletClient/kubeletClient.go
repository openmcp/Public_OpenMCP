package kubeletClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"openmcp/openmcp/omcplog"
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
	omcplog.V(4).Info( "Func NewKubeletClient Called")
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
	omcplog.V(4).Info( "Func GetSummary Called")
	omcplog.V(2).Info( "GetSummary HTTP New Request, only_cpu_and_memory=false")

	scheme := "https"
	port := kc.port
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/stats/summary",
		//RawQuery: "only_cpu_and_memory=true",
	}
	omcplog.V(3).Info( "GetSummary URL: ", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		omcplog.V(0).Info(err)
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
	omcplog.V(4).Info( "Func makeRequestAndGetValue Called")

	// TODO(directxman12): support validating certs by hostname

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	//klog.V(0).Info( "makeRequestAndGetValue1")


	omcplog.V(3).Info( "Request Host:", req.Host)
	response, err := client.Do(req)
	omcplog.V(3).Info( "Status: ", response.Status)
	if err != nil {
		omcplog.V(0).Info( err)
		return err
	}
	//klog.V(0).Info( "makeRequestAndGetValue2")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body - %v", err)
	}
	//klog.V(0).Info( "makeRequestAndGetValue3")
	if response.StatusCode == http.StatusNotFound {
		return &ErrNotFound{req.URL.String()}
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	}
	//klog.V(0).Info( "makeRequestAndGetValue4")
	kubeletAddr := "[unknown]"
	if req.URL != nil {
		kubeletAddr = req.URL.Host
	}
	//klog.V(0).Info( "makeRequestAndGetValue5")
	omcplog.V(5).Infof("Raw response from Kubelet at %s: %s", kubeletAddr, string(body))

	//////////////////////////////////////////
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		omcplog.V(0).Info( err)
		panic(err.Error())
	}
	omcplog.V(3).Infof("%s%s\n", prettyJSON.Bytes()[:400],"...................")
	omcplog.V(3).Infof("%s\n", prettyJSON.Bytes())
	//////////////////////////////////////////


	err = json.Unmarshal(body, value)
	if err != nil {
		return fmt.Errorf("failed to parse output. Response: %q. Error: %v", string(body), err)
	}

	return nil
}
