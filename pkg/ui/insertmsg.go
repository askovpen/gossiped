package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/types"
	"github.com/askovpen/gossiped/pkg/ui/editor"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	//"log"
)

const (
	newMsgTypeAnswer        = 1
	newMsgTypeAnswerNewArea = 2
	newMsgTypeForward       = 4
)

// IM struct
type IM struct {
	eb         *editor.View
	eh         *EditHeader
	newMsg     *msgapi.Message
	curArea    int
	postArea   int
	newMsgType int
	buffer     *editor.Buffer
}

// InsertMsgMenu modal menu
func (a *App) InsertMsgMenu() (string, tview.Primitive, bool, bool) {
	modal := NewModalMenu().
		SetY(6).
		SetText("Save?").
		AddButtons([]string{"Yes", "No, Drop", "Continue Writing", "Edit Header"}).
		SetDoneFunc(func(buttonIndex int) {
			switch b := buttonIndex; b {
			case 0:
				//a.im.newMsg.Body = a.im.eb.GetText(false)
				a.im.newMsg.Body = a.im.buffer.String()
				msgapi.Areas[a.im.postArea].SaveMsg(a.im.newMsg.MakeBody())
				a.Pages.HidePage("InsertMsgMenu")
				a.Pages.RemovePage("InsertMsgMenu")
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[a.im.curArea].GetName(), msgapi.Areas[a.im.curArea].GetLast()))
				a.Pages.RemovePage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[a.im.curArea].GetName()))
				a.App.SetFocus(a.Pages)
			case 1:
				a.Pages.HidePage("InsertMsgMenu")
				a.Pages.RemovePage("InsertMsgMenu")
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[a.im.curArea].GetName(), msgapi.Areas[a.im.curArea].GetLast()))
				a.Pages.RemovePage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[a.im.curArea].GetName()))
				a.App.SetFocus(a.Pages)
			case 2:
				a.Pages.HidePage("InsertMsgMenu")
				a.App.SetFocus(a.im.eb)
			case 3:
				a.Pages.HidePage("InsertMsgMenu")
				a.App.SetFocus(a.im.eh)
			}
		})
	return "InsertMsgMenu", modal, false, false
}

// InsertMsg widget
func (a *App) InsertMsg(areaID int, msgType int) (string, tview.Primitive, bool, bool) {
	var omsg *msgapi.Message
	a.im.curArea = areaID
	a.im.newMsgType = msgType
	if a.im.newMsgType == 0 || a.im.newMsgType == newMsgTypeAnswer {
		a.im.postArea = areaID
	}
	a.im.newMsg = &msgapi.Message{From: config.Config.Username, FromAddr: config.Config.Address, AreaID: a.im.postArea}
	a.im.newMsg.Kludges = make(map[string]string)
	a.im.newMsg.Kludges["PID:"] = config.PID
	a.im.newMsg.Kludges["CHRS:"] = config.Config.Chrs.Default
	if msgapi.Areas[a.im.postArea].GetChrs() != "" {
		a.im.newMsg.Kludges["CHRS:"] = msgapi.Areas[a.im.postArea].GetChrs()
	}
	if msgapi.Areas[a.im.postArea].GetType() != msgapi.EchoAreaTypeNetmail && (a.im.newMsgType == 0 || a.im.newMsgType == newMsgTypeForward) {
		a.im.newMsg.To = "All"
	}
	if (a.im.newMsgType&newMsgTypeAnswer) != 0 || (a.im.newMsgType&newMsgTypeAnswerNewArea) != 0 {
		omsg, _ = msgapi.Areas[areaID].GetMsg(msgapi.Areas[a.im.curArea].GetLast())
		a.im.newMsg.To = omsg.From
		a.im.newMsg.ToAddr = omsg.FromAddr
		a.im.newMsg.Kludges["REPLY:"] = omsg.Kludges["MSGID:"]
		a.im.newMsg.Subject = omsg.Subject
	} else if (a.im.newMsgType & newMsgTypeForward) != 0 {
		omsg, _ = msgapi.Areas[areaID].GetMsg(msgapi.Areas[a.im.curArea].GetLast())
		a.im.newMsg.Subject = omsg.Subject
	}
	a.im.eh = NewEditHeader(a.im.newMsg)
	a.im.eh.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" " + msgapi.Areas[a.im.postArea].GetName() + " ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorYellow)
	a.im.eb = editor.NewView(editor.NewBufferFromString(""))
	//	a.im.eb = NewEditBody().
	a.im.eb.SetDoneFunc(func() {
		a.Pages.ShowPage("InsertMsgMenu")
		//			//log.Printf("%q",a.App.GetFocus())
	})
	a.im.eh.SetDoneFunc(func(r [5][]rune) {
		a.im.newMsg.From = string(r[0])
		a.im.newMsg.FromAddr = types.AddrFromString(string(r[1]))
		a.im.newMsg.To = string(r[2])
		a.im.newMsg.ToAddr = types.AddrFromString(string(r[3]))
		a.im.newMsg.Subject = string(r[4])
		/*
			if len(a.im.eb.GetText(false)) == 0 {
		*/
		var mv string
		//var p int
		if a.im.newMsgType == 0 {
			mv = a.im.newMsg.ToEditNewView()
		} else if a.im.newMsgType == newMsgTypeAnswer || a.im.newMsgType == newMsgTypeAnswerNewArea {
			mv = a.im.newMsg.ToEditAnswerView(omsg)
		} else if a.im.newMsgType == newMsgTypeForward {
			mv = a.im.newMsg.ToEditForwardView(omsg)
		}
		a.im.buffer = editor.NewBufferFromString(mv)
		//p = p
		a.im.eb.OpenBuffer(a.im.buffer)
		/*
			a.im.eb.SetText(mv, p)
		}*/
		a.App.SetFocus(a.im.eb)
	})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.im.eh, 6, 1, true).
		AddItem(a.im.eb, 0, 1, false)
	return fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaID].GetName()), layout, true, true
}
