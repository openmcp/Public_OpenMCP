package v1alpha1

import (
	"context"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type DestinationRuleInterface interface {
	List(opts metav1.ListOptions) (*v1alpha3.DestinationRuleList, error)
	Get(name string) (*v1alpha3.DestinationRule, error)
	Create(destinationrule *v1alpha3.DestinationRule) (*v1alpha3.DestinationRule, error)
	Update(destinationrule *v1alpha3.DestinationRule) (*v1alpha3.DestinationRule, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type DestinationRuleClient struct {
	restClient rest.Interface
	ns         string
}

func (c *DestinationRuleClient) List(opts metav1.ListOptions) (*v1alpha3.DestinationRuleList, error) {
	result := v1alpha3.DestinationRuleList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("destinationrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DestinationRuleClient) Get(name string) (*v1alpha3.DestinationRule, error) {
	result := v1alpha3.DestinationRule{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("destinationrules").
		Name(name).
		//VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DestinationRuleClient) Update(destinationrule *v1alpha3.DestinationRule) (*v1alpha3.DestinationRule, error) {
	result := v1alpha3.DestinationRule{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("destinationrules").
		Body(destinationrule).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DestinationRuleClient) Create(destinationrule *v1alpha3.DestinationRule) (*v1alpha3.DestinationRule, error) {
	result := v1alpha3.DestinationRule{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("destinationrules").
		Body(destinationrule).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DestinationRuleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("destinationrule").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
