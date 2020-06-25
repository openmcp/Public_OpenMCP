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
	fmt.Println("Func NewKubeletClient Called")
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
	//fmt.Println("Func GetSummary Called")
	scheme := "https"
	port := kc.port
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/stats/summary",
		//RawQuery: "only_cpu_and_memory=true",
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println("check2")
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
	//	fmt.Println("Check2", err)
	//	// handle err
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Println("Check3", err)
	//	panic(err.Error())
	//}
	//
	//var prettyJSON bytes.Buffer
	//err = json.Indent(&prettyJSON, body, "", "\t")
	//if err != nil {
	//	fmt.Println("Check4", err)
	//	panic(err.Error())
	//}
	//fmt.Printf("%s\n", prettyJSON.Bytes())

	return summary, nil
}

func (kc *KubeletClient) makeRequestAndGetValue(client *http.Client, req *http.Request, token string, value interface{}) error {
	fmt.Println("Get Metric Using Kubelet API")
	//fmt.Println("Func makeRequestAndGetValue Called")
	// TODO(directxman12): support validating certs by hostname

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	//fmt.Println("makeRequestAndGetValue1")

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("check3")
		return err
	}
	//fmt.Println("makeRequestAndGetValue2")
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("check4")
		return fmt.Errorf("failed to read response body - %v", err)
	}
	//fmt.Println("makeRequestAndGetValue3")
	if response.StatusCode == http.StatusNotFound {
		fmt.Println("check55")
		return &ErrNotFound{req.URL.String()}
	} else if response.StatusCode != http.StatusOK {
		fmt.Println("check5")
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	}
	//fmt.Println("makeRequestAndGetValue4")
	kubeletAddr := "[unknown]"
	if req.URL != nil {
		kubeletAddr = req.URL.Host
	}
	//fmt.Println("makeRequestAndGetValue5")
	klog.V(10).Infof("Raw response from Kubelet at %s: %s", kubeletAddr, string(body))

	//////////////////////////////////////////
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		fmt.Println("Check11111111111111", err)
		panic(err.Error())
	}
	fmt.Printf("%s\n", prettyJSON.Bytes())
	//////////////////////////////////////////
	//fmt.Println("makeRequestAndGetValue6")

	err = json.Unmarshal(body, value)
	if err != nil {
		fmt.Println("check6")
		return fmt.Errorf("failed to parse output. Response: %q. Error: %v", string(body), err)
	}
	//fmt.Println("makeRequestAndGetValue7")
	return nil
}
