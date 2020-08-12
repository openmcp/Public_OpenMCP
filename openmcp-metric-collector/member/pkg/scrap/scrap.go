package scrap

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/clock"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/decode"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/kubeletClient"
	"openmcp/openmcp/openmcp-metric-collector/member/pkg/storage"
)

func Scrap(config *rest.Config, kubelet_client *kubeletClient.KubeletClient, nodes []corev1.Node) (*storage.Collection, error) {
	klog.V(0).Info("Scrap Start")

	responseChannel := make(chan *storage.MetricsBatch, len(nodes))
	errChannel := make(chan error, len(nodes))
	defer close(responseChannel)
	defer close(errChannel)

	startTime := clock.MyClock.Now()

	//var wait serviceDNS.WaitGroup
	//wait.Add(len(cm.Node_list.Items))

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
	//wait.Wait()
	//time.Sleep(1 * time.Second)

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
		res.Matricsbatchs = append(res.Matricsbatchs, *srcBatch)
		//res.Matricsbatchs[i].Node = srcBatch.Node
		//res.Matricsbatchs[i].Pods = append(res.Matricsbatchs[i].Pods, srcBatch.Pods...)

		nodeNum += 1
		podNum += len(srcBatch.Pods)
	}
	klog.V(0).Infof("ScrapeMetrics: time: %s, nodes: %v, pods: %v", clock.MyClock.Since(startTime), nodeNum, podNum)
	return res, utilerrors.NewAggregate(errs)
}

func CollectNode(config *rest.Config, kubelet_client *kubeletClient.KubeletClient, node corev1.Node) (*storage.MetricsBatch, error) {
	klog.V(0).Info("Collect Node Start goroutine : '", node.Name, "'")
	host := node.Status.Addresses[0].Address
	token := config.BearerToken
	summary, err := kubelet_client.GetSummary(host, token)
	if err != nil {
		klog.V(0).Info("check1")
		return nil, fmt.Errorf("unable to fetch metrics from Kubelet %s (%s): %v", node.Name, node.Status.Addresses[0].Address, err)
	}

	return decode.DecodeBatch(summary)
}