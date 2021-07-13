package scrap

import (
	"fmt"
	"openmcp/openmcp/openmcp-metric-collector/member/src/clock"
	"openmcp/openmcp/openmcp-metric-collector/member/src/decode"
	"openmcp/openmcp/openmcp-metric-collector/member/src/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/src/storage"
	"os"

	corev1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
)

func Scrap(config *rest.Config, kubelet_client *kubeletClient.KubeletClient, nodes []corev1.Node) (*storage.Collection, error) {
	fmt.Println("Func Scrap Called")

	responseChannel := make(chan *storage.MetricsBatch, len(nodes))
	errChannel := make(chan error, len(nodes))
	defer close(responseChannel)
	defer close(errChannel)

	startTime := clock.MyClock.Now()

	for _, node := range nodes {
		go func(node corev1.Node) {
			//defer wait.Done()
			metrics, err := CollectNode(config, kubelet_client, node)
			if err != nil {
				err = fmt.Errorf("unable to fully scrape metrics from node %s: %v", node.Name, err)
			}
			responseChannel <- metrics
			errChannel <- err
		}(node)

	}

	var errs []error
	res := &storage.Collection{}
	nodeNum := 0
	podNum := 0
	for range nodes {
		err := <-errChannel
		srcBatch := <-responseChannel
		if err != nil {
			errs = append(errs, err)
			// NB: partial node results are still worth saving, so
			// don't skip storing results if we got an error
		}
		if srcBatch == nil {
			continue
		}
		res.Metricsbatchs = append(res.Metricsbatchs, *srcBatch)

		nodeNum += 1
		podNum += len(srcBatch.Pods)
	}

	res.ClusterName = os.Getenv("CLUSTER_NAME") //config.Username

	fmt.Println("ScrapeMetrics: time: ", clock.MyClock.Since(startTime), "nodes: ", nodeNum, "pods: ", podNum)
	return res, utilerrors.NewAggregate(errs)
}

func CollectNode(config *rest.Config, kubelet_client *kubeletClient.KubeletClient, node corev1.Node) (*storage.MetricsBatch, error) {
	fmt.Println("Func CollectNode Called")
	fmt.Println("Collect Node Start goroutine : '", node.Name, "'")
	host := node.Status.Addresses[0].Address
	token := config.BearerToken
	summary, err := kubelet_client.GetSummary(host, token)
	fmt.Println("summary : ", summary)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch metrics from Kubelet %s (%s): %v", node.Name, node.Status.Addresses[0].Address, err)
	}

	return decode.DecodeBatch(summary)
}
