package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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

func initAreaListHeader(a *App) {
	a.al.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue)
	a.al.SetSelectedStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy).Bold(true))
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
}

func (a *App) RefreshAreaList() {
	var currentArea = ""
	if a.CurrentArea != nil {
		currentArea = (*a.CurrentArea).GetName()
	}
	refreshAreaList(a, currentArea)
}

func refreshAreaList(a *App, currentArea string) {
	msgapi.SortAreas()
	a.al.Clear()
	initAreaListHeader(a)
	var selectIndex = -1
	for i, ar := range msgapi.Areas {
		var areaStyle = " "
		if msgapi.AreaHasUnreadMessages(&ar) {
			areaStyle = "[::b]+"
		}
		a.al.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i), 10)+areaStyle).
			SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		a.al.SetCell(i+1, 1, tview.NewTableCell(ar.GetName()).SetTextColor(tcell.ColorSilver))
		a.al.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()), 10)).
			SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		a.al.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()-ar.GetLast()), 10)).
			SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		if currentArea != "" && currentArea == ar.GetName() {
			selectIndex = i + 1
		}
	}
	if selectIndex != -1 {
		a.al.Select(selectIndex, 0)
	}

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
			var area = msgapi.Areas[row-1]
			a.sb.SetStatus(fmt.Sprintf("%s: %d msgs, %d unread",
				area.GetName(),
				area.GetCount(),
				area.GetCount()-area.GetLast(),
			))
		})
	a.al.SetSelectedFunc(func(row int, column int) {
		a.onSelected(row, column)
	})
	a.al.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyEsc:
			searchString.Clear()
			a.Pages.ShowPage("AreaListQuit")
		case tcell.KeyF1:
			a.Pages.ShowPage("AreaListHelp")
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
	refreshAreaList(a, "")
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
	a.CurrentArea = &msgapi.Areas[row-1]
	if a.Pages.HasPage(fmt.Sprintf("ViewMsg-%s-%d", (*a.CurrentArea).GetName(), (*a.CurrentArea).GetLast())) {
		a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*a.CurrentArea).GetName(), (*a.CurrentArea).GetLast()))
	} else {
		a.Pages.AddPage(a.ViewMsg(a.CurrentArea, (*a.CurrentArea).GetLast()))
		a.Pages.SwitchToPage(fmt.Sprintf("ViewMsg-%s-%d", (*a.CurrentArea).GetName(), (*a.CurrentArea).GetLast()))
	}
}
