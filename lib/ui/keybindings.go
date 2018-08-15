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
  if err := App.SetKeybinding("AreaList", gocui.KeyEsc, gocui.ModNone, quitAreaList); err != nil {
    return err
  }
  if err := App.SetKeybinding("AreaList", gocui.KeyArrowUp, gocui.ModNone, areaPrev); err != nil {
    return err
  }
  if err := App.SetKeybinding("AreaList", gocui.KeyEnter, gocui.ModNone, viewArea); err != nil {
    return err
  }
  if err := App.SetKeybinding("AreaList", gocui.KeyArrowRight, gocui.ModNone, viewArea); err != nil {
    return err
  }
  if err := App.SetKeybinding("QuitMsg", gocui.KeyArrowDown, gocui.ModNone, quitUp); err != nil {
    return err
  }
  if err := App.SetKeybinding("QuitMsg", gocui.KeyArrowUp, gocui.ModNone, quitUp); err != nil {
    return err
  }
  if err := App.SetKeybinding("QuitMsg", gocui.KeyEnter, gocui.ModNone, quitEnter); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyArrowLeft, gocui.ModNone, prevMsg); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyArrowRight, gocui.ModNone, nextMsg); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", '<', gocui.ModNone, firstMsg); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", '>', gocui.ModNone, lastMsg); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyEsc, gocui.ModNone, quitMsgView); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", 'k', gocui.ModAlt, toggleKludges); err != nil {
    return err
  }
  if err := App.SetKeybinding("MsgBody", gocui.KeyCtrlK, gocui.ModNone, toggleKludges); err != nil {
    return err
  }
  if err := App.SetKeybinding("ErrorMsg", gocui.KeyEnter, gocui.ModNone, exitError); err != nil {
    return err
  }

  return nil
}
