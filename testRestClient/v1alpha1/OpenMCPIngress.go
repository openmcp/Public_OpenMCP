package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"openmcp/openmcp/testRestClient/rest"

	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

type OpenMCPIngressInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPIngressList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPIngress, error)
	Create(deployment *v1alpha1.OpenMCPIngress) (*v1alpha1.OpenMCPIngress, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPIngressClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPIngressClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPIngressList, error) {
	result := v1alpha1.OpenMCPIngressList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresss").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPIngress, error) {
	result := v1alpha1.OpenMCPIngress{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresss").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Create(deployment *v1alpha1.OpenMCPIngress) (*v1alpha1.OpenMCPIngress, error) {
	result := v1alpha1.OpenMCPIngress{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpingresss").
		Body(deployment).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPIngressClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpingresss").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
