package v1alpha1

import (
	"context"
	clusterv1alpha1 "openmcp/openmcp/apis/cluster/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type OpenMCPClusterInterface interface {
	List(opts metav1.ListOptions) (*clusterv1alpha1.OpenMCPClusterList, error)
	Get(name string, options metav1.GetOptions) (*clusterv1alpha1.OpenMCPCluster, error)
	Create(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error)
	Update(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error)
	UpdateStatus(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type OpenMCPClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *OpenMCPClusterClient) List(opts metav1.ListOptions) (*clusterv1alpha1.OpenMCPClusterList, error) {
	result := clusterv1alpha1.OpenMCPClusterList{}
	//c.restClient.Get().Namespace(c.ns).Resource()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPClusterClient) Get(name string, opts metav1.GetOptions) (*clusterv1alpha1.OpenMCPCluster, error) {
	result := clusterv1alpha1.OpenMCPCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpclusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPClusterClient) Create(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error) {
	result := clusterv1alpha1.OpenMCPCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("openmcpclusters").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPClusterClient) Update(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error) {
	result := clusterv1alpha1.OpenMCPCluster{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpclusters").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}
func (c *OpenMCPClusterClient) UpdateStatus(cluster *clusterv1alpha1.OpenMCPCluster) (*clusterv1alpha1.OpenMCPCluster, error) {
	result := clusterv1alpha1.OpenMCPCluster{}
	err := c.restClient.
		Put().
		Name(cluster.Name).
		Namespace(c.ns).
		Resource("openmcpclusters").
		SubResource("status").
		Body(cluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *OpenMCPClusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("openmcpclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
