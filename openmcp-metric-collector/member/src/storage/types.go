// Copyright 2018 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
)

type Collection struct {
	Metricsbatchs []MetricsBatch
	ClusterName string
	LatencyTime   *timestamppb.Timestamp
}

// MetricsBatch is a single batch of pod, container, and node metrics from some source.
type MetricsBatch struct {
	IP   string
	Node NodeMetricsPoint
	Pods []PodMetricsPoint
}

// NodeMetricsPoint contains the metrics for some node at some point in time.
type NodeMetricsPoint struct {
	Name string
	MetricsPoint
}

// PodMetricsPoint contains the metrics for some pod's containers.
type PodMetricsPoint struct {
	Name      string
	Namespace string
	MetricsPoint
	Containers []ContainerMetricsPoint
}

// ContainerMetricsPoint contains the metrics for some container at some point in time.
type ContainerMetricsPoint struct {
	Name string
	MetricsPoint
}

// MetricsPoint represents the a set of specific metrics at some point in time.
type MetricsPoint struct {
	Timestamp time.Time

	// Cpu
	CPUUsageNanoCores resource.Quantity

	// Memory
	MemoryUsageBytes      resource.Quantity
	MemoryAvailableBytes  resource.Quantity
	MemoryWorkingSetBytes resource.Quantity

	// Network
	NetworkRxBytes resource.Quantity
	NetworkTxBytes resource.Quantity

	// Fs
	FsAvailableBytes resource.Quantity
	FsCapacityBytes  resource.Quantity
	FsUsedBytes      resource.Quantity
}
