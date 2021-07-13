package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPServiceInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPServiceList, error)
	Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPService, error)
	Create(deployment *resourcev1alpha1.OpenMCPService) (*resourcev1alpha1.OpenMCPService, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPServiceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPServiceClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPServiceList, error) {
	result := resourcev1alpha1.OpenMCPServiceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPServiceClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPService, error) {
	result := resourcev1alpha1.OpenMCPService{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpservices").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPServiceClient) Create(deployment *resourcev1alpha1.OpenMCPService) (*resourcev1alpha1.OpenMCPService, error) {
	result := resourcev1alpha1.OpenMCPService{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpservices").
		Body(deployment).
		Do(context.TODO()).
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
		Watch(context.TODO())
}
