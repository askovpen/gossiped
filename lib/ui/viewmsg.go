package ui

import(
  "fmt"
  "github.com/askovpen/goated/lib/msgapi"
)

func viewMsg(areaId int, msgNum uint32) {
  maxX, maxY := App.Size()
  msg, _:=msgapi.Areas[areaId].GetMsg(msgNum)
  MsgHeader, _:= App.SetView("MsgHeader", 0, 0, maxX-1, 6);
  MsgHeader.Frame = true
  MsgBody, _:= App.SetView("MsgBody", -1, 6, maxX, maxY-1);
  MsgBody.Frame = false
  fmt.Fprintf(MsgBody, "%s",msg.ToView(true))
}
