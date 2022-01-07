package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPPersistentVolumeClaimInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeClaimList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error)
	Create(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error)
	Update(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error)
	UpdateStatus(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPPersistentVolumeClaimClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPPersistentVolumeClaimClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeClaimList, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeClaimList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClaimClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeClaim{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClaimClient) Create(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeClaim{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPPersistentVolumeClaimClient) Update(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeClaim{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPPersistentVolumeClaimClient) UpdateStatus(cluster *resourcev1alpha1.OpenMCPPersistentVolumeClaim) (*resourcev1alpha1.OpenMCPPersistentVolumeClaim, error) {
	result := resourcev1alpha1.OpenMCPPersistentVolumeClaim{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		SubResource("status").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPPersistentVolumeClaimClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcppersistentvolumeclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
