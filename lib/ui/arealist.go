package ui

import (
  "fmt"
  "github.com/gdamore/tcell"
  "github.com/rivo/tview"
  "github.com/askovpen/goated/lib/msgapi"
  "log"
  "strconv"
)

func AreaList() (string, tview.Primitive, bool, bool) {
  table:= tview.NewTable().
    SetFixed(1,0).
    SetSelectable(true, false).
    SetSelectionChangedFunc(func(row int, column int) {
      if row<1 { row=1 }
      Status.SetText(fmt.Sprintf(" [::b]%s: %d msgs, %d unread", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetCount(), msgapi.Areas[row-1].GetCount() -  msgapi.Areas[row-1].GetLast()))
    })
  table.SetSelectedFunc(func(row int, column int) {
      log.Printf("selected %d %d", row, column)
      if row<1 { row=1 }
      if Pages.HasPage(fmt.Sprintf("%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast())) {
        Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast()))
      } else {
        Pages.AddPage(ViewMsg(row-1,msgapi.Areas[row-1].GetLast()))
        Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[row-1].GetName(), msgapi.Areas[row-1].GetLast()))
      }
    })
  table.SetBorder(true).
    SetBorderAttributes(tcell.AttrBold).
    SetBorderColor(tcell.ColorBlue)
  
  table.SetCell(
    0,0,tview.NewTableCell(" Area").
      SetTextColor(tcell.ColorYellow).
      SetAttributes(tcell.AttrBold).
      SetSelectable(false))
  table.SetCell(
    0,1,tview.NewTableCell("EchoID").
      SetTextColor(tcell.ColorYellow).
      SetAttributes(tcell.AttrBold).
      SetExpansion(1).
      SetSelectable(false))
  table.SetCell(
    0,2,tview.NewTableCell("Msgs").
      SetTextColor(tcell.ColorYellow).
      SetAttributes(tcell.AttrBold).
      SetSelectable(false).
      SetAlign(tview.AlignRight))
  table.SetCell(
    0,3,tview.NewTableCell("   New").
      SetTextColor(tcell.ColorYellow).
      SetAttributes(tcell.AttrBold).
      SetSelectable(false).
      SetAlign(tview.AlignRight))
  for i, a := range msgapi.Areas {
    if a.GetCount()-a.GetLast()>0 {
      table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i),10)+"[::b]+").SetAlign(tview.AlignRight))
    } else {
      table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i),10)+" ").SetAlign(tview.AlignRight))
    }
    table.SetCell(i+1, 1, tview.NewTableCell(a.GetName()))
    table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(a.GetCount()),10)).SetAlign(tview.AlignRight))
    table.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(a.GetCount()-a.GetLast()),10)).SetAlign(tview.AlignRight))
  }
  return "AreaList", table, true, true
}
