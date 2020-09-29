package clusterregistry

import (
	"errors"
	"openmcp/openmcp/omcplog"
	"strconv"
	"sync"
)

// DefaultClusterInfo is a basic registry using the following format:
//{
//	"cluster1" : {
//		"Latitude" : "10.32" ,
//		"Longitude" : " 20.44",
//		"IngressIP" : "10.0.3.203",
//		"Country" : "US",
//		"Continent" : "Asia",
//		"ResourceScore" : "80",
//		"HopScore" : "60",
//	},
//}

var lock sync.RWMutex

// Common errors.
var (
	ErrClusterNotFound = errors.New("Cluster not found")
)

// for a given service name / version pair.
type Registry interface {
	Add(ClusterName, Latitude, Longitude, IngressIP, Country, Continent, ResourceScore, HopScore string) // Add an endpoint to our registry
	Longitude(ClusterName string) (float64, error)
	Latitude(ClusterName string) (float64, error)
	IngressIP(ClusterName string) (string, error)
	Country(ClusterName string) (string, error)
	Continent(ClusterName string) (string, error)
	ResourceScore(ClusterName string) (float64, error)
	HopScore(ClusterName string) (float64, error)
}

type DefaultClusterInfo map[string]map[string]string

func (c DefaultClusterInfo) Lookup(cluster string) (bool, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Lookup")
	lock.RLock()
	_, ok := c[cluster]
	lock.RUnlock()
	if !ok {
		return false, ErrClusterNotFound
	}
	return true, nil
}

func (c DefaultClusterInfo) IngressIP(ClusterName string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster IngressIP")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return "", ErrClusterNotFound
	}
	IngressIP := cluster["IngressIP"]
	return IngressIP, nil
}

func (c DefaultClusterInfo) Longitude(ClusterName string) (float64, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Longitude")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return 0, ErrClusterNotFound
	}
	longitude, _ := strconv.ParseFloat(cluster["Longitude"], 64)
	return longitude, nil
}

func (c DefaultClusterInfo) Latitude(ClusterName string) (float64, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Latitude")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return 0, ErrClusterNotFound
	}
	Latitude, _ := strconv.ParseFloat(cluster["Latitude"], 64)
	return Latitude, nil
}

func (c DefaultClusterInfo) Country(ClusterName string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Country")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return "", ErrClusterNotFound
	}
	country := cluster["Country"]
	return country, nil
}

func (c DefaultClusterInfo) Continent(ClusterName string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Continent")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return "", ErrClusterNotFound
	}
	continent := cluster["Continent"]
	return continent, nil
}

func (c DefaultClusterInfo) ResourceScore(ClusterName string) (float64, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster ResourceScore")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return 0, ErrClusterNotFound
	}
	resourceScore, _ := strconv.ParseFloat(cluster["ResourceScore"], 64)
	return resourceScore, nil
}

func (c DefaultClusterInfo) HopScore(ClusterName string) (float64, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster HopScore")
	lock.RLock()
	cluster, ok := c[ClusterName]
	lock.RUnlock()
	if !ok {
		return 0, ErrClusterNotFound
	}
	resourceScore, _ := strconv.ParseFloat(cluster["HopScore"], 64)
	return resourceScore, nil
}

func (c DefaultClusterInfo) Add(ClusterName, Latitude, Longitude, IngressIP, Country, Continent, ResourceScore, HopScore string) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(ClusterRegistry)] Function Cluster Add")
	lock.Lock()
	defer lock.Unlock()

	cluster, ok := c[ClusterName]
	if !ok {
		cluster = map[string]string{}
		c[ClusterName] = cluster
	}
	cluster["Latitude"] = Latitude
	cluster["Longitude"] = Longitude
	cluster["IngressIP"] = IngressIP
	cluster["Country"] = Country
	cluster["Continent"] = Continent
	cluster["ResourceScore"] = ResourceScore
	cluster["HopScore"] = HopScore
}
