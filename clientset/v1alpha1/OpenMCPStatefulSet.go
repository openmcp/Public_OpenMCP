package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPStatefulSetInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPStatefulSetList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPStatefulSet, error)
	Create(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error)
	Update(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error)
	UpdateStatus(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPStatefulSetClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPStatefulSetClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPStatefulSetList, error) {
	result := resourcev1alpha1.OpenMCPStatefulSetList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPStatefulSetClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPStatefulSet, error) {
	result := resourcev1alpha1.OpenMCPStatefulSet{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPStatefulSetClient) Create(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error) {
	result := resourcev1alpha1.OpenMCPStatefulSet{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPStatefulSetClient) Update(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error) {
	result := resourcev1alpha1.OpenMCPStatefulSet{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPStatefulSetClient) UpdateStatus(cluster *resourcev1alpha1.OpenMCPStatefulSet) (*resourcev1alpha1.OpenMCPStatefulSet, error) {
	result := resourcev1alpha1.OpenMCPStatefulSet{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		SubResource("status").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPStatefulSetClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpstatefulsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
