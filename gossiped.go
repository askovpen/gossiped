package main

import (
	"github.com/askovpen/gossiped/pkg/areasconfig"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/ui"
	"github.com/askovpen/gossiped/pkg/utils"
	"log"
	"os"
	"path/filepath"
)

var (
	version = "2.0"
	commit  = "dev"
	date    = "unknown"
)

func tryFindConfig() string {
	for _, fn := range []string{
		filepath.Join(os.Getenv("HOME"), "gossiped.yml"),
		filepath.Join(os.Getenv("HOME"), ".config", "gossiped.yml"),
		"/usr/local/etc/ftn/gossiped.yml",
		"/etc/ftn/gossiped.yml",
		"gossiped.yml",
	} {
		if utils.FileExists(fn) {
			return fn
		}
	}
	return ""
}

func main() {
	config.Version = version + "-" + commit
	config.InitVars()
	log.Printf("%s started", config.LongPID)
	var fn string
	if len(os.Args) == 1 {
		fn = tryFindConfig()
		if fn == "" {
			log.Printf("Usage: %s <config.yml>", os.Args[0])
			return
		}
	} else {
		if utils.FileExists(os.Args[1]) {
			fn = os.Args[1]
		} else {
			log.Printf("Usage: %s <config.yml>", os.Args[0])
			return
		}
	}

	err := config.Read(fn)
	if err != nil {
		log.Print(err)
		return
	}
	f, _ := os.OpenFile(config.Config.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	err = areasconfig.Read()
	if err != nil {
		log.Print(err)
		return
	}
	// ui.App, err = gocui.NewGui(gocui.OutputNormal)
	app := ui.NewApp()
	log.Print("start")
	if err = app.Run(); err != nil {
		log.Print("started ui")
		log.Print(err)
		return
	}
}
