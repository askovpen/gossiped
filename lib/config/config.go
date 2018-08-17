package config

import (
	"errors"
	"github.com/askovpen/goated/lib/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type config_s struct {
	Username   string
	FidoConfig string
	Log        string
	Address    *types.FidoAddr
	Origin     string
	Template   string
}

var (
	Config   config_s
	Version  string
	PID      string
	LongPID  string
	Template []string
)

func Read() error {
	Version = "0.0.1"
	PID = "ATED+" + runtime.GOOS[0:3] + " " + Version
	LongPID = "goAtEd-" + runtime.GOOS + "/" + runtime.GOARCH + " " + Version
	yamlFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		return err
	}
	if Config.Address == nil {
		return errors.New("Address not defined")
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
	return nil
}
