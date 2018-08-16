package ui

import (
	"github.com/jroimartin/gocui"
	"github.com/askovpen/goated/lib/msgapi"
)

var (
	App						*gocui.Gui
	AreaList			*gocui.View
	AreaPosition	uint16
	ActiveWindow	string
	parentWindow	string
	curAreaId			int
	curMsgNum			uint32
	showKludges		bool
	StatusLine		string
	newMsg				*msgapi.Message
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
