package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPConfigMapInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPConfigMapList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPConfigMap, error)
	Create(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error)
	Update(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error)
	UpdateStatus(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPConfigMapClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPConfigMapClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPConfigMapList, error) {
	result := resourcev1alpha1.OpenMCPConfigMapList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPConfigMapClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPConfigMap, error) {
	result := resourcev1alpha1.OpenMCPConfigMap{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPConfigMapClient) Create(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error) {
	result := resourcev1alpha1.OpenMCPConfigMap{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		Body(ocm).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPConfigMapClient) Update(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error) {
	result := resourcev1alpha1.OpenMCPConfigMap{}
	err := c.restClient.
		Put().
		Name(ocm.Name).
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		Body(ocm).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPConfigMapClient) UpdateStatus(ocm *resourcev1alpha1.OpenMCPConfigMap) (*resourcev1alpha1.OpenMCPConfigMap, error) {
	result := resourcev1alpha1.OpenMCPConfigMap{}
	err := c.restClient.
		Put().
		Name(ocm.Name).
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		SubResource("status").
		Body(ocm).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPConfigMapClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpconfigmaps").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
