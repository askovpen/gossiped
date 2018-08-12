package main

import (
//  "fmt"
  "github.com/gdamore/tcell"
  "github.com/rivo/tview"
  "log"
  "os"
  "github.com/askovpen/goated/lib/config"
  "github.com/askovpen/goated/lib/fidoconfig"
  "github.com/askovpen/goated/lib/ui"
  "time"
//  "strconv"
)

var (
  el []EchoList
)
type EchoList struct {
  AreaNum uint16
  Name string
  Count uint32
  New   uint32
}

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
  f, _ := os.OpenFile(config.Config.Log, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
  defer f.Close()
  log.SetOutput(f)
  fidoconfig.Read()
  //return
  ui.App = tview.NewApplication()

  ui.Pages = tview.NewPages()
  ui.Pages.AddPage(ui.AreaList())

  ui.Status = tview.NewTextView().SetWrap(false)
  ui.Status.SetBackgroundColor(tcell.ColorBlue)
  ui.Status.SetTextColor(tcell.ColorWhite)
  ui.Status.SetDynamicColors(true)
  ui.Status.SetChangedFunc(func() {
    ui.App.Draw()
  })

  ui.StatusTime = tview.NewTextView().SetWrap(false)
  ui.StatusTime.SetBackgroundColor(tcell.ColorBlue)
  ui.StatusTime.SetTextColor(tcell.ColorWhite)
  ui.StatusTime.SetDynamicColors(true)
  ui.StatusTime.SetChangedFunc(func() {
    ui.App.Draw()
  })
  clock := time.NewTicker(1 * time.Second)
  go func() {
    for t := range clock.C {
      //fmt.Fprintf(ui.StatusTime,"%s",t.Format("15:04:05"))
      ui.StatusTime.SetText("[::b]"+t.Format("15:04:05"))
    }
  }()
//  return
  layout := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(ui.Pages, 0, 1, true).
    AddItem(tview.NewFlex().
      AddItem(ui.Status, 0, 1, false).
      AddItem(ui.StatusTime, 10, 1, false),1,1,false)
  if err := ui.App.SetRoot(layout, true).Run(); err != nil {
    panic(err)
  }
}

