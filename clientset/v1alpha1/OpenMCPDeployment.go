package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPDeploymentInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPDeploymentList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPDeployment, error)
	Create(deployment *resourcev1alpha1.OpenMCPDeployment) (*resourcev1alpha1.OpenMCPDeployment, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPDeploymentClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPDeploymentClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPDeploymentList, error) {
	result := resourcev1alpha1.OpenMCPDeploymentList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDeploymentClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPDeployment, error) {
	result := resourcev1alpha1.OpenMCPDeployment{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDeploymentClient) Create(deployment *resourcev1alpha1.OpenMCPDeployment) (*resourcev1alpha1.OpenMCPDeployment, error) {
	result := resourcev1alpha1.OpenMCPDeployment{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		Body(deployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDeploymentClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdeployment").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
