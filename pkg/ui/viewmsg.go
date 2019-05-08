package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
	"strings"
)

func (a *App) ViewMsg(areaId int, msgNum uint32) (string, tview.Primitive, bool, bool) {
	msg, err := msgapi.Areas[areaId].GetMsg(msgNum)
	if err != nil {
		modal := tview.NewModal().
			SetText(err.Error()).
			AddButtons([]string{"Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				a.Pages.SwitchToPage("AreaList")
			})
		return fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum), modal, true, true
	}
	msgapi.Areas[areaId].SetLast(msgNum)
	if msgapi.Areas[areaId].GetCount()-msgapi.Areas[areaId].GetLast() > 0 {
		a.al.SetCell(areaId+1, 0, tview.NewTableCell(strconv.FormatInt(int64(areaId), 10)+"[::b]+").SetAlign(tview.AlignRight))
	} else {
		a.al.SetCell(areaId+1, 0, tview.NewTableCell(strconv.FormatInt(int64(areaId), 10)+" ").SetAlign(tview.AlignRight))
	}
	a.al.SetCell(areaId+1, 2, tview.NewTableCell(strconv.FormatInt(int64(msgapi.Areas[areaId].GetCount()), 10)).SetAlign(tview.AlignRight))
	a.al.SetCell(areaId+1, 3, tview.NewTableCell(strconv.FormatInt(int64(msgapi.Areas[areaId].GetCount()-msgapi.Areas[areaId].GetLast()), 10)).SetAlign(tview.AlignRight))
	a.sb.SetStatus(fmt.Sprintf("Msg %d of %d (%d left)",
		msgNum,
		msgapi.Areas[areaId].GetCount(),
		msgapi.Areas[areaId].GetCount()-msgNum,
	))
	header := tview.NewTextView().
		SetWrap(false).
		SetTextColor(tcell.ColorSilver)
	header.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" " + msgapi.Areas[areaId].GetName() + " ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorYellow)
	repl := ""
	if msg.ReplyTo > 0 {
		repl = fmt.Sprintf("-%d ", msg.ReplyTo)
	}
	for _, rn := range msg.Replies {
		repl += fmt.Sprintf("+%d ", rn)
	}
	htxt := fmt.Sprintf(" Msg  : %-34s %-36s\n",
		fmt.Sprintf("%d of %d %s", msgNum, msgapi.Areas[areaId].GetCount(), repl), strings.Join(msg.Attrs, " "))
	htxt += fmt.Sprintf(" From : %-34s %-15s %-18s\n",
		msg.From,
		msg.FromAddr.String(),
		msg.DateWritten.Format("02 Jan 06 15:04:05"))
	htxt += fmt.Sprintf(" To   : %-34s %-15s %-18s\n",
		msg.To,
		msg.ToAddr.String(),
		msg.DateArrived.Format("02 Jan 06 15:04:05"))
	htxt += fmt.Sprintf(" Subj : %-50s",
		msg.Subject)
	header.SetText(htxt)
	body := tview.NewTextView().SetWrap(true).SetWordWrap(true).SetTextColor(tcell.ColorSilver)
	body.SetDynamicColors(true)
	body.SetText(msg.ToView(a.showKludges))
	body.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			a.Pages.SwitchToPage("AreaList")
			a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
		}
	})
	body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRight {
			if msgNum == msgapi.Areas[areaId].GetCount() {
				a.Pages.SwitchToPage("AreaList")
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
				} else {
					a.Pages.AddPage(a.ViewMsg(areaId, msgNum+1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
				}
			}
		} else if event.Key() == tcell.KeyLeft {
			if msgNum <= 1 {
				a.Pages.SwitchToPage("AreaList")
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
			} else {
				if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1)) {
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
				} else {
					a.Pages.AddPage(a.ViewMsg(areaId, msgNum-1))
					a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1))
					a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
				}
			}
		} else if event.Key() == tcell.KeyCtrlK || (event.Rune() == 'k' && event.Modifiers()&tcell.ModAlt > 0) {
			a.showKludges = !a.showKludges
			body.SetText(msg.ToView(a.showKludges))
		} else if event.Key() == tcell.KeyCtrlQ || event.Key() == tcell.KeyF3 || (event.Rune() == 'q') {
			a.Pages.AddPage(a.InsertMsg(areaId, newMsgTypeAnswer))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaId].GetName()))
		} else if event.Key() == tcell.KeyCtrlN || (event.Rune() == 'n' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(areaId, newMsgTypeAnswerNewArea))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyCtrlF || (event.Rune() == 'f' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showAreaList(areaId, newMsgTypeForward))
			a.Pages.ShowPage("AreaListModal")
		} else if event.Key() == tcell.KeyCtrlL || (event.Rune() == 'l' && event.Modifiers()&tcell.ModAlt > 0) {
			a.Pages.AddPage(a.showMessageList(areaId))
			a.Pages.ShowPage("MessageListModal")
		} else if event.Key() == tcell.KeyInsert || event.Key() == tcell.KeyCtrlI {
			a.Pages.AddPage(a.InsertMsg(areaId, 0))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaId].GetName()))
		} else if event.Rune() == '<' {
			if msgNum != 1 {
				a.Pages.AddPage(a.ViewMsg(areaId, 1))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), 1))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
			}
		} else if event.Rune() == '>' {
			if msgNum != msgapi.Areas[areaId].GetCount() {
				a.Pages.AddPage(a.ViewMsg(areaId, msgapi.Areas[areaId].GetCount()))
				a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgapi.Areas[areaId].GetCount()))
				a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
			}
		}

		return event
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 6, 1, false).
		AddItem(body, 0, 1, true)
	return fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum), layout, true, true
}

func (a *App) showMessageList(areaId int) (string, tview.Primitive, bool, bool) {
	modal := NewModalMessageList(areaId).
		SetDoneFunc(func(msgNum uint32) {
			a.Pages.HidePage("MessageListModal")
			a.Pages.RemovePage("MessageListModal")
			a.Pages.RemovePage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgapi.Areas[areaId].GetLast()))
			a.Pages.AddPage(a.ViewMsg(areaId, msgNum))
			a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[areaId].GetName(), msgNum))
			a.App.SetFocus(a.Pages)
		})
	return "MessageListModal", modal, true, true
}
func (a *App) showAreaList(areaId int, newMsgType int) (string, tview.Primitive, bool, bool) {
	modal := NewModalAreaList().
		SetDoneFunc(func(buttonIndex int) {
			a.im.postArea = buttonIndex - 1
			a.Pages.HidePage("AreaListModal")
			a.Pages.RemovePage("AreaListModal")
			a.Pages.AddPage(a.InsertMsg(areaId, newMsgType))
			a.Pages.AddPage(a.InsertMsgMenu())
			a.Pages.SwitchToPage(fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaId].GetName()))
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
