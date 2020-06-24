package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"resource-controller/apis/keti/v1alpha1"
)

type OpenMCPPolicyEngineInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.OpenMCPPolicyEngineList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.OpenMCPPolicyEngine, error)
	Create(deployment *v1alpha1.OpenMCPPolicyEngine) (*v1alpha1.OpenMCPPolicyEngine, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPPolicyEngineClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPPolicyEngineClient) List(opts metav1.ListOptions) (*v1alpha1.OpenMCPPolicyEngineList, error) {
	result := v1alpha1.OpenMCPPolicyEngineList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicyengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyEngineClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.OpenMCPPolicyEngine, error) {
	result := v1alpha1.OpenMCPPolicyEngine{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicyengines").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyEngineClient) Create(deployment *v1alpha1.OpenMCPPolicyEngine) (*v1alpha1.OpenMCPPolicyEngine, error) {
	result := v1alpha1.OpenMCPPolicyEngine{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcppolicyengines").
		Body(deployment).
		Do().
		Into(&result)

	return &result, err
}

func (c *OpenMCPPolicyEngineClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppolicyengines").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
