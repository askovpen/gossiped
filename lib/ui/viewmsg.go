package ui

import(
  "fmt"
  "github.com/askovpen/goated/lib/msgapi"
  "github.com/gdamore/tcell"
  "github.com/rivo/tview"
)

func ViewMsg(areaId int, msgNum uint32) (string, tview.Primitive, bool, bool) {
  msg, err:=msgapi.Areas[areaId].GetMsg(msgNum)
  if err!=nil {
    modal := tview.NewModal().
      SetText(err.Error()).
      AddButtons([]string{"Quit"}).
      SetDoneFunc(func(buttonIndex int, buttonLabel string) {
        Pages.SwitchToPage("AreaList")
      })
    return fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum), modal, true, true
  }
  header:=tview.NewTextView().
    SetWrap(false)
  header.SetBorder(true).
    SetBorderAttributes(tcell.AttrBold).
    SetBorderColor(tcell.ColorBlue).
    SetTitle(" "+msgapi.Areas[areaId].GetName()+" ").
    SetTitleAlign(tview.AlignLeft).
    SetTitleColor(tcell.ColorYellow)
  htxt:=fmt.Sprintf(" Msg  : %-34s %-36s\n",
    fmt.Sprintf("%d of %d", msgNum, msgapi.Areas[areaId].GetCount()),"Pvt")
  htxt+=fmt.Sprintf(" From : %-34s %-15s %-18s\n",
    msg.From,
    msg.FromAddr.String(),
    msg.DateWritten.Format("02 Jan 06 15:04:05"))
  htxt+=fmt.Sprintf(" To   : %-34s %-15s %-18s\n",
    msg.To,
    msg.ToAddr.String(),
    msg.DateArrived.Format("02 Jan 06 15:04:05"))
  htxt+=fmt.Sprintf(" Subj : %-50s",
    msg.Subject)
  header.SetText(htxt)
  body:=tview.NewTextView().SetWrap(true)
  body.SetDynamicColors(true)
  body.SetText(msg.ToView(showKludges))
  body.SetDoneFunc(func(key tcell.Key) {
    if key == tcell.KeyEscape {
      Pages.SwitchToPage("AreaList")
    }
  })
  body.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyRight {
      if msgNum==msgapi.Areas[areaId].GetCount() {
        Pages.SwitchToPage("AreaList")
      } else {
        if Pages.HasPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1)) {
          Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1))
        } else {
          Pages.AddPage(ViewMsg(areaId, msgNum+1))
          Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum+1))
        }
      }
    } else if event.Key() == tcell.KeyLeft {
      if msgNum<=1 {
        Pages.SwitchToPage("AreaList")
      } else {
        if Pages.HasPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1)) {
          Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1))
        } else {
          Pages.AddPage(ViewMsg(areaId, msgNum-1))
          Pages.SwitchToPage(fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum-1))
        }
      }
    } else if event.Key() == tcell.KeyCtrlH {
      showKludges=!showKludges
    }

    return event
  })

  layout := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(header,6,1,false).
    AddItem(body,0,1,true)
  return fmt.Sprintf("%s-%d", msgapi.Areas[areaId].GetName(), msgNum), layout, true, true
}
