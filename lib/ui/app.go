package ui

import (
	"github.com/askovpen/goated/lib/msgapi"
	"github.com/askovpen/gocui"
)

var (
	// App gui
	App *gocui.Gui
	// AreaList variable
	AreaList *gocui.View
	// AreaPosition variable
	AreaPosition uint16
	// ActiveWindow name
	ActiveWindow string
	parentWindow string
	curAreaID    int
	curMsgNum    uint32
	showKludges  bool
	// StatusLine variable
	StatusLine string
	// StatusTime variable
	StatusTime string
	newMsg     *msgapi.Message
	newMsgType string
)

// Quit application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
