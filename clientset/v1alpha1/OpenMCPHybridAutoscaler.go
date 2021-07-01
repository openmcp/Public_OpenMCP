package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPHybridAutoScalerInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPHybridAutoScalerList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPHybridAutoScaler, error)
	Create(deployment *resourcev1alpha1.OpenMCPHybridAutoScaler) (*resourcev1alpha1.OpenMCPHybridAutoScaler, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPHybridAutoScalerClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPHybridAutoScalerClient) reate(deployment resourcev1alpha1.OpenMCPHybridAutoScaler) (*resourcev1alpha1.OpenMCPHybridAutoScaler, error) {
	panic("implement me")
}

func (c *OpenMCPHybridAutoScalerClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPHybridAutoScalerList, error) {
	result := resourcev1alpha1.OpenMCPHybridAutoScalerList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPHybridAutoScalerClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPHybridAutoScaler, error) {
	result := resourcev1alpha1.OpenMCPHybridAutoScaler{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPHybridAutoScalerClient) Create(deployment *resourcev1alpha1.OpenMCPHybridAutoScaler) (*resourcev1alpha1.OpenMCPHybridAutoScaler, error) {
	result := resourcev1alpha1.OpenMCPHybridAutoScaler{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcphybridautoscalers").
		Body(deployment).
		Do(context.TODO()).
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
		Watch(context.TODO())
}
