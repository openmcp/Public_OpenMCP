package util

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)
type OmcpctlConf struct {
	OpenmcpApiServer string `yaml:"openmcpAPIServer"`
	OpenmcpDir string `yaml:"openmcpDir"`
	NfsServer  string `yaml:"nfsServer"`
}

func GetOmcpctlConf(filepath string) *OmcpctlConf {
	confStruct := &OmcpctlConf{}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &confStruct)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return confStruct
}
