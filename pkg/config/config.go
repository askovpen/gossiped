package config

import (
	"errors"
	"github.com/askovpen/gossiped/pkg/types"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"runtime"
	"strings"
)

type (
	ColorMap map[string]string
	configS  struct {
		Username string
		AreaFile struct {
			Path string
			Type string
		}
		Areas []struct {
			Name     string
			Path     string
			Type     string
			BaseType string
			Chrs     string
		}
		Colors   map[string]ColorMap
		Log      string
		Address  *types.FidoAddr
		Origin   string
		Tearline string
		Template string
		Chrs     struct {
			Default string
			IBMPC   string
		}
	}
)

// vars
var (
	Version  string
	PID      string
	LongPID  string
	Config   configS
	Template []string
	city     map[string]string
)

// InitVars define version variables
func InitVars() {
	PID = "gossipEd+" + runtime.GOOS[0:3] + " " + Version
	LongPID = "gossipEd-" + runtime.GOOS + "/" + runtime.GOARCH + " " + Version
}

// Read config
func Read(fn string) error {
	yamlFile, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return err
	}
	if Config.Address == nil {
		return errors.New("Config.Address not defined")
	}
	if Config.Chrs.Default == "" {
		return errors.New("Config.Chrs.Default not defined")
	}
	tpl, err := ioutil.ReadFile(Config.Template)
	if err != nil {
		return err
	}
	for _, l := range strings.Split(string(tpl), "\n") {
		if len(l) > 0 && l[0] == ';' {
			continue
		}
		Template = append(Template, l)
	}
	if len(Config.Tearline) == 0 {
		Config.Tearline = LongPID
	}
	readCity()
	return nil
}
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

func (defaultColors ColorMap) ProduceFromConfig(colorArea string) ColorMap {
	var out = make(ColorMap)
	var validKeys = make(map[string]bool)
	for k, v := range defaultColors {
		validKeys[k] = true
		out[k] = v
	}
	if Config.Colors[colorArea] != nil {
		for element, colorValue := range Config.Colors[colorArea] {
			if validKeys[element] {
				out[element] = colorValue
			}
		}
	}
	return out
}
