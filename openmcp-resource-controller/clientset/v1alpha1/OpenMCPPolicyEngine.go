package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"openmcp/openmcp/openmcp-resource-controller/apis/keti/v1alpha1"
)

type OpenMCPPolicyInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPPolicyList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPPolicy, error)
	Create(deployment *v1alpha1.OpenMCPPolicy) (*v1alpha1.OpenMCPPolicy, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPPolicyClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPPolicyClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPPolicyList, error) {
	result := v1alpha1.OpenMCPPolicyList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPPolicy, error) {
	result := v1alpha1.OpenMCPPolicy{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicys").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyClient) Create(deployment *v1alpha1.OpenMCPPolicy) (*v1alpha1.OpenMCPPolicy, error) {
	result := v1alpha1.OpenMCPPolicy{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcppolicys").
		Body(deployment).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
