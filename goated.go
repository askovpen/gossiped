package main

import (
//  "fmt"
//  "github.com/gdamore/tcell"
//  "github.com/rivo/tview"
  "log"
  "os"
  "github.com/askovpen/goated/lib/config"
  "github.com/askovpen/goated/lib/fidoconfig"
  "github.com/askovpen/goated/lib/ui"
  "github.com/jroimartin/gocui"
  "time"
//  "strconv"
)
var clock *time.Ticker
func main() {
  if len(os.Args)==1 {
    log.Printf("Usage: %s <config.yml>",os.Args[0])
    return
  }
  
  err:=config.Read()
  if err!=nil {
    log.Print(err)
    return
  }
  log.Printf("%s started",config.LongPID)
  f, _ := os.OpenFile(config.Config.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  defer f.Close()
  log.SetOutput(f)
  fidoconfig.Read()
  ui.App, err = gocui.NewGui(gocui.OutputNormal)
  if err != nil {
    log.Panicln(err)
  }
  defer ui.App.Close()
  ui.App.Highlight = true
  ui.App.SelFgColor = gocui.ColorBlue | gocui.AttrBold
  ui.App.SetManagerFunc(ui.Layout)
  //ui.App.SetCurrentView("AreaList")
  ui.ActiveWindow="AreaList"
  //ui.CreateAreaList()
  if err := ui.App.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
    log.Panicln(err)
  }
  if err := ui.Keybindings(ui.App); err != nil {                                                                                                                                                                    
    log.Panicln(err)                                                                                                                                                                                        
  }
  if err := ui.App.MainLoop(); err != nil && err != gocui.ErrQuit {
    log.Panicln(err)
  }
}
