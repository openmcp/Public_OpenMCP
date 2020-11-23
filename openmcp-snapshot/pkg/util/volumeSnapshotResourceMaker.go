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

// RunType : Snapshot, SnapshotRestore 동작 구분
type RunType string

const (
	// RunTypeSnapshot : 스냅샷 타입
	RunTypeSnapshot RunType = "sns"
	// RunTypeSnapshotRestore : 스냅샷 복구 타입
	RunTypeSnapshotRestore RunType = "snr"
)

func getJobName(snapshotKey string, runType RunType) string {
	return string(runType) + "-" + snapshotKey + "-job"
}
func getExternalNfsPVCName(snapshotKey string, runType RunType) string {
	return string(runType) + "-" + snapshotKey + "-epvc"
}
func getExternalNfsPVName(snapshotKey string, runType RunType) string {
	return string(runType) + "-" + snapshotKey + "-epv"
}
func getPVCName(snapshotKey string, runType RunType) string {
	return string(runType) + "-" + snapshotKey + "-pvc"
}
func getPVName(snapshotKey string, runType RunType) string {
	return string(runType) + "-" + snapshotKey + "-pv"
}

//GetJobAPI  job 을 배포한다.
func GetJobAPI(snapshotKey string, cmd string, runType RunType) *batchv1.Job {
	//sns-KEY-job
	//
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getJobName(snapshotKey, runType),
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
							Name:  "ubuntu-cmds",
							Image: "ubuntu:16.04",
							Command: []string{
								"/bin/sh",
							},
							Args: []string{
								"-c", cmd,
								//"-c", "echo 'test'; docker commit -a gen -m commitTest 598bc3e5efe7 10.0.0.224:4999/test:2.0; do echo running; sleep 1000;done"
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									MountPath: "/storage",
									Name:      "external-nfs",
								},
								{
									MountPath: "/data",
									Name:      "target-volume",
								},
								{
									MountPath: "/etc/localtime",
									Name:      "tz-seoul",
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
							Name: "external-nfs",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: getExternalNfsPVCName(snapshotKey, runType),
								},
							},
						},
						{

							Name: "target-volume",
							VolumeSource: apiv1.VolumeSource{
								PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
									ClaimName: getPVCName(snapshotKey, runType),
								},
							},
						},
						{
							Name: "tz-seoul",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/etc/localtime",
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

func GetExternalNfsPVCAPI(snapshotKey string, runType RunType) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getExternalNfsPVCName(snapshotKey, runType),
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
					"app": getExternalNfsPVName(snapshotKey, runType),
				},
			},
		},
	}
	return pvc
}

func GetExternalNfsPVAPI(snapshotKey string, pvResource apiv1.PersistentVolume, runType RunType) (*v1.PersistentVolume, string) {
	pv := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			//Name: "demo",
			Name: getExternalNfsPVName(snapshotKey, runType),
			Labels: map[string]string{
				"app": getExternalNfsPVName(snapshotKey, runType),
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
	mountPath := makeMountPath(pvResource.ClusterName, pvResource.Name)
	return pv, mountPath
}
func makeMountPath(clusterName string, resourceName string) string {
	mountPath := strings.Join([]string{clusterName, TypeVolumeSnapshot, resourceName}, string(filepath.Separator))
	return mountPath
}
func getExternalNfsPVPathForVolumeSnapshot(clusterName string, resourceName string) string {
	return getExternalNfsPVPath(clusterName, TypeVolumeSnapshot, resourceName)
}

func getExternalNfsPVPath(clusterName string, snapshotType string, resourceName string) string {
	//externalNfsPVPath := strings.Join([]string{EXTERNAL_NFS_PATH_STORAGE, clusterName, snapshotType, resourceName}, string(filepath.Separator))
	externalNfsPVPath := EXTERNAL_NFS_PATH_STORAGE
	return externalNfsPVPath
}

//PV 정보를 이용하여 PVC 정보를 만들어내는 부분. Label 을 snapshot Key로 둠
func GetPVCAPI(snapshotKey string, pvResource apiv1.PersistentVolume, runType RunType) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPVCName(snapshotKey, runType),
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
					"name": getPVName(snapshotKey, runType),
				},
			},

			//VolumeName: getPVCName(snapshotKey),
		},
	}
	return pvc
}

func GetPVAPI(snapshotKey string, pvResource apiv1.PersistentVolume, runType RunType) *v1.PersistentVolume {
	pv := &apiv1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			//Name: "demo",
			Name: getPVName(snapshotKey, runType),
			Labels: map[string]string{
				"name": getPVName(snapshotKey, runType),
			},
		},
		Spec: apiv1.PersistentVolumeSpec{},
	}
	pv.Spec.Capacity = pvResource.Spec.Capacity
	pv.Spec.PersistentVolumeSource = pvResource.Spec.PersistentVolumeSource
	pv.Spec.PersistentVolumeReclaimPolicy = pvResource.Spec.PersistentVolumeReclaimPolicy
	pv.Spec.AccessModes = pvResource.Spec.AccessModes

	return pv
}

func GetPVAPIOri(snapshotKey string, pvResource apiv1.PersistentVolume) *v1.PersistentVolume {
	pv := &pvResource

	pv.ResourceVersion = ""
	pv.Spec.ClaimRef = nil

	return pv
}

//PV 정보를 이용하여 PVC 정보를 만들어내는 부분. Label 을 통째로 watch 하도록 함.
func GetPVCAPIOri(snapshotKey string, pvResource apiv1.PersistentVolume, runType RunType) *v1.PersistentVolumeClaim {
	pvc := &v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPVCName(snapshotKey, runType),
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
