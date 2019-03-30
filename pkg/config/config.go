package config

import (
	"errors"
	"github.com/askovpen/goated/pkg/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type configS struct {
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

//go:generate go run ../../util/autoversion/main.go
var (
	Config   configS
	PID      = "ATED+" + runtime.GOOS[0:3] + " " + Version
	LongPID  = "goAtEd-" + runtime.GOOS + "/" + runtime.GOARCH + " " + Version
	Template []string
)

// Read config
func Read() error {
	yamlFile, err := ioutil.ReadFile(os.Args[1])
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
	for _, l := range strings.Split(string(tpl[:]), "\n") {
		if len(l) > 0 && l[0] == ';' {
			continue
		}
		Template = append(Template, l)
	}
	if len(Config.Tearline) == 0 {
		Config.Tearline = LongPID
	}
	return nil
}
