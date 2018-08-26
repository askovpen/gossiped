package main

import (
	"fmt"
	"github.com/askovpen/goated/lib/areasconfig"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/ui"
	"github.com/askovpen/gocui"
	"log"
	"os"
	"time"
)

func main() {
	log.Printf("%s started", config.LongPID)
	if len(os.Args) == 1 {
		log.Printf("Usage: %s <config.yml>", os.Args[0])
		return
	}

	err := config.Read()
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
	ui.App, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer ui.App.Close()

	ui.App.InputEsc = true
	ui.App.SetManagerFunc(ui.Layout)
	ui.App.BgColor = gocui.ColorBlack
	ui.App.FgColor = gocui.ColorWhite
	ui.ActiveWindow = "AreaList"
	if err := ui.App.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		log.Panicln(err)
	}
	if err := ui.Keybindings(ui.App); err != nil {
		log.Panicln(err)
	}
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			ui.StatusTime = fmt.Sprintf("│ %s ", t.Format("15:04:05"))
			ui.App.Update(func(*gocui.Gui) error {
				return nil
			})
		}
	}()
	if err := ui.App.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
