package ingressregistry

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

// Registry is an interface used to lookup the target host
// for a given service name / version pair.
type Registry interface {
	Add(ingressName, url string)
	Delete(ingressName string)
	//Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(ingressName string) ([]string, error)
	CheckURL(url string) (bool, error)
}

// DefaultRegistry is a basic registry using the following format:
// {
//   "IngressName": [
//       "keti.test.com/test",
//       "lb_test.com/service",
//     ],
// }

//type DefaultRegistry map[string]map[string]map[string]stringzmgma
type DefaultRegistry map[string][]string

// Lookup return the endpoint list for the given service name/version.

func (r DefaultRegistry) Add(ingressName, url string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(IngressRegistry)] Function Add")
	lock.Lock()
	defer lock.Unlock()

	omcplog.V(5).Info("[OpenMCP Loadbalancing Controller(IngressRegistry)] " + ingressName)

	service, ok := r[ingressName]
	if !ok {
		service = []string{}
		r[ingressName] = service
	}
	service = append(service, url)
	r[ingressName] = append(r[ingressName], url)
	omcplog.V(5).Info(r)
}

func (r DefaultRegistry) Lookup(ingressName string) ([]string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] Function Lookup")
	lock.RLock()
	targets, ok := r[ingressName]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

func (r DefaultRegistry) CheckURL(url string) (bool, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] Function CheckURL")
	lock.RLock()
	targets := r
	lock.RUnlock()

	for _, ingressName := range targets {
		for _, ingressURL := range ingressName {
			omcplog.V(5).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] " + ingressURL)
			if ingressURL == url {
				return true, nil
			}
		}
	}
	return false, ErrServiceNotFound
}

func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] Function Failure")
	// Would be used to remove an endpoint from the rotation, log the failure, etc.
	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Delete(ingressName string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] Function Delete")
	lock.Lock()
	defer lock.Unlock()

	_, ok := r[ingressName]
	if !ok {
		return
	}

	delete(r, ingressName)
}

//// Delete removes the given endpoit for the service name/version.
//func (r DefaultRegistry) Delete(host, path, endpoint string) {
//	fmt.Println("----Delete----")
//	lock.Lock()
//	defer lock.Unlock()
//
//	service, ok := r[host]
//	if !ok {
//		return
//	}
//
//begin:
//	for i, svc := range service[path] {
//		if svc == endpoint {
//			copy(service[path][i:], service[path][i+1:])
//			service[path] = service[path][:len(service[path])-1]
//			goto begin
//		}
//	}
//}
