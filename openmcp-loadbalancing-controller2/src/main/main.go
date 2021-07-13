/*
Copyright 2018 The Multicluster-Controller Authors.

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

package main

import (
	//"sync"
	//"time"

	"openmcp/openmcp/openmcp-loadbalancing-controller2/src/controller"
	"openmcp/openmcp/openmcp-loadbalancing-controller2/src/controller/DestinationRule/DestinationRuleWeight"
	"openmcp/openmcp/openmcp-loadbalancing-controller2/src/reverseProxy"
	"sync"
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go reverseProxy.ReverseProxy()
	go controller.ServiceMeshController()
	//go OpenMCPVirtualService.SyncWeight()
	time.Sleep(time.Second * 2)
	go DestinationRuleWeight.AnalyticWeight()

	wg.Wait()

}
