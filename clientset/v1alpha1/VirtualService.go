package v1alpha1

import (
	"context"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type VirtualServiceInterface interface {
	List(opts metav1.ListOptions) (*v1alpha3.VirtualServiceList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha3.VirtualService, error)
	Create(vs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error)
	Update(vs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type VirtualServiceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *VirtualServiceClient) List(opts metav1.ListOptions) (*v1alpha3.VirtualServiceList, error) {
	result := v1alpha3.VirtualServiceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("virtualservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *VirtualServiceClient) Get(name string, opts metav1.GetOptions) (*v1alpha3.VirtualService, error) {
	result := v1alpha3.VirtualService{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("virtualservices").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *VirtualServiceClient) Create(vs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
	result := v1alpha3.VirtualService{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("virtualservices").
		Body(vs).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *VirtualServiceClient) Update(vs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
	result := v1alpha3.VirtualService{}
	err := c.restClient.
		Put().
		Name(vs.Name).
		Namespace(c.ns).
		Resource("virtualservices").
		Body(vs).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *VirtualServiceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("virtualservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
