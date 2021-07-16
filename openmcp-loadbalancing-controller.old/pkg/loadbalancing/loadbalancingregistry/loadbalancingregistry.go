package loadbalancingregistry

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
	Add(host, path, endpoint string)                // Add an endpoint to our registry
	Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(host, path string) (string, error)       // Return the endpoint list for the given service name/version
	IngressDelete(host, path string)
	Init()
}


// DefaultRegistry is a basic registry using the following format:
// {
//   "Host": {
//     "Path": [
//       "serviceName",
//     ],
//   },
// }

type DefaultRegistry map[string]map[string]string


func (r DefaultRegistry) Init() {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Cluster Init")
	lock.RLock()
	for k := range r {
		delete(r, k)
	}
	lock.RUnlock()
}

func (r DefaultRegistry) Lookup(host string, path string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Lookup")
	lock.RLock()
	target, ok := r[host][path]
	lock.RUnlock()
	if !ok {
		return "", ErrServiceNotFound
	}
	return target, nil
}


func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Failure")
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Add(host, path, endpoint string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Add")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		service = map[string]string{}
		r[host] = service
	}
	service[path] = endpoint
}

func (r DefaultRegistry) Delete(host, path, endpoint string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function Delete")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		return
	}
	omcplog.V(5).Info(service)
}


func (r DefaultRegistry) IngressDelete(host, path string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(LoadbalancingRegistry)] Function IngressDelete")
	lock.Lock()
	defer lock.Unlock()

	service, ok := r[host]
	if !ok {
		return
	}
	delete(service, path)

	if len(r[host]) == 0 {
		delete(r, host)
	}
}
