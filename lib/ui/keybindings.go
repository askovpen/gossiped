package ui

import (
//  "fmt"
//  "github.com/askovpen/goated/lib/msgapi"
  "github.com/jroimartin/gocui"
//  "strconv"
//  "log"
)

func Keybindings(g *gocui.Gui) error {
  if err := App.SetKeybinding("AreaList", gocui.KeyArrowDown, gocui.ModNone, areaNext); err != nil {
    return err
  }
  if err := App.SetKeybinding("AreaList", gocui.KeyArrowUp, gocui.ModNone, areaPrev); err != nil {
    return err
  }
  if err := App.SetKeybinding("AreaList", gocui.KeyEnter, gocui.ModNone, viewArea); err != nil {
    return err
  }
  return nil
}
