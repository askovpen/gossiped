package ui

import (
	"github.com/askovpen/goated/pkg/msgapi"
	"github.com/askovpen/gocui"
)

var (
	// App gui
	App *gocui.Gui
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
	StatusTime   string
	newMsg       *msgapi.Message
	newMsgType   int
	newMsgAreaID int
)

const (
	newMsgTypeAnswer        = 1
	newMsgTypeAnswerNewArea = 2
	newMsgTypeForward       = 4
)

// Quit application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
