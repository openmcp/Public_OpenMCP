package protobuf

const (
	GRPC_PORT = "8080"
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
