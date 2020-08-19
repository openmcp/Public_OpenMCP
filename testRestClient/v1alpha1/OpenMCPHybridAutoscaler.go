package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"openmcp/openmcp/testRestClient/rest"

	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

type OpenMCPHybridAutoScalerInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPHybridAutoScalerList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPHybridAutoScaler, error)
	Create(deployment *v1alpha1.OpenMCPHybridAutoScaler) (*v1alpha1.OpenMCPHybridAutoScaler, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPHybridAutoScalerClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPHybridAutoScalerClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPHybridAutoScalerList, error) {
	result := v1alpha1.OpenMCPHybridAutoScalerList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPHybridAutoScalerClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPHybridAutoScaler, error) {
	result := v1alpha1.OpenMCPHybridAutoScaler{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPHybridAutoScalerClient) Create(deployment *v1alpha1.OpenMCPHybridAutoScaler) (*v1alpha1.OpenMCPHybridAutoScaler, error) {
	result := v1alpha1.OpenMCPHybridAutoScaler{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		Body(deployment).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPHybridAutoScalerClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
