package serviceregistry

import (
	"errors"
	"log"
	"openmcp/openmcp/omcplog"
	"sync"
)

var lock sync.RWMutex

// Common errors.
var (
	ErrServiceNotFound = errors.New("service name/version not found")
)

type Registry interface {
	Add(serviceName, endpoint string)
	Delete(serviceName string)
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(serviceName string) ([]string, error)
	EndpointCheck(serviceName string, endpoint string) bool
	Init()
}

// DefaultRegistry is a basic registry using the following format:
// {
//   "ServiceName": [
//       "cluster1",
//       "cluster2",
//     ],
// }

type DefaultRegistry map[string][]string

func (r DefaultRegistry) Init() {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Cluster Init")
	lock.RLock()
	for k := range r {
		delete(r, k)
	}
	lock.RUnlock()
}


func (r DefaultRegistry) EndpointCheck(serviceName string, endpoint string) bool {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Lookup")
	lock.RLock()
	targets, ok := r[serviceName]
	lock.RUnlock()
	if !ok {
		return true
	}
	for _, endpoints := range targets {
		omcplog.V(5).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] " + endpoints)
		if endpoint == endpoints {
			return false
		}
	}
	return true
}

func (r DefaultRegistry) Add(serviceName, endpoint string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Add")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[serviceName]
	if !ok {
		service = []string{}
		r[serviceName] = service
	}
	service = append(service, endpoint)
	r[serviceName] = append(r[serviceName], endpoint)
	omcplog.V(5).Info(service)
	omcplog.V(5).Info(r)
}

func (r DefaultRegistry) Lookup(serviceName string) ([]string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Lookup")
	lock.RLock()
	targets, ok := r[serviceName]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Failure")
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Delete(serviceName string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ServiceRegistry)] Function Delete")
	lock.Lock()
	defer lock.Unlock()

	_, ok := r[serviceName]
	if !ok {
		return
	}

	delete(r, serviceName)
}
