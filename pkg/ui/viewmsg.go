package ui

import (
	//"errors"
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/ui/editor"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
	//"strings"
)

func (a *App) SwitchToAreaListPage() {
	a.RefreshAreaList()
	a.Pages.SwitchToPage("AreaList")
}

// ViewMsg widget
func (a *App) ViewMsg(area *msgapi.AreaPrimitive, msgNum uint32) (string, tview.Primitive, bool, bool) {
	msg, err := (*area).GetMsg(msgNum)
	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				a.SwitchToAreaListPage()
			})
		return fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum), modal, true, true
	}
	if msg != nil {
		if msgNum == 0 {
			msgNum = 1
		}
		(*area).SetLast(msgNum)
	}
	a.sb.SetStatus(fmt.Sprintf("%s: message %d of %d (%d left)",
		(*area).GetName(),
		msgNum,
		(*area).GetCount(),
		(*area).GetCount()-msgNum,
	))
	styleBorder := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementBorder)
	fgTitle, bgTitle, titleAttrs := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementTitle).Decompose()
	header := NewViewHeader(msg)
	header.SetBorder(true).
		SetBorderStyle(styleBorder).
		SetTitle(fmt.Sprintf("[%s:%s:%s] %s ",
			fgTitle.String(),
			bgTitle.String(),
			config.MaskToStringStyle(titleAttrs),
			(*area).GetName(),
		)).
		SetTitleAlign(tview.AlignLeft)
	var body *editor.View
	if msg != nil {
		body = editor.NewView(editor.NewBufferFromString(msg.ToView(a.showKludges)))
	} else {
		body = editor.NewView(editor.NewBufferFromString(""))
	}
	header.SetDoneFunc(func(s string) {
		num, _ := strconv.ParseUint(s, 10, 32)
		if uint32(num) >= (*area).GetCount() {
			a.App.SetFocus(body)
		} else {
			if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), num)) {
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), num))
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			} else {
				a.Pages.AddPage(a.ViewMsg(area, uint32(num)))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), num))
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			}
		}
	})

	body.Readonly = true
	body.SetDoneFunc(func() {
		a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
		a.SwitchToAreaListPage()
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var area = a.CurrentArea
		if event.Key() == tcell.KeyF1 {
			a.Pages.AddPage(a.ViewMsgHelp())
		} else if event.Key() == tcell.KeyRight {
			if msgNum == (*area).GetCount() {
	                        a.RefreshAreaList()
                                a.CurrentArea = &msgapi.Areas[0]
				a.SwitchToAreaListPage()
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum+1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum+1))
					go (func() {
						a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
					})()
				} else {
					a.Pages.AddPage(a.ViewMsg(area, msgNum+1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum+1))
					go (func() {
						a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
					})()
				}
			}
		} else if event.Key() == tcell.KeyLeft {
			if msgNum <= 1 {
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				a.SwitchToAreaListPage()
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum-1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum-1))
					go (func() {
						a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
					})()
				} else {
					a.Pages.AddPage(a.ViewMsg(area, msgNum-1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum-1))
					go (func() {
						a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
					})()
				}
			}
		} else if event.Key() == tcell.KeyInsert || event.Key() == tcell.KeyCtrlI {
			a.Pages.AddPage(a.InsertMsg(area, 0))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", (*area).GetName()))
		} else if msg == nil {
			return event
		} else if event.Key() == tcell.KeyCtrlK || (event.Rune() == 'k' && event.Modifiers()&tcell.ModAlt > 0) {
			a.showKludges = !a.showKludges
			//body.SetText(msg.ToView(a.showKludges))
			body.OpenBuffer(editor.NewBufferFromString(msg.ToView(a.showKludges)))
		} else if event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyF3 || (event.Rune() == 'q') {
			a.Pages.AddPage(a.InsertMsg(area, newMsgTypeAnswer))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", (*area).GetName()))
		} else if event.Key() == tcell.KeyCtrlN || (event.Rune() == 'n' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(area, newMsgTypeAnswerNewArea))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyCtrlF || (event.Rune() == 'f' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(area, newMsgTypeForward))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyDelete {
			a.Pages.AddPage(a.showDelMsg(area, msgNum))
			a.Pages.ShowPage("DelMsgModal")
		} else if event.Key() == tcell.KeyCtrlL || event.Rune() == 'l' {
			a.Pages.AddPage(a.showMessageList(area))
			a.Pages.ShowPage("MessageListModal")
		} else if event.Key() == tcell.KeyCtrlG || event.Rune() == 'g' {
			a.App.SetFocus(header)
			//a.Pages.AddPage(a.showMessageList(area))
			//a.Pages.ShowPage("MessageListModal")
		} else if event.Rune() == '<' {
			if msgNum != 1 {
				a.Pages.AddPage(a.ViewMsg(area, 1))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), 1))
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			}
		} else if event.Rune() == '>' {
			if msgNum != (*area).GetCount() {
				a.Pages.AddPage(a.ViewMsg(area, (*area).GetCount()))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), (*area).GetCount()))
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			}
		}

		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(body, 0, 1, true)
	return fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum), layout, true, true
}

func (a *App) showMessageList(area *msgapi.AreaPrimitive) (string, tview.Primitive, bool, bool) {
	modal := NewModalMessageList(area).
		SetDoneFunc(func(msgNum uint32) {
			a.Pages.HidePage("MessageListModal")
			a.Pages.RemovePage("MessageListModal")
			a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), (*area).GetLast()))
			a.Pages.AddPage(a.ViewMsg(area, msgNum))
			a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
			a.App.SetFocus(a.Pages)
		})
	return "MessageListModal", modal, true, true
}

func (a *App) showAreaList(area *msgapi.AreaPrimitive, newMsgType int) (string, tview.Primitive, bool, bool) {
	modal := NewModalAreaList().
		SetDoneFunc(func(buttonIndex int) {
			a.im.postArea = area
			a.Pages.HidePage("AreaListModal")
			a.Pages.RemovePage("AreaListModal")
			a.Pages.AddPage(a.InsertMsg(area, newMsgType))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", (*area).GetName()))
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
func (a *App) showDelMsg(area *msgapi.AreaPrimitive, msgNum uint32) (string, tview.Primitive, bool, bool) {
	modal := NewModalMenu().
		SetY(6).
		SetText("Delete?").
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int) {
			a.Pages.HidePage("DelMsgModal")
			a.Pages.RemovePage("DelMsgModal")
			if buttonIndex == 0 {
				(*area).DelMsg(msgNum)
				a.Pages.AddPage(a.ViewMsg(area, msgNum-1))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum-1))
				go (func() {
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", (*area).GetName(), msgNum))
				})()
			}
			a.App.SetFocus(a.Pages)
		})
	return "DelMsgModal", modal, true, true
}
