package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPDaemonSetInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPDaemonSetList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPDaemonSet, error)
	Create(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error)
	Update(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error)
	UpdateStatus(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPDaemonSetClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPDaemonSetClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPDaemonSetList, error) {
	result := resourcev1alpha1.OpenMCPDaemonSetList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDaemonSetClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPDaemonSet, error) {
	result := resourcev1alpha1.OpenMCPDaemonSet{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDaemonSetClient) Create(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error) {
	result := resourcev1alpha1.OpenMCPDaemonSet{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPDaemonSetClient) Update(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error) {
	result := resourcev1alpha1.OpenMCPDaemonSet{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPDaemonSetClient) UpdateStatus(cluster *resourcev1alpha1.OpenMCPDaemonSet) (*resourcev1alpha1.OpenMCPDaemonSet, error) {
	result := resourcev1alpha1.OpenMCPDaemonSet{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		SubResource("status").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPDaemonSetClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpdaemonsets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
