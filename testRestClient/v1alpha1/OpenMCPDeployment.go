package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"openmcp/openmcp/testRestClient/rest"

	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

type OpenMCPDeploymentInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPDeploymentList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPDeployment, error)
	Create(deployment *v1alpha1.OpenMCPDeployment) (*v1alpha1.OpenMCPDeployment, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPDeploymentClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPDeploymentClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPDeploymentList, error) {
	result := v1alpha1.OpenMCPDeploymentList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPDeploymentClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPDeployment, error) {
	result := v1alpha1.OpenMCPDeployment{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPDeploymentClient) Create(deployment *v1alpha1.OpenMCPDeployment) (*v1alpha1.OpenMCPDeployment, error) {
	result := v1alpha1.OpenMCPDeployment{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpdeployments").
		Body(deployment).
		Do().
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
		Watch()
}
