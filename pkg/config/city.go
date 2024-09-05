package config

import (
	"os"

	"github.com/askovpen/gossiped/pkg/types"
	"gopkg.in/yaml.v3"
)

func readCity() {
	yamlFile, err := os.ReadFile(Config.CityPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &city)
	if err != nil {
		panic(err)
	}
}

// GetCity return city
func GetCity(sfa *types.FidoAddr) string {
	if val, ok := city[sfa.ShortString()]; ok {
		return val
	}
	return "unknown"
}
