package ui

import (
	"fmt"
	"github.com/askovpen/gocui"
	"log"
)

func errorMsg(msg string, parent string) error {
	log.Printf("error msg: %s (parent: %s)", msg, parent)
	maxX, maxY := App.Size()
	parentWindow = parent
	var v *gocui.View
	o := 0
	if len(msg)/2 != 0 {
		o = 1
	}
	if len(msg) < maxX-10 {
		v, _ = App.SetView("ErrorMsg", maxX/2-(len(msg)/2)-1, maxY/2-1, maxX/2+(len(msg))/2+1+o, maxY/2+1)
	} else {
		v, _ = App.SetView("ErrorMsg", maxX/2-30, maxY/2-2, maxX/2+30, maxY/2+2)
	}
	v.Wrap = true
	v.BgColor = gocui.ColorRed
	v.FgColor = gocui.ColorWhite | gocui.AttrBold
	fmt.Fprintf(v, " %s", msg)
	ActiveWindow = "ErrorMsg"
	App.SetCurrentView("ErrorMsg")
	return nil
}
func exitError(g *gocui.Gui, v *gocui.View) error {
	ActiveWindow = parentWindow
	g.DeleteView("ErrorMsg")
	ActiveWindow = parentWindow
	log.Printf("Aw: " + ActiveWindow)
	return nil
}
