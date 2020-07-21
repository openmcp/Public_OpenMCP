package serviceregistry

import (
	"errors"
	"fmt"
	"log"
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
	Add(serviceName, endpoint string)
	Delete(serviceName string)
	//Delete(host, path, endpoint string)             // Remove an endpoint to our registry
	Failure(host, path, endpoint string, err error) // Mark an endpoint as failed.
	Lookup(serviceName string) ([]string, error)
	EndpointCheck(serviceName string, endpoint string) bool
}

// DefaultRegistry is a basic registry using the following format:
// {
//   "ServiceName": [
//       "cluster1",
//       "cluster2",
//     ],
// }

//type DefaultRegistry map[string]map[string]map[string]stringzmgma
type DefaultRegistry map[string][]string

// Lookup return the endpoint list for the given service name/version.

func (r DefaultRegistry) EndpointCheck(serviceName string, endpoint string) bool {
	fmt.Println("*****Ingress Lookup*****")
	lock.RLock()
	targets, ok := r[serviceName]
	lock.RUnlock()
	if !ok {
		return true
	}
	fmt.Println(endpoint)
	for _, endpoints := range targets {
		fmt.Println(endpoints)
		if endpoint == endpoints {
			return false
		}
	}
	return true
}

func (r DefaultRegistry) Add(serviceName, endpoint string) {
	fmt.Println("*****Service Name Add*****")
	lock.Lock()
	defer lock.Unlock()

	fmt.Println(r)
	fmt.Println(serviceName)
	fmt.Println(endpoint)

	service, ok := r[serviceName]
	fmt.Println(service)
	if !ok {
		service = []string{}
		r[serviceName] = service
	}
	service = append(service, endpoint)
	r[serviceName] = append(r[serviceName], endpoint)
	fmt.Println(service)
	fmt.Println(r)
}

func (r DefaultRegistry) Lookup(serviceName string) ([]string, error) {
	//fmt.Println("----Lookup----")
	lock.RLock()
	targets, ok := r[serviceName]
	lock.RUnlock()
	if !ok {
		return nil, ErrServiceNotFound
	}
	return targets, nil
}

func (r DefaultRegistry) Failure(host, path, endpoint string, err error) {
	// Would be used to remove an endpoint from the rotation, log the failure, etc.
	//log.Printf("Error accessing %s/%s (%s): %s", path, endpoint, err)
	log.Printf("Error accessing %s %s (%s): %s", host, path, endpoint, err)
}

func (r DefaultRegistry) Delete(serviceName string) {
	fmt.Println("*****Delete*****")
	lock.Lock()
	defer lock.Unlock()

	_, ok := r[serviceName]
	if !ok {
		return
	}

	delete(r, serviceName)
}
