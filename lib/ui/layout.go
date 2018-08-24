package ui

import (
	"fmt"
	"github.com/askovpen/gocui"
	"log"
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	status, err := g.SetView("status", -1, maxY-2, maxX-11, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	statusTime, err := g.SetView("statusTime", maxX-12, maxY-2, maxX, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	status.Frame = false
	status.Wrap = false
	status.BgColor = gocui.ColorBlue
	status.FgColor = gocui.ColorWhite | gocui.AttrBold
	status.Clear()
	statusTime.Frame = false
	statusTime.Wrap = false
	statusTime.BgColor = gocui.ColorBlue
	statusTime.FgColor = gocui.ColorWhite | gocui.AttrBold
	statusTime.Clear()
	fmt.Fprintf(status, StatusLine)
	fmt.Fprintf(statusTime, StatusTime)
	err = CreateAreaList()
	if err != nil {
		log.Print(err)
	}
	setCurrentViewOnTop(g, ActiveWindow)
	return nil
}
