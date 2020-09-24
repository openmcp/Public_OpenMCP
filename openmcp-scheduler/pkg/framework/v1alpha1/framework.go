package v1alpha1

import (
	// "openmcp/openmcp/omcplog"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/openmcp-scheduler/pkg/framework/plugins/predicates"
	"openmcp/openmcp/openmcp-scheduler/pkg/framework/plugins/priorities"
	"openmcp/openmcp/openmcp-scheduler/pkg/protobuf"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
)

type openmcpFramework struct {
	filterPlugins     []OpenmcpFilterPlugin
	scorePlugins      []OpenmcpScorePlugin
	prefilterPlugins  []OpenmcpPreFilterPlugin
	prescorePlugins   []OpenmcpPreScorePlugin
	postfilterPlugins []OpenmcpPostFilterPlugin
	IspreScore        bool
	preScore          int64
	betweenScores     int64
	preselectedName   string
	preClusterName    string
}

// The appearance of the blank identifier in this construct indicates
// that the declaration exists only for the type checking, not to create a variable.
var _ OpenmcpFramework = &openmcpFramework{}

func (f *openmcpFramework) EndPod() {
	f.IspreScore = false
	f.preScore = 0
	f.betweenScores = 0
	f.preselectedName = ""
	f.preClusterName = ""
}
func NewFramework(grpcClient protobuf.RequestAnalysisClient) OpenmcpFramework {

	f := &openmcpFramework{
		filterPlugins: []OpenmcpFilterPlugin{
			&predicates.MatchClusterSelector{},
			&predicates.PodFitsResources{},
			&predicates.CheckNeededResources{},
			&predicates.MatchClusterAffinity{},
			&predicates.PodFitsHostPorts{},
			&predicates.NoDiskConflict{},
		},
		scorePlugins: []OpenmcpScorePlugin{
			&priorities.MostRequested{},
			&priorities.DominantResource{},
			&priorities.RequestedToCapacityRatio{},
			&priorities.BalancedNetworkAllocation{},
			&priorities.QosPriority{},
		},
		prefilterPlugins: []OpenmcpPreFilterPlugin{
			&predicates.PodFitsResources{},
			&predicates.MatchClusterAffinity{},
		},
		prescorePlugins: []OpenmcpPreScorePlugin{
			&priorities.MostRequested{},
			&priorities.DominantResource{},
			&priorities.RequestedToCapacityRatio{},
			&priorities.BalancedNetworkAllocation{},
			&priorities.QosPriority{},
		},
		postfilterPlugins: []OpenmcpPostFilterPlugin{
			&predicates.PodFitsResources{},
		},
	}
	return f
}

func (f *openmcpFramework) RunFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, postpods []*ketiresource.Pod) OpenmcpClusterFilteredStatus {
	result := make(map[string]bool)

	if clusters == nil {
		return nil
	}
	for _, cluster := range clusters {
		cluster.PreFilterA = false
		cluster.PreFilter = false
		result[cluster.ClusterName] = true
		for _, pl := range f.prefilterPlugins {
			pl.PreFilter(pod, cluster)
		}
		if cluster.PreFilter == false || cluster.PreFilterA == false {
			//omcplog.V(0).Info("cluster differ", cluster.ClusterName, cluster.PreFilter, cluster.PreFilterA)
			result[cluster.ClusterName] = false
			continue
		}
		for _, pl := range f.filterPlugins {

			isFiltered := pl.Filter(pod, cluster)

			result[cluster.ClusterName] = result[cluster.ClusterName] && isFiltered

			if !result[cluster.ClusterName] {
				break
			}
		}
	}
	omcplog.V(0).Info("filter result=> ", result)
	return result
}
func selectCluster(scoreResult OpenmcpPluginToClusterScores) string {
	var selectedCluster string
	var maxScore int64

	for clusterName, scoreList := range scoreResult {
		var clusterScore int64
		for _, score := range scoreList {
			clusterScore += score.Score
		}

		if clusterScore > maxScore {
			selectedCluster = clusterName
			maxScore = clusterScore

		}
	}

	return selectedCluster
}

// func (f *openmcpFramework) RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, replicas int32) OpenmcpPluginToClusterScores {
func (f *openmcpFramework) RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster, allclusters map[string]*ketiresource.Cluster, replicas int32) string {

	if !f.IspreScore {
		f.preScore = 0
		preresult := make(map[string]OpenmcpPluginScoreList)

		for _, cluster := range clusters {
			preresult[cluster.ClusterName] = make([]OpenmcpPluginScore, 0)
			for _, pl := range f.prescorePlugins {
				scoring := pl.PreScore(pod, cluster, false)
				//omcplog.V(0).Infof("[PreScore]")
				transScore := OpenmcpPluginScore{
					Name:  pl.Name(),
					Score: scoring,
				}
				f.preScore += scoring

				preresult[cluster.ClusterName] = append(preresult[cluster.ClusterName], transScore)
			}

		}
		f.IspreScore = true

		pr := selectCluster(preresult)
		f.preselectedName = pr
		f.preClusterName = pr
		//omcplog.V(0).Info("one preScore", f.preselectedName)

		return pr
	}
	if f.IspreScore && f.preselectedName != "" {
		for _, pl := range f.prescorePlugins {
			pl.PreScore(pod, allclusters[f.preselectedName], true)
		}
		f.preselectedName = ""
	}
	result := make(map[string]OpenmcpPluginScoreList)
	for _, cluster := range clusters {

		result[cluster.ClusterName] = make([]OpenmcpPluginScore, 0)

		for _, pl := range f.scorePlugins {

			plScore := OpenmcpPluginScore{
				Name:  pl.Name(),
				Score: pl.Score(pod, cluster, replicas, f.preClusterName),
			}
			// omcplog.V(0).Infof("[%v] %-22vScore:%5v", cluster.ClusterName, pl.Name(), plScore.Score)
			// Update the result of this cluster
			result[cluster.ClusterName] = append(result[cluster.ClusterName], plScore)
		}
	}
	pr := selectCluster(result)
	f.preClusterName = pr
	return pr
}
func (f *openmcpFramework) HasFilterPlugins() bool {
	return len(f.filterPlugins) > 0
}

func (f *openmcpFramework) HasScorePlugins() bool {
	return len(f.scorePlugins) > 0
}
