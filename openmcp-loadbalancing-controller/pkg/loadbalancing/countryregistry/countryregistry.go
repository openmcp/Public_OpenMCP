package countryregistry

import (
	"errors"
	"openmcp/openmcp/omcplog"
	"sync"
)

var lock sync.RWMutex

// Common errors.
var (
	ErrClusterNotFound = errors.New("Cluster not found")
)

type Registry interface {
	Lookup(country string) (string, error)
}

type DefaultCountryInfo map[string]string

func (c DefaultCountryInfo) Lookup(country string) (string, error) {
	omcplog.V(4).Info("[OpenMCP Loadbalancing Controller(CountryRegistry)] Function Lookup")
	lock.RLock()
	continent, ok := c[country]
	lock.RUnlock()
	if !ok {
		return "", ErrClusterNotFound
	}
	return continent, nil
}
