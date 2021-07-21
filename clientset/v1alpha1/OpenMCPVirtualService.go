package v1alpha1

import (
	"context"

	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPVirtualServiceInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPVirtualServiceList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPVirtualService, error)
	Create(vs *resourcev1alpha1.OpenMCPVirtualService) (*resourcev1alpha1.OpenMCPVirtualService, error)
	Update(vs *resourcev1alpha1.OpenMCPVirtualService) (*resourcev1alpha1.OpenMCPVirtualService, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPVirtualServiceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPVirtualServiceClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPVirtualServiceList, error) {
	result := resourcev1alpha1.OpenMCPVirtualServiceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpvirtualservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPVirtualServiceClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPVirtualService, error) {
	result := resourcev1alpha1.OpenMCPVirtualService{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpvirtualservices").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPVirtualServiceClient) Create(vs *resourcev1alpha1.OpenMCPVirtualService) (*resourcev1alpha1.OpenMCPVirtualService, error) {
	result := resourcev1alpha1.OpenMCPVirtualService{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpvirtualservices").
		Body(vs).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPVirtualServiceClient) Update(vs *resourcev1alpha1.OpenMCPVirtualService) (*resourcev1alpha1.OpenMCPVirtualService, error) {
	result := resourcev1alpha1.OpenMCPVirtualService{}
	err := c.restClient.
		Put().
		Name(vs.Name).
		Namespace(c.ns).
		Resource("openmcpvirtualservices").
		Body(vs).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPVirtualServiceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("virtualservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
