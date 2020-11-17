package resources

import (
	"encoding/json"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	"openmcp/openmcp/openmcp-snapshot/pkg/run/etcd/util"
)

// Dynamic resource
type Dynamic struct {
	apiCaller dynamic.ResourceInterface
}

// convertResourceObj : json String 을 obj 로 변환
func (dynamicResource Dynamic) convertResourceObj(resourceInfoJSON string) (*unstructured.Unstructured, error) {

	// jsonStr 에서 marshal 하기
	jsonBytes := []byte(resourceInfoJSON)

	// JSON 디코딩
	var unstructured *unstructured.Unstructured
	jsonEerr := json.Unmarshal(jsonBytes, &unstructured)
	if jsonEerr != nil {
		return nil, jsonEerr
	}
	return unstructured, nil
}

// CreateResourceForJSON : deploy 생성
func (dynamicResource Dynamic) CreateResourceForJSON(clientset dynamic.Interface, resourceInfoJSON string) (bool, error) {
	resourceInfo, convertErr := dynamicResource.convertResourceObj(resourceInfoJSON)
	if convertErr != nil {
		return false, convertErr
	}

	namespace := apiv1.NamespaceDefault
	gvk := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(namespace)

	//dynamic.apiCaller = clientset.AppsV1().Deployments(nameSpace)

	/*
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				//Name: "demo",
				Name:      resourceInfo.GetObjectMeta().GetName() + SnapshotTailName,
				Namespace: resourceInfo.ObjectMeta.Namespace,
				Labels:    resourceInfo.ObjectMeta.Labels,
			},

			Spec: resourceInfo.Spec,

			Spec: appsv1.DeploymentSpec{
				Replicas: resourceInfo.Spec.Replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: resourceInfo.Spec.Selector.MatchLabels,
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: resourceInfo.Spec.Template.ObjectMeta.Labels,
					},
					Spec: apiv1.PodSpec{
						Containers: resourceInfo.Spec.Template.Spec.Containers,
					},
				},
			},
		}
	*/
	/*
		resourceInfo := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "demo-deployment",
				},
				"spec": map[string]interface{}{
					"replicas": 2,
					"selector": map[string]interface{}{
						"matchLabels": map[string]interface{}{
							"app": "demo",
						},
					},
					"template": map[string]interface{}{
						"metadata": map[string]interface{}{
							"labels": map[string]interface{}{
								"app": "demo",
							},
						},

						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name":  "web",
									"image": "nginx:1.12",
									"ports": []map[string]interface{}{
										{
											"name":          "http",
											"protocol":      "TCP",
											"containerPort": 80,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	*/

	// Create Deployment
	fmt.Println("Creating deployment 2...")

	resourceInfo.SetResourceVersion("")
	oldName := resourceInfo.GetName()
	newName := oldName + "-snapshot"
	resourceInfo.SetName(newName)
	result, apiCallErr := dynamicResource.apiCaller.Create(resourceInfo, metav1.CreateOptions{})
	if apiCallErr != nil {
		return false, apiCallErr
	}
	fmt.Printf("Created deployment 2 %q.\n", result)

	return true, nil
}

// GetJSON : return json string
func (dynamicResource Dynamic) GetJSON(clientset dynamic.Interface, resourceName string, resourceNamespace string) (string, error) {
	gvk := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	dynamicResource.apiCaller = clientset.Resource(gvk).Namespace(resourceNamespace)

	fmt.Printf("Listing Resource in namespace %q:\n", resourceNamespace)

	result, apiCallErr := dynamicResource.apiCaller.Get(resourceName, metav1.GetOptions{})
	if apiCallErr != nil {
		return "", apiCallErr
	}

	return util.Obj2JsonString(result)
}
