package agent

//https://github.com/docker/go-docker/blob/master/ 참조하여 작성

//PushGlobalRegistryImage 는 node 에 등록된 이미지를 global repository 에 push 시키는 명령어이다.
func (r *RegistryNodeManager) PushGlobalRegistryImage(clusterName string, nodeName string, imageName string, Tag string) error {

	err := r.init(clusterName, nodeName)
	if err != nil {
		return err
	}
	return nil

	// node

}

/*
//PullGlobalRegistryImage 는 global repository 에 존재하는 이미지를 pull 하는 명령어이다.
func (r *RegistryNodeManager) PullGlobalRegistryImage(clusterName string, nodeName string, imageName string, Tag string) error {

}

//PushGlobalRegistryImageJob push 기능을 하는 job 을 node 에 배포하는 기능
func (r *RegistryNodeManager) PushGlobalRegistryImageJob(nodeName string) error {

		appName := r.getAppName(nodeName)
		labelName := r.getLabelName(nodeName)

		jobClient := r.clientset.AppsV1().(utils.ProjectNamespace)
		deployment, _ := deploymentsClient.Get(appName, metav1.GetOptions{})
		if deployment.ObjectMeta.Name != "" {
			fmt.Printf("deployment exist : " + deployment.ObjectMeta.Name)
			return nil
		}

		deployment = &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: appName,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": appName,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": appName,
						},
					},
					Spec: apiv1.PodSpec{
						Containers: []apiv1.Container{
							{
								Name:  "web",
								Image: "nginx:1.12",
								Ports: []apiv1.ContainerPort{
									{
										Name:          "http",
										Protocol:      apiv1.ProtocolTCP,
										ContainerPort: 80,
									},
								},
							},
						},
						NodeSelector: map[string]string{
							labelName: "true",
						},
					},
				},
			},
		}

		// Create Deployment
		fmt.Println("Creating deployment...")
		result, err := deploymentsClient.Create(deployment)
		if err != nil {
			return err
		}
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}
*/
