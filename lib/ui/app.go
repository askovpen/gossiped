package ui

import(
  "github.com/jroimartin/gocui"
)

var (
  App  *gocui.Gui
  AreaList  *gocui.View
  Status  *gocui.View
  AreaPosition uint16
  ActiveWindow string
  parentWindow string
  curAreaId int
  curMsgNum uint32
  showKludges bool
)

func Quit(g *gocui.Gui, v *gocui.View) error {
  return gocui.ErrQuit
}

