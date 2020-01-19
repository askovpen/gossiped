package ui

import (
	//"errors"
	"fmt"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/ui/editor"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	//"strings"
)

// ViewMsg widget
func (a *App) ViewMsg(areaID int, msgNum uint32) (string, tview.Primitive, bool, bool) {
	msg, err := msgapi.Areas[areaID].GetMsg(msgNum)
	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				a.Pages.SwitchToPage("AreaList")
			})
		return fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum), modal, true, true
	}
	if msg != nil {
		msgapi.Areas[areaID].SetLast(msgNum)
	}
	if msgapi.Areas[areaID].GetCount()-msgapi.Areas[areaID].GetLast() > 0 {
		a.al.SetCell(areaID+1, 0, tview.NewTableCell(strconv.FormatInt(int64(areaID), 10)+"[::b]+").SetAlign(tview.AlignRight))
	} else {
		a.al.SetCell(areaID+1, 0, tview.NewTableCell(strconv.FormatInt(int64(areaID), 10)+" ").SetAlign(tview.AlignRight))
	}
	a.al.SetCell(areaID+1, 2, tview.NewTableCell(strconv.FormatInt(int64(msgapi.Areas[areaID].GetCount()), 10)).SetAlign(tview.AlignRight))
	a.al.SetCell(areaID+1, 3, tview.NewTableCell(strconv.FormatInt(int64(msgapi.Areas[areaID].GetCount()-msgapi.Areas[areaID].GetLast()), 10)).SetAlign(tview.AlignRight))
	a.sb.SetStatus(fmt.Sprintf("Msg %d of %d (%d left)",
		msgNum,
		msgapi.Areas[areaID].GetCount(),
		msgapi.Areas[areaID].GetCount()-msgNum,
	))
	header := NewViewHeader(msg)
	header.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" " + msgapi.Areas[areaID].GetName() + " ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorYellow)
	var body *editor.View
	if msg != nil {
		body = editor.NewView(editor.NewBufferFromString(msg.ToView(a.showKludges)))
	} else {
		body = editor.NewView(editor.NewBufferFromString(""))
	}
	header.SetDoneFunc(func(s string) {
		num, _ := strconv.ParseUint(s, 10, 32)
		if uint32(num) >= msgapi.Areas[areaID].GetCount() {
			a.App.SetFocus(body)
		} else {
			if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), num)) {
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), num))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			} else {
				a.Pages.AddPage(a.ViewMsg(areaID, uint32(num)))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), num))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			}
		}
	})

	body.Readonly = true
	body.SetDoneFunc(func() {
		//		if key == tcell.KeyEscape {
		a.Pages.SwitchToPage("AreaList")
		a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
		//		}
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			if msgNum == msgapi.Areas[areaID].GetCount() {
				a.Pages.SwitchToPage("AreaList")
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum+1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum+1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
				} else {
					a.Pages.AddPage(a.ViewMsg(areaID, msgNum+1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum+1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
				}
			}
		} else if event.Key() == tcell.KeyLeft {
			if msgNum <= 1 {
				a.Pages.SwitchToPage("AreaList")
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum-1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum-1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
				} else {
					a.Pages.AddPage(a.ViewMsg(areaID, msgNum-1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum-1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
				}
			}
		} else if event.Key() == tcell.KeyInsert || event.Key() == tcell.KeyCtrlI {
			a.Pages.AddPage(a.InsertMsg(areaID, 0))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaID].GetName()))
		} else if msg == nil {
			return event
		} else if event.Key() == tcell.KeyCtrlK || (event.Rune() == 'k' && event.Modifiers()&tcell.ModAlt > 0) {
			a.showKludges = !a.showKludges
			//body.SetText(msg.ToView(a.showKludges))
			body.OpenBuffer(editor.NewBufferFromString(msg.ToView(a.showKludges)))
		} else if event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyF3 || (event.Rune() == 'q') {
			a.Pages.AddPage(a.InsertMsg(areaID, newMsgTypeAnswer))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaID].GetName()))
		} else if event.Key() == tcell.KeyCtrlN || (event.Rune() == 'n' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(areaID, newMsgTypeAnswerNewArea))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyCtrlF || (event.Rune() == 'f' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(areaID, newMsgTypeForward))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyDelete {
			a.Pages.AddPage(a.showDelMsg(areaID, msgNum))
			a.Pages.ShowPage("DelMsgModal")
		} else if event.Key() == tcell.KeyCtrlL || event.Rune() == 'l' {
			a.Pages.AddPage(a.showMessageList(areaID))
			a.Pages.ShowPage("MessageListModal")
		} else if event.Key() == tcell.KeyCtrlG || event.Rune() == 'g' {
			a.App.SetFocus(header)
			//a.Pages.AddPage(a.showMessageList(areaID))
			//a.Pages.ShowPage("MessageListModal")
		} else if event.Rune() == '<' {
			if msgNum != 1 {
				a.Pages.AddPage(a.ViewMsg(areaID, 1))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), 1))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			}
		} else if event.Rune() == '>' {
			if msgNum != msgapi.Areas[areaID].GetCount() {
				a.Pages.AddPage(a.ViewMsg(areaID, msgapi.Areas[areaID].GetCount()))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgapi.Areas[areaID].GetCount()))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			}
		}

		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(body, 0, 1, true)
	return fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum), layout, true, true
}

func (a *App) showMessageList(areaID int) (string, tview.Primitive, bool, bool) {
	modal := NewModalMessageList(areaID).
		SetDoneFunc(func(msgNum uint32) {
			a.Pages.HidePage("MessageListModal")
			a.Pages.RemovePage("MessageListModal")
			a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgapi.Areas[areaID].GetLast()))
			a.Pages.AddPage(a.ViewMsg(areaID, msgNum))
			a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			a.App.SetFocus(a.Pages)
		})
	return "MessageListModal", modal, true, true
}
func (a *App) showAreaList(areaID int, newMsgType int) (string, tview.Primitive, bool, bool) {
	modal := NewModalAreaList().
		SetDoneFunc(func(buttonIndex int) {
			a.im.postArea = buttonIndex - 1
			a.Pages.HidePage("AreaListModal")
			a.Pages.RemovePage("AreaListModal")
			a.Pages.AddPage(a.InsertMsg(areaID, newMsgType))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaID].GetName()))
			a.App.SetFocus(a.Pages)
		})
	if newMsgType == newMsgTypeAnswerNewArea {
		modal.SetText("Answer In Area:")
	}
	if newMsgType == newMsgTypeForward {
		modal.SetText("Forward To Area:")
	}
	return "AreaListModal", modal, true, true
}
func (a *App) showDelMsg(areaID int, msgNum uint32) (string, tview.Primitive, bool, bool) {
	modal := NewModalMenu().
		SetY(6).
		SetText("Delete?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int) {
			a.Pages.HidePage("DelMsgModal")
			a.Pages.RemovePage("DelMsgModal")
			if buttonIndex == 0 {
				msgapi.Areas[areaID].DelMsg(msgNum)
				a.Pages.AddPage(a.ViewMsg(areaID, msgNum-1))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum-1))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaID].GetName(), msgNum))
			} else {
			}
			a.App.SetFocus(a.Pages)
		})
	return "DelMsgModal", modal, true, true
}
