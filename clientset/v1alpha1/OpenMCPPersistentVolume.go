package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPPersistentVolumeInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPPersistentVolume, error)
	Create(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error)
	Update(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error)
	UpdateStatus(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPPersistentVolumeClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPPersistentVolumeClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeList, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPPersistentVolume, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolume{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumes").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClient) Create(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolume{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcppersistentvolumes").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPPersistentVolumeClient) Update(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolume{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcppersistentvolumes").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPPersistentVolumeClient) UpdateStatus(cluster *resourcev1alpha1.OpenMCPPersistentVolume) (*resourcev1alpha1.OpenMCPPersistentVolume, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolume{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("OpenMCPPersistentVolumes").
		SubResource("status").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumes").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
