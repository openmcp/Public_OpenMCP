package v1alpha1

import (
	resourcev1alpha1 "openmcp/openmcp/apis/resource/v1alpha1"

	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleV1Alpha1Interface interface {
	OpenMCPDeployment(namespace string) OpenMCPDeploymentInterface
	OpenMCPHybridAutoScaler(namespace string) OpenMCPHybridAutoScalerInterface
	OpenMCPPolicy(namespace string) OpenMCPPolicyInterface
	OpenMCPService(namespace string) OpenMCPServiceInterface
	OpenMCPIngress(namespace string) OpenMCPIngressInterface
	OpenMCPNamespace(namespace string) OpenMCPNamespaceInterface
	VirtualService(namespace string) VirtualServiceInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	//apis.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: resourcev1alpha1.GroupName, Version: resourcev1alpha1.GroupVersion}
	config.APIPath = "/apis"
	//config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func NewIstioForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	//apis.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha3.GroupName, Version: "v1alpha3"}
	config.APIPath = "/apis"
	//config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
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
func (c *ExampleV1Alpha1Client) OpenMCPPolicy(namespace string) OpenMCPPolicyInterface {
	return &OpenMCPPolicyClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) OpenMCPService(namespace string) OpenMCPServiceInterface {
	return &OpenMCPServiceClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) OpenMCPIngress(namespace string) OpenMCPIngressInterface {
	return &OpenMCPIngressClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) OpenMCPNamespace(namespace string) OpenMCPNamespaceInterface {
	return &OpenMCPNamespaceClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) VirtualService(namespace string) VirtualServiceInterface {
	return &VirtualServiceClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
