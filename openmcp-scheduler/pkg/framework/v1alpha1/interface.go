package v1alpha1

import (
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
	Name	string
	Score	int64
}

// OpenmcpPluginToClusterScores declare map from cluster name to its OpenmcpClusterScoreList
type OpenmcpPluginToClusterScores map[string]OpenmcpPluginScoreList

// OpenmcpPluginFilteredStatus declare map from cluster name to its filtering result
type OpenmcpClusterFilteredStatus map[string]bool

type OpenmcpFramework interface {
	// RunFilterPluginsOnClusters runs the set of configured filtering plugins. 
	// It returns a map that stores for each filtering plugin name the corresponding 
	RunFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster) OpenmcpClusterFilteredStatus
	// RunScorePluginsOnClusters runs the set of configured scoring plugins.
	// It returns a map that stores for each 
	RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster) OpenmcpPluginToClusterScores
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

type OpenmcpScorePlugin interface {
	OpenmcpPlugin

	Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster) int64
}
