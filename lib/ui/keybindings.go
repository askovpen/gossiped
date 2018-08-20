package ui

import (
	//  "fmt"
	//  "github.com/askovpen/goated/lib/msgapi"
	"github.com/askovpen/gocui"
	//  "strconv"
	//  "log"
)

func Keybindings(g *gocui.Gui) error {
	g.SetKeybinding("AreaList", gocui.KeyArrowDown, gocui.ModNone, areaNext)
	g.SetKeybinding("AreaList", gocui.KeyEsc, gocui.ModNone, quitAreaList)
	g.SetKeybinding("AreaList", gocui.KeyArrowUp, gocui.ModNone, areaPrev)
	g.SetKeybinding("AreaList", gocui.KeyEnter, gocui.ModNone, viewArea)
	g.SetKeybinding("AreaList", gocui.KeyArrowRight, gocui.ModNone, viewArea)

	g.SetKeybinding("QuitMsg", gocui.KeyArrowDown, gocui.ModNone, quitUp)
	g.SetKeybinding("QuitMsg", gocui.KeyArrowUp, gocui.ModNone, quitUp)
	g.SetKeybinding("QuitMsg", gocui.KeyEnter, gocui.ModNone, quitEnter)

	g.SetKeybinding("MsgBody", gocui.KeyArrowDown, gocui.ModNone, scrollDown)
	g.SetKeybinding("MsgBody", gocui.KeyArrowUp, gocui.ModNone, scrollUp)
	g.SetKeybinding("MsgBody", gocui.KeyArrowLeft, gocui.ModNone, prevMsg)
	g.SetKeybinding("MsgBody", gocui.KeyArrowRight, gocui.ModNone, nextMsg)
	g.SetKeybinding("MsgBody", '<', gocui.ModNone, firstMsg)
	g.SetKeybinding("MsgBody", '>', gocui.ModNone, lastMsg)
	g.SetKeybinding("MsgBody", gocui.KeyEsc, gocui.ModNone, quitMsgView)
	g.SetKeybinding("MsgBody", 'k', gocui.ModAlt, toggleKludges)
	g.SetKeybinding("MsgBody", gocui.KeyCtrlK, gocui.ModNone, toggleKludges)
	g.SetKeybinding("MsgBody", gocui.KeyCtrlG, gocui.ModNone, editMsgNum)
	g.SetKeybinding("MsgBody", gocui.KeyInsert, gocui.ModNone, editMsg)
	g.SetKeybinding("MsgBody", gocui.KeyCtrlQ, gocui.ModNone, answerMsg)
	g.SetKeybinding("MsgBody", gocui.KeyF3, gocui.ModNone, answerMsg)
	g.SetKeybinding("MsgBody", 'q', gocui.ModAlt, answerMsg)

	g.SetKeybinding("editToName", gocui.KeyEnter, gocui.ModNone, editToNameNext)
	g.SetKeybinding("editToName", gocui.KeyTab, gocui.ModNone, editToNameNext)
	g.SetKeybinding("editToAddr", gocui.KeyEnter, gocui.ModNone, editToAddrNext)
	g.SetKeybinding("editToAddr", gocui.KeyTab, gocui.ModNone, editToAddrNext)
	g.SetKeybinding("editSubj", gocui.KeyEnter, gocui.ModNone, editToSubjBody)
	g.SetKeybinding("editSubj", gocui.KeyTab, gocui.ModNone, editToSubjNext)
	g.SetKeybinding("editFromName", gocui.KeyEnter, gocui.ModNone, editFromNameNext)
	g.SetKeybinding("editFromName", gocui.KeyTab, gocui.ModNone, editFromNameNext)
	g.SetKeybinding("editFromAddr", gocui.KeyEnter, gocui.ModNone, editFromAddrNext)
	g.SetKeybinding("editFromAddr", gocui.KeyTab, gocui.ModNone, editFromAddrNext)
	g.SetKeybinding("editMsgBody", gocui.KeyCtrlS, gocui.ModNone, editMsgBodyMenu)
	g.SetKeybinding("editMsgBody", gocui.KeyF2, gocui.ModNone, editMsgBodyMenu)
	g.SetKeybinding("editMsgBody", gocui.KeyEsc, gocui.ModNone, editMsgBodyMenu)
	g.SetKeybinding("editMenuMsg", gocui.KeyEnter, gocui.ModNone, saveMessage)
	g.SetKeybinding("editMenuMsg", gocui.KeyArrowUp, gocui.ModNone, editMsgBodyMenuUp)
	g.SetKeybinding("editMenuMsg", gocui.KeyArrowDown, gocui.ModNone, editMsgBodyMenuDown)

	g.SetKeybinding("editNumber", gocui.KeyEnter, gocui.ModNone, editMsgNumEnter)
	g.SetKeybinding("ErrorMsg", gocui.KeyEnter, gocui.ModNone, exitError)

	return nil
}
