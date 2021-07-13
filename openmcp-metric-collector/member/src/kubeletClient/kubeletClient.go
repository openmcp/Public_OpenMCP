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
	"openmcp/openmcp/openmcp-metric-collector/member/src/stats"
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
	fmt.Println("Func GetSummary Called")
	fmt.Println("GetSummary HTTP New Request, only_cpu_and_memory=false")

	scheme := "https"
	port := kc.port
	url := url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		Path:   "/stats/summary",
		//RawQuery: "only_cpu_and_memory=true",
	}
	fmt.Println("GetSummary URL: ", url.String())
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	summary := &stats.Summary{}
	client := kc.client
	if client == nil {
		client = http.DefaultClient
	}
	err = kc.makeRequestAndGetValue(client, req.WithContext(context.TODO()), token, summary)

	summary.IP = host

	return summary, nil
}

func (kc *KubeletClient) makeRequestAndGetValue(client *http.Client, req *http.Request, token string, value interface{}) error {
	fmt.Println("Func makeRequestAndGetValue Called")

	// TODO(directxman12): support validating certs by hostname

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("Request Host:", req.Host)
	response, err := client.Do(req)
	fmt.Println("Status: ", response.Status)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body - %v", err)
	}
	if response.StatusCode == http.StatusNotFound {
		return &ErrNotFound{req.URL.String()}
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed - %q, response: %q", response.Status, string(body))
	}
	kubeletAddr := "[unknown]"
	if req.URL != nil {
		kubeletAddr = req.URL.Host
	}
	fmt.Println("Raw response from Kubelet at ", kubeletAddr)
	//fmt.Println("Raw response from Kubelet at %s: %s", kubeletAddr, string(body))

	//////////////////////////////////////////
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	//fmt.Println("%s%s\n", prettyJSON.Bytes()[:400],"...................")
	//fmt.Printf("%s\n", prettyJSON.Bytes())
	//////////////////////////////////////////

	err = json.Unmarshal(body, value)
	if err != nil {
		return fmt.Errorf("failed to parse output. Response: %q. Error: %v", string(body), err)
	}

	return nil
}
