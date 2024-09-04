package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func readCity() {
	yamlFile, err := ioutil.ReadFile("city.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &city)
	if err != nil {
		return
	}
}

// GetCity return city
func GetCity(sa string) string {
	if val, ok := city[sa]; ok {
		return val
	}
	return "unknown"
}
