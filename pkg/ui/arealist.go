package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	// "log"
	"strconv"
)

// AreaListQuit exit app
func (a *App) AreaListQuit() (string, tview.Primitive, bool, bool) {
	modal := NewModalMenu().
		SetText("Quit GOssipEd?").
		AddButtons([]string{
			"    Quit   ",
			"   Cancel  ",
		}).
		SetDoneFunc(func(buttonIndex int) {
			if buttonIndex == 0 {
				a.App.Stop()
			} else {
				a.Pages.HidePage("AreaListQuit")
				a.App.SetFocus(a.al)
				//a.Pages.SwitchToPage("AreaList")
			}
		})
	return "AreaListQuit", modal, false, false
}

// AreaList - arealist widget
func (a *App) AreaList() (string, tview.Primitive, bool, bool) {
	searchString := NewSearchString()
	a.al = tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetSelectionChangedFunc(func(row int, column int) {
			if row < 1 {
				row = 1
			}
			a.sb.SetStatus(fmt.Sprintf("%s: %d msgs, %d unread",
				msgapi.Areas[row-1].GetName(),
				msgapi.Areas[row-1].GetCount(),
				msgapi.Areas[row-1].GetCount()-msgapi.Areas[row-1].GetLast(),
			))
		})
	a.al.SetSelectedFunc(func(row int, column int) {
		a.onSelected(row, column)
	})
	a.al.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue)
	a.al.SetSelectedStyle(tcell.ColorWhite, tcell.ColorNavy, tcell.AttrBold)
	a.al.SetCell(
		0, 0, tview.NewTableCell(" Area").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false))
	a.al.SetCell(
		0, 1, tview.NewTableCell("EchoID").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1).
			SetSelectable(false))
	a.al.SetCell(
		0, 2, tview.NewTableCell("Msgs").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	a.al.SetCell(
		0, 3, tview.NewTableCell("   New").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	a.al.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyEsc:
			searchString.Clear()
			a.Pages.ShowPage("AreaListQuit")
		case tcell.KeyRight:
			searchString.Clear()
			a.onSelected(a.al.GetSelection())
		case tcell.KeyDown, tcell.KeyUp, tcell.KeyEnter:
			searchString.Clear()
		case tcell.KeyRune:
			searchString.AddChar(event.Rune())
			row := msgapi.Search(searchString.GetText())
			if row > 0 {
				a.al.Select(row, 0)
			}
		}
		return event
	})
	for i, ar := range msgapi.Areas {
		if ar.GetCount()-ar.GetLast() > 0 {
			a.al.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i), 10)+"[::b]+").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		} else {
			a.al.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i), 10)+" ").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		}
		a.al.SetCell(i+1, 1, tview.NewTableCell(ar.GetName()).SetTextColor(tcell.ColorSilver))
		a.al.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()), 10)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		a.al.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()-ar.GetLast()), 10)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
	}
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchString, 1, 1, false).
		AddItem(a.al, 0, 1, true)
	return "AreaList", layout, true, true
}
func (a *App) onSelected(row int, column int) {
	if row < 1 {
		row = 1
	}
	if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast())) {
		a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast()))
	} else {
		a.Pages.AddPage(a.ViewMsg(row-1, msgapi.Areas[row-1].GetLast()))
		a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast()))
	}
}
