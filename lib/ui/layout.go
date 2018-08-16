package ui

import (
	"fmt"
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
	status, err := g.SetView("status", -1, maxY-2, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	status.Frame = false
	status.Wrap = false
	status.BgColor = gocui.ColorBlue
	status.FgColor = gocui.ColorWhite | gocui.AttrBold
	status.Clear()
	fmt.Fprintf(status, StatusLine)
	err = CreateAreaList()
	if err != nil {
		log.Print(err)
	}
	setCurrentViewOnTop(g, ActiveWindow)
	return nil
}
