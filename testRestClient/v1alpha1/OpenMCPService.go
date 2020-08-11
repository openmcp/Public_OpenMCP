package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"openmcp/openmcp/testRestClient/rest"

	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

type OpenMCPServiceInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPServiceList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPService, error)
	Create(deployment *v1alpha1.OpenMCPService) (*v1alpha1.OpenMCPService, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPServiceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPServiceClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPServiceList, error) {
	result := v1alpha1.OpenMCPServiceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPServiceClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPService, error) {
	result := v1alpha1.OpenMCPService{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpservices").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPServiceClient) Create(deployment *v1alpha1.OpenMCPService) (*v1alpha1.OpenMCPService, error) {
	result := v1alpha1.OpenMCPService{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpservices").
		Body(deployment).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPServiceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
