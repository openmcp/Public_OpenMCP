package snapshot

import (
	"context"
	"fmt"
	snapshotv1alpha1 "openmcp/openmcp/apis/snapshot/v1alpha1"
	"openmcp/openmcp/omcplog"

	corev1 "k8s.io/api/core/v1"
)

func (r *reconciler) makeStatusRun(instance *snapshotv1alpha1.Snapshot, status corev1.ConditionStatus, description string, elapsedTime string, err error) {

	if elapsedTime == "" {
		elapsedTime = "-"
	}

	instance.Status.ElapsedTime = elapsedTime
	instance.Status.Status = status
	instance.Status.Description = description
	instance.Status.ConditionProgress = fmt.Sprintf("%f", float64(r.progressCurrent)/float64(r.progressMax)*100) + "%"

	omcplog.V(3).Info(instance.Status)

	omcplog.V(3).Info("progressCurrent : ", r.progressCurrent)
	omcplog.V(3).Info("progressMax : ", r.progressMax)

	omcplog.V(3).Info("isVolumeSnapshot : ", instance.Status.IsVolumeSnapshot)
	omcplog.V(3).Info("elapsedTime : ", instance.Status.ElapsedTime)
	omcplog.V(3).Info("Status : ", instance.Status.Status)
	omcplog.V(3).Info("Description : ", instance.Status.Description)
	omcplog.V(3).Info("ConditionProgress : ", instance.Status.ConditionProgress)
	omcplog.V(3).Info("---SnapshotSources---")
	omcplog.V(3).Info(instance.Status.SnapshotSources)

	err = r.live.Status().Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
	err = r.live.Update(context.TODO(), instance)
	if err != nil {
		omcplog.V(3).Info(err, "-----------")
	}
}
