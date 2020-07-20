package v1alpha1

import (
	"time"
	"k8s.io/klog"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"
	"openmcp/openmcp/openmcp-scheduler/pkg/framework/plugins/predicates"
	"openmcp/openmcp/openmcp-scheduler/pkg/framework/plugins/priorities"
)

type openmcpFramework struct {
	filterPlugins		[]OpenmcpFilterPlugin
	scorePlugins		[]OpenmcpScorePlugin
}


// Interface checks
//The appearance of the blank identifier in this construct indicates 
// that the declaration exists only for the type checking, not to create a variable.
var _ OpenmcpFramework = &openmcpFramework{}

// have to change.. argument should be config_file
func NewFramework() OpenmcpFramework{
	klog.Info("Create NewFramework...")

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
	}

	return f
}

func (f *openmcpFramework) RunFilterPluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster) OpenmcpClusterFilteredStatus{
	result := make(map[string]bool)

	klog.Info("[FILTERING] Run FilterPlugins on Cluster Level")

	if clusters == nil{
		return nil
	}

	for _, cluster := range clusters {
		result[cluster.ClusterName] = true

		for _, pl := range f.filterPlugins {
			startTime := time.Now()
			isFiltered := pl.Filter(pod, cluster)
			elapsedTime := time.Since(startTime)

			klog.Infof("[%v] %v\tResult(T/F): %v\t%v", 
						cluster.ClusterName, pl.Name(), isFiltered, elapsedTime)

			// Update the result of this cluster
			result[cluster.ClusterName] = result[cluster.ClusterName] && isFiltered
			if !result[cluster.ClusterName]{
				break
			}
		}
	}
	return result
}

func (f *openmcpFramework) RunScorePluginsOnClusters(pod *ketiresource.Pod, clusters map[string]*ketiresource.Cluster) OpenmcpPluginToClusterScores{
	result := make(map[string]OpenmcpPluginScoreList)

	klog.Info("[SCORING] Run ScorePlugins on Cluster Level")

	for _, cluster := range clusters {
		result[cluster.ClusterName] = make([]OpenmcpPluginScore, 0)

		for _, pl := range f.scorePlugins {
			startTime := time.Now()

			plScore := OpenmcpPluginScore{
				Name:	pl.Name(),
				Score:	pl.Score(pod, cluster),
			}

			elapsedTime := time.Since(startTime)

			klog.Infof("[%v] %v\tResult: %v\t%v", 
						cluster.ClusterName, pl.Name(), plScore, elapsedTime)

			// Update the result of this cluster
			result[cluster.ClusterName] = append(result[cluster.ClusterName], plScore)
		}
	}
	return result
}

func (f *openmcpFramework) HasFilterPlugins() bool {
	return len(f.filterPlugins) > 0
}

func (f *openmcpFramework) HasScorePlugins() bool {
	return len(f.scorePlugins) > 0
}
