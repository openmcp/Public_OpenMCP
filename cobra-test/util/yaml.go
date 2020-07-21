package util

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
)

func GetYaml(filepath string) interface{} {
	var yamlStruct interface{}

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &yamlStruct)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return yamlStruct
}
