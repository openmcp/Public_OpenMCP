package cmd

import (
	"context"
	"fmt"
	"log"
	"openmcp/openmcp/omcpctl/apiServerMethod"
	"openmcp/openmcp/omcpctl/resource"
	cobrautil "openmcp/openmcp/omcpctl/util"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func Test_CacheResource(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.0.0.226:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	localTime := time.Now().Local()
	date := localTime.Format("20060102")
	prefixString := "cache/" + string(date) + "/"
	fmt.Print(prefixString)
	// fmt.Print(date)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, prefixString, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		fmt.Print(err)
	}
	cancel()
	var m map[string]int
	m = make(map[string]int)
	for _, ev := range resp.Kvs {
		keys := strings.Split(strings.TrimSpace(string(ev.Key)), "/")
		key := ""
		for i := 3; i < len(keys); i++ {
			if i == 3 {
				key += keys[i]
			} else {
				key += "/" + keys[i]
			}
		}
		value, _ := strconv.Atoi(strings.TrimSpace(string(ev.Value)))
		m[key] += value
	}
	clusterName := "cluster2"
	LINK := cobrautil.GetLinkParser("nodes", "", "default", clusterName)
	fmt.Println(LINK)

	body, err := apiServerMethod.GetAPIServer(LINK)
	if err != nil {
		fmt.Println("error: the server doesn't have a resource type 'node'")
	}
	result := resource.GetNodeList(body)
	// result := util.CmdExec2("omcpctl get node --context " + clusterName)
	fmt.Println(result)
	imagelist := ""
	for key, val := range m {
		if val >= 6 {
			keys := strings.Split(key, ":")
			str := "  - imageName: \"" + keys[0] + "\"\n" + "    tagName: \"" + keys[1] + "\"\n"
			imagelist += str
		}
	}
	now := time.Now().Unix()
	nodeIP := "10.0.0.226"
	yaml := "apiVersion: openmcp.k8s.io/v1alpha1\n" +
		"kind: NodeRegistry\n" +
		"metadata:\n" +
		"  name: " + strconv.FormatInt(now, 10) + "-image-pull\n" +
		"spec:\n" +
		"  command: \"pull\"\n" +
		"  clusterName: \"" + clusterName + "\"\n" +
		"  nodeName: \"" + nodeIP + "\"\n" +
		"  imageList:\n" +
		imagelist

	fmt.Print(yaml)
}
