package v1alpha1

import (
	"context"
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPJobInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPJobList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.OpenMCPJob, error)
	Create(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error)
	Update(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error)
	UpdateStatus(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPJobClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPJobClient) List(opts metav1.ListOptions) (*resourcev1alpha1.OpenMCPJobList, error) {
	result := resourcev1alpha1.OpenMCPJobList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPJobClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.OpenMCPJob, error) {
	result := resourcev1alpha1.OpenMCPJob{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpjobs").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPJobClient) Create(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error) {
	result := resourcev1alpha1.OpenMCPJob{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpjobs").
		Body(ojob).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPJobClient) Update(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error) {
	result := resourcev1alpha1.OpenMCPJob{}
	err := c.restClient.
		Put().
		Name(ojob.Name).
		Namespace(c.ns).
		Resource("openmcpjobs").
		Body(ojob).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPJobClient) UpdateStatus(ojob *resourcev1alpha1.OpenMCPJob) (*resourcev1alpha1.OpenMCPJob, error) {
	result := resourcev1alpha1.OpenMCPJob{}
	err := c.restClient.
		Put().
		Name(ojob.Name).
		Namespace(c.ns).
		Resource("openmcpjobs").
		SubResource("status").
		Body(ojob).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPJobClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
