package util

import (
	"path/filepath"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getJobName(snapshotKey string) string {
	return "sns-" + snapshotKey + "-job"
}
func getExternalNfsPVCName(snapshotKey string) string {
	return "sns-" + snapshotKey + "-epvc"
}
func getExternalNfsPVName(snapshotKey string) string {
	return "sns-" + snapshotKey + "-epv"
}
func getPVCName(snapshotKey string) string {
	return "sns-" + snapshotKey + "-pvc"
}

//GetJobAPI  job 을 배포한다.
func GetJobAPI(snapshotKey string, cmd string) *batchv1.Job {
	//sns-KEY-job
	//
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getJobName(snapshotKey),
			Namespace: "openmcp",
		},
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					//"app": "openmcp",
				},
			},
			//ManualSelector: false,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "openmcp",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "docker-cmds",
							Image: "docker:19.03.8",
							Command: []string{
								"/bin/sh",
							},
							Args: []string{
								"-c", cmd,
								//"-c", "echo 'test'; docker commit -a gen -m commitTest 598bc3e5efe7 10.0.0.224:4999/test:2.0; do echo running; sleep 1000;done"
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/var/run",
									Name:      "docker-sock",
								},
								{
									MountPath: "/etc/docker/",
									Name:      "docker-cert",
								},
								{
									MountPath: "/storage",
									Name:      "external-nfs",
								},
								{
									MountPath: "/data",
									Name:      "target-volume",
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					//NodeSelector: map[string]string{
					//	labelName: "true",
					//},
					Volumes: []apiv1.Volume{
						{
							Name: "docker-sock",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/var/run",
								},
							},
						},
						{
							Name: "docker-cert",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/docker/",
								},
							},
						},
						{
							Name: "external-nfs",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: getExternalNfsPVCName(snapshotKey),
								},
							},
						},
						{

							Name: "target-volume",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: getPVCName(snapshotKey),
								},
							},
						},
					},
				},
			},
		},
	}
	return job
}

func GetExternalNfsPVCAPI(snapshotKey string) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getExternalNfsPVCName(snapshotKey),
			Namespace: "openmcp",
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("10Gi"),
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": getExternalNfsPVName(snapshotKey),
				},
			},
			VolumeName: getExternalNfsPVCName(snapshotKey),
		},
	}
	return pvc
}

func GetExternalNfsPVAPI(snapshotKey string, pvResource apiv1.PersistentVolume) *v1.PersistentVolume {
	pv := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			//Name: "demo",
			Name: getExternalNfsPVName(snapshotKey),
			Labels: map[string]string{
				"app": getExternalNfsPVName(snapshotKey),
			},
		},
		Spec: apiv1.PersistentVolumeSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Capacity: apiv1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resource.MustParse("10Gi"),
			},
			PersistentVolumeSource: apiv1.PersistentVolumeSource{
				NFS: &v1.NFSVolumeSource{
					Path:     getExternalNfsPVPathForVolumeSnapshot(pvResource.ClusterName, pvResource.Name),
					Server:   EXTERNAL_NFS,
					ReadOnly: false,
				},
			},
			PersistentVolumeReclaimPolicy: apiv1.PersistentVolumeReclaimRetain,
		},
	}
	return pv
}

func getExternalNfsPVPathForVolumeSnapshot(clusterName string, resourceName string) string {
	return getExternalNfsPVPath(clusterName, TypeVolumeSnapshot, resourceName)
}

func getExternalNfsPVPath(clusterName string, snapshotType string, resourceName string) string {
	externalNfsPVPath := strings.Join([]string{EXTERNAL_NFS_PATH_STORAGE, clusterName, snapshotType, resourceName}, string(filepath.Separator))
	return externalNfsPVPath
}

//PV 정보를 이용하여 PVC 정보를 만들어내는 부분. Label 을 통째로 watch 하도록 함.
func GetPVCAPI(snapshotKey string, pvResource apiv1.PersistentVolume) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPVCName(snapshotKey),
			Namespace: "openmcp",
		},
		Spec: apiv1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
			Resources: apiv1.ResourceRequirements{
				Requests: apiv1.ResourceList{
					v1.ResourceName(v1.ResourceStorage): resource.MustParse("10Gi"),
				},
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: pvResource.Labels,
			},

			//VolumeName: getPVCName(snapshotKey),
		},
	}
	return pvc
}
