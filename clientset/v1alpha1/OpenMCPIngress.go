package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPIngressInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPIngressList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPIngress, error)
	Create(oing *resourcev1alpha1.OpenMCPIngress) (*resourcev1alpha1.OpenMCPIngress, error)
	UpdateStatus(oing *resourcev1alpha1.OpenMCPIngress) (*resourcev1alpha1.OpenMCPIngress, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPIngressClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPIngressClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPIngressList, error) {
	result := resourcev1alpha1.OpenMCPIngressList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPIngress, error) {
	result := resourcev1alpha1.OpenMCPIngress{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresses").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Create(oing *resourcev1alpha1.OpenMCPIngress) (*resourcev1alpha1.OpenMCPIngress, error) {
	result := resourcev1alpha1.OpenMCPIngress{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpingresses").
		Body(oing).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPIngressClient) UpdateStatus(oing *resourcev1alpha1.OpenMCPIngress) (*resourcev1alpha1.OpenMCPIngress, error) {
	result := resourcev1alpha1.OpenMCPIngress{}
	err := c.restClient.
		Put().
		Name(oing.Name).
		Namespace(c.ns).
		Resource("openmcpingresses").
		SubResource("status").
		Body(oing).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
