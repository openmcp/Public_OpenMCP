package v1alpha1

import (
	"container/list"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

const (
	MaxNodeScore int64 = 100
	MinNodeScore int64 = 0
)

// OpenmcpClusterScoreList declares a list of plugins and their scores.
type OpenmcpPluginScoreList []OpenmcpPluginScore

// OpenmcpClusterScore is a struct with plugin name and score
type OpenmcpPluginScore struct {
	Name  string
	Score int64
}

// OpenmcpPluginToClusterScores declare map from cluster name to its OpenmcpClusterScoreList
type OpenmcpPluginToClusterScores map[string]OpenmcpPluginScoreList

// OpenmcpPluginFilteredStatus declare map from cluster name to its filtering result
type OpenmcpClusterFilteredStatus map[string]bool
type OpenmcpClusterPostFilteredStatus map[string]bool
type OpenmcpFramework interface {
	// RunFilterPluginsOnClusters runs the set of configured filtering plugins.
	// It returns a map that stores for each filtering plugin name the corresponding
	RunFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster) OpenmcpClusterFilteredStatus
	EraseFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, request map[string]int32) string
	// RunScorePluginsOnClusters runs the set of configured scoring plugins.
	// It returns a map that stores for each
	RunPostFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, postpods *list.List) OpenmcpClusterPostFilteredStatus
	RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, allclusters map[string]*ketiresource.Cluster, replicas int32) string
	//RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, replicas int32) OpenmcpPluginToClusterScores
	EndPod()
}

// OpenmcpPlugin is the parent type for all the scheduling framework plugins
type OpenmcpPlugin interface {
	Name() string
}

// OpenmcpFilterPlugin is an interface for Filter plugins.
// This concept used to be called 'predicate' in the original scheduler.
// This plugins should return "true" if the pod can be deployed into the cluster and
// return "false" if the pod can not be deployed into the cluster
type OpenmcpFilterPlugin interface {
	OpenmcpPlugin
	Filter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool
}

// type EraseFilterPlugin interface {
// 	OpenmcpPlugin
// 	PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64
// }
type OpenmcpScorePlugin interface {
	OpenmcpPlugin
	Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64
	PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64
}

type OpenmcpEraseScorePlugin interface {
	OpenmcpPlugin
	Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64
	PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64
}

// OpenmcpPreFilterPlugin and OpenmcpPreScore are interfaces for PreFilter and PreScore
//this concept used to ..
type OpenmcpPreFilterPlugin interface {
	OpenmcpPlugin
	PreFilter(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) bool
}

//
type OpenmcpPostFilterPlugin interface {
	OpenmcpPlugin
	PostFilter(newPod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, postpods *list.List) (bool, error)
}
