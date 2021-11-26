/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logic

import (
	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
)

// VpaPreProcessor processes the VPAs before applying default .
type VpaPreProcessor interface {
	Process(vpa *vpa_types.VerticalPodAutoscaler, isCreate bool) (*vpa_types.VerticalPodAutoscaler, error)
}

// noopVpaPreProcessor leaves pods unchanged when processing
type noopVpaPreProcessor struct{}

// Process leaves the pod unchanged
func (p *noopVpaPreProcessor) Process(vpa *vpa_types.VerticalPodAutoscaler, isCreate bool) (*vpa_types.VerticalPodAutoscaler, error) {
	return vpa, nil
}

// NewDefaultVpaPreProcessor creates a VpaPreProcessor that leaves VPAs unchanged and returns no error
func NewDefaultVpaPreProcessor() VpaPreProcessor {
	return &noopVpaPreProcessor{}
}
