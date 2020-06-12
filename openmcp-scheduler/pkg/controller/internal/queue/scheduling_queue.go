// This file contains structures that implement scheduling queue types.
// Scheduling queues hold OpenMCPDeployments waiting to be scheduled.
// This file implements a priority queue which has two sub queues.
// (1) activeQ : holds OpenMCPDeployments that are being considered for sheduling.
// (2) unschedulableQ : holds OpenMCPDeployments that are already tried and 
// are determined to be unschedulable.

package queue

import (
	"sync"
	"k8s.io/klog"
	v1 "k8s.io/api/core/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/scheduler/internal/heap"
	ketiv1alpha1 "openmcpscheduler/pkg/apis/keti/v1alpha"
)

type SchedulingQueue interface {
	Add(dep *ketiv1alpha1.OpenMCPDeployment) error
	AddUnschedulableIfNotPresent(dep *ketiv1alpha1.OpenMCPDeployment, podSchedulingCycle int64) error
	SchedulingCycle() int64
	Pop() // have to change
	Update(oldDep, newDep *ketiv1alpha1.OpenMCPDeployment) error
	Delete(dep *ketiv1alpha1.OpenMCPDeployment) error
	MoveAllToActiveOr... // have to change
	AssignedDepAdded(dep *ketiv1alpha1.OpenMCPDeployment)
	AssignedDepUpdated(dep *ketiv1alpha1.OpenMCPDeployment)
	NominatedPodsForNode(nodeName string) []*v1.Pod
	NominatedNodesForCluster(clusterName string) []*v1.Node
	PendingDeployments() []*ketiv1alpha1.OpenMCPDeployment
	Close()
	// check one more time...
	UpdateNominatedPodForNode(pod *v1.Pod, nodeName string)
	UpdateNominatedNodeForCluster(node *v1.Node, clusterName string)
	// have to change....
	DeleteNominatedPodIfExists(po
	// NumUnschedulablePods returns the number of unschedulable Deployments exist in the SchedulingQueue.
	NumUnschedulableDeployments() int
}

// PriorityQueue implements a scheduling queue.
// The head of PriorityQueue is the highest prirotiy pending pod. 
// The structure has two sub queues. 
type PriorityQueue struct {
	stop <- chan struct{}

	lock sync.RWMutex
	cond sync.Cond

	activeQ *heap.Heap
	unschedulableQ *heap.Heap
	nominatedDeployments *nominatedDeploymentMap
	schedulingCyclce int64
	closed bool
}

// NewSchedulingQueue initializes a priority queue as a new scheduling queue.
func NewSchedulingQueue(stop <-chan struct{}) SchedulingQueue {
	return newPriorityQueue(stop)
}

func NewPriorityQueue(stop <-chan struct{}) *PriorityQueue {
	pq := &PriorityQueue{
		stop:		stop,
		activeQ:	heap.NewWith
	}
	pq.cond.L = &pq.lock
	pq.run()

	return pq
}
