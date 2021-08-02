package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPSecretInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPSecretList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPSecret, error)
	Create(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error)
	Update(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error)
	UpdateStatus(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPSecretClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPSecretClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPSecretList, error) {
	result := resourcev1alpha1.OpenMCPSecretList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPSecretClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPSecret, error) {
	result := resourcev1alpha1.OpenMCPSecret{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpsecrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPSecretClient) Create(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error) {
	result := resourcev1alpha1.OpenMCPSecret{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpsecrets").
		Body(osec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPSecretClient) Update(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error) {
	result := resourcev1alpha1.OpenMCPSecret{}
	err := c.restClient.
		Put().
		Name(osec.Name).
		Namespace(c.ns).
		Resource("openmcpsecrets").
		Body(osec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPSecretClient) UpdateStatus(osec *resourcev1alpha1.OpenMCPSecret) (*resourcev1alpha1.OpenMCPSecret, error) {
	result := resourcev1alpha1.OpenMCPSecret{}
	err := c.restClient.
		Put().
		Name(osec.Name).
		Namespace(c.ns).
		Resource("openmcpsecrets").
		SubResource("status").
		Body(osec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPSecretClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpsecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
