package countryregistry

import (
	"errors"
	"fmt"
	//"strconv"
	"sync"
)

var lock sync.RWMutex

// Common errors.
var (
	ErrClusterNotFound = errors.New("Cluster not found")
)

// for a given service name / version pair.
type Registry interface {
	Lookup(country string) (string, error)
}

type DefaultCountryInfo map[string]string

func (c DefaultCountryInfo) Lookup(country string) (string, error) {
	fmt.Println("----Country Lookup----")
	lock.RLock()
	continent, ok := c[country]
	lock.RUnlock()
	if !ok {
		return "", ErrClusterNotFound
	}
	return continent, nil
}
