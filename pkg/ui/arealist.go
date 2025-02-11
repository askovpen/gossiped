package ui

import (
	"fmt"
	"strconv"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	borderStyle := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementBorder)
	headerStyle := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementHeader)
	fgHeader, bgHeader, attrHeader := headerStyle.Decompose()
	selStyle := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementSelection)
	a.al.SetBorder(true).
		SetBorderStyle(borderStyle)
	a.al.SetSelectedStyle(selStyle)
	a.al.SetCell(
		0, 0, tview.NewTableCell(" Area").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false))
	a.al.SetCell(
		0, 1, tview.NewTableCell("EchoID").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetExpansion(1).
			SetSelectable(false))
	a.al.SetCell(
		0, 2, tview.NewTableCell("Msgs").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	a.al.SetCell(
		0, 3, tview.NewTableCell("   New").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
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
	styleItem := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementItem)
	styleHighligt := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementHighlight)
	fgItem, bgItem, attrItem := styleItem.Decompose()
	fgHigh, bgHigh, attrHigh := styleHighligt.Decompose()
	var selectIndex = -1
	for i, ar := range msgapi.Areas {
		fg, bg, attr := fgItem, bgItem, attrItem
		areaStyle := ""
		if msgapi.AreaHasUnreadMessages(&ar) {
			areaStyle = "+"
			fg, bg, attr = fgHigh, bgHigh, attrHigh
		}
		a.al.SetCell(i+1, 0, tview.NewTableCell(areaStyle+strconv.FormatInt(int64(i), 10)).
			SetAlign(tview.AlignRight).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		a.al.SetCell(i+1, 1, tview.NewTableCell(ar.GetName()).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		a.al.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()), 10)).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr).
			SetAlign(tview.AlignRight))
		a.al.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()-ar.GetLast()), 10)).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr).
			SetAlign(tview.AlignRight))
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
	_, defBg, _ := config.StyleDefault.Decompose()
	a.al.SetBackgroundColor(defBg)
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
