package ui

import (
//  "fmt"
//  "github.com/askovpen/goated/lib/msgapi"
  "github.com/jroimartin/gocui"
//  "strconv"
  "log"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
  if _, err := g.SetCurrentView(name); err != nil {
    return nil, err
  }
  return g.SetViewOnTop(name)
}
/*
func getAreaNew(m msgapi.AreaPrimitive) string {
  if m.GetCount()-m.GetLast()>0 {
    return "\033[37;1m+\033[0m"
  } else {
    return " "
  }
}*/
func Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  Status, err:= g.SetView("status", -1, maxY-2, maxX, maxY);
  if err!=nil && err!=gocui.ErrUnknownView { 
    return err
  }
  Status.Frame=false
  Status.Wrap=false
  Status.BgColor=gocui.ColorBlue
  Status.FgColor=gocui.ColorWhite|gocui.AttrBold
  Status.Clear()
//  fmt.Fprintf(Status," Loading...")
  err=CreateAreaList()
  if err!=nil {
    log.Print(err)
  }
  setCurrentViewOnTop(g, ActiveWindow)
  return nil
  
}
