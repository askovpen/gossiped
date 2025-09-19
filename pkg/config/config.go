package config

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/askovpen/gossiped/pkg/nodelist"
	"github.com/askovpen/gossiped/pkg/types"
	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

type (
	ColorMap    map[string]string
	SortTypeMap map[string]string
	configS     struct {
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
		Colorscheme string
		Log         string
		Address     *types.FidoAddr
		Origin      string
		Tearline    string
		Template    string
		Chrs        struct {
			Default string
			IBMPC   string
		}
		Statusbar struct {
			Clock bool
		}
		Sorting      SortTypeMap
		Colors       map[string]ColorMap
		CityPath     string
		NodelistPath string
	}
)

// vars
var (
	Version      string
	PID          string
	LongPID      string
	Config       configS
	Template     []string
	city         map[string]string
	StyleDefault tcell.Style
)

// InitVars define version variables
func InitVars() {
	PID = "gossipEd+" + runtime.GOOS[0:3] + " " + Version
	LongPID = "gossipEd-" + runtime.GOOS + "/" + runtime.GOARCH + " " + Version
}
func tryPath(rootPath string, filePath string) string {
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	}
	if _, err := os.Stat(path.Join(rootPath, filePath)); err == nil {
		return path.Join(rootPath, filePath)
	}
	return ""
}

// Read config
func Read(fn string) error {
	yamlFile, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	rootPath := filepath.Dir(fn)

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
	Config.Template = tryPath(rootPath, Config.Template)
	tpl, err := os.ReadFile(Config.Template)
	if err != nil {
		return err
	}
	readTemplate(tpl)
	if len(Config.Tearline) == 0 {
		Config.Tearline = LongPID
	}
	errColors := readColors(rootPath)
	if errColors != nil {
		return errColors
	}
	if Config.CityPath == "" {
		return errors.New("Config.CityPath not defined")
	}
	Config.CityPath = tryPath(rootPath, Config.CityPath)
	err = readCity()
	Config.NodelistPath = tryPath(rootPath, Config.NodelistPath)
	nodelist.Read(Config.NodelistPath)
	if err != nil {
		return err
	}
	return nil
}

func readTemplate(tpl []byte) {
	for _, l := range strings.Split(string(tpl), "\n") {
		if len(l) > 0 && l[0] == ';' {
			continue
		}
		Template = append(Template, l)
	}
}
