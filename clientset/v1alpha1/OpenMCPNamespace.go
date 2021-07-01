package v1alpha1

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"
)

type OpenMCPNamespaceInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPNamespaceList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPNamespace, error)
	Create(deployment *resourcev1alpha1.OpenMCPNamespace) (*resourcev1alpha1.OpenMCPNamespace, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPNamespaceClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPNamespaceClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPNamespaceList, error) {
	result := resourcev1alpha1.OpenMCPNamespaceList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpnamespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPNamespaceClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPNamespace, error) {
	result := resourcev1alpha1.OpenMCPNamespace{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpnamespaces").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPNamespaceClient) Create(deployment *resourcev1alpha1.OpenMCPNamespace) (*resourcev1alpha1.OpenMCPNamespace, error) {
	result := resourcev1alpha1.OpenMCPNamespace{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpnamespaces").
		Body(deployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPNamespaceClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpnamespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
