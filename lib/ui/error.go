package ui

import (
  "fmt"
  "github.com/jroimartin/gocui"
  "log"
)

func errorMsg(msg string, parent string) error {
  maxX, maxY := App.Size()
  parentWindow=parent
  v, _ := App.SetView("ErrorMsg", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
  v.Wrap = true
  v.BgColor=gocui.ColorRed
  v.FgColor=gocui.ColorWhite | gocui.AttrBold
  fmt.Fprintf(v, msg)
  ActiveWindow="ErrorMsg"
  App.SetCurrentView("ErrorMsg")
  return nil
}
func exitError(g *gocui.Gui, v *gocui.View) error {
  ActiveWindow=parentWindow
  g.DeleteView("ErrorMsg")
  ActiveWindow=parentWindow
  log.Printf("Aw: "+ActiveWindow)
  return nil
}
