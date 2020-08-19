package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"resource-controller/apis/keti/v1alpha1"
	"resource-controller/apis"
)

type ExampleV1Alpha1Interface interface {
	OpenMCPDeployment(namespace string) OpenMCPDeploymentInterface
	OpenMCPHybridAutoScaler(namespace string) OpenMCPHybridAutoScalerInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	apis.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) OpenMCPDeployment(namespace string) OpenMCPDeploymentInterface {
	return &OpenMCPDeploymentClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) OpenMCPHybridAutoScaler(namespace string) OpenMCPHybridAutoScalerInterface {
	return &OpenMCPHybridAutoScalerClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}