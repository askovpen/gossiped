package ui

import(
  "fmt"
  "github.com/askovpen/goated/lib/msgapi"
  "github.com/jroimartin/gocui"
  "log"
)

func viewMsg(areaId int, msgNum uint32) error {
  if msgNum==0 && msgapi.Areas[areaId].GetCount()!=0 {
    msgNum=1
  }
  curAreaId=areaId
  curMsgNum=msgNum
  maxX, maxY := App.Size()
  msg, err:=msgapi.Areas[areaId].GetMsg(msgNum)
  if err!=nil {
    log.Printf("vM err: %s",err.Error())
    return err
  }
  StatusLine=fmt.Sprintf("Msg %d of %d (%d left)",
    msgNum,
    msgapi.Areas[areaId].GetCount(),
    msgapi.Areas[areaId].GetCount()-msgNum)
  msgapi.Areas[areaId].SetLast(msgNum)
  MsgHeader, _:= App.SetView("MsgHeader", 0, 0, maxX-1, 5);
  MsgHeader.Wrap = false
  MsgHeader.Title=msgapi.Areas[areaId].GetName()
  fmt.Fprintf(MsgHeader, " Msg  : %-34s %-36s\n",
    fmt.Sprintf("%d of %d", msgNum, msgapi.Areas[areaId].GetCount()),"Pvt")
  fmt.Fprintf(MsgHeader, " From : %-34s %-15s %-18s\n",
    msg.From,
    msg.FromAddr.String(),
    msg.DateWritten.Format("02 Jan 06 15:04:05"))
  fmt.Fprintf(MsgHeader, " To   : %-34s %-15s %-18s\n",
      msg.To,
      msg.ToAddr.String(),
      msg.DateArrived.Format("02 Jan 06 15:04:05"))
  fmt.Fprintf(MsgHeader, " Subj : %-50s",
      msg.Subject)
  MsgBody, _:= App.SetView("MsgBody", -1, 5, maxX, maxY-1);
  MsgBody.Frame = false
  fmt.Fprintf(MsgBody, "%s",msg.ToView(showKludges))
  return nil
}
func prevMsg(g *gocui.Gui, v *gocui.View) error {
  quitMsgView(g,v)
  if curMsgNum>1 {
    err:=viewMsg(curAreaId, curMsgNum-1)
    if err!=nil {
      errorMsg(err.Error(),"AreaList")
    }
    ActiveWindow="MsgBody"
  }
  return nil
}

func nextMsg(g *gocui.Gui, v *gocui.View) error {
  quitMsgView(g,v)
  if curMsgNum<msgapi.Areas[curAreaId].GetCount() {
    viewMsg(curAreaId, curMsgNum+1)
    ActiveWindow="MsgBody"
  }
  return nil
}

func firstMsg(g *gocui.Gui, v *gocui.View) error {
  quitMsgView(g,v)
  viewMsg(curAreaId, 1)
  ActiveWindow="MsgBody"
  return nil
}

func lastMsg(g *gocui.Gui, v *gocui.View) error {
  quitMsgView(g,v)
  viewMsg(curAreaId, msgapi.Areas[curAreaId].GetCount())
  ActiveWindow="MsgBody"
  return nil
}

func quitMsgView(g *gocui.Gui, v *gocui.View) error {
  log.Printf("Delete")
  ActiveWindow="AreaList"
  if err := g.DeleteView("MsgHeader"); err != nil {
    return err
  }
  if err := g.DeleteView("MsgBody"); err != nil {
    return err
  }
  return nil
}
func scrollDown(g *gocui.Gui, v *gocui.View) error {
  log.Printf("test")
  if v != nil {
    ox, oy := v.Origin()
    _,sy:=v.Size()
    if oy>=len(v.BufferLines())-sy-1 {
      return nil
    }
    if err := v.SetOrigin(ox, oy+1); err != nil {
      return err
    }
  }
  return nil
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
  log.Printf("test")
  if v != nil {
    ox, oy := v.Origin()
    if oy==0 {
      return nil
    }
    if err := v.SetOrigin(ox, oy-1); err != nil {
      return err
    }
  }
  return nil
}
func toggleKludges(g *gocui.Gui, v *gocui.View) error {
  log.Printf("togglekl")
  showKludges=!showKludges
  quitMsgView(g,v)
  viewMsg(curAreaId, curMsgNum)
  ActiveWindow="MsgBody"
  return nil
}
