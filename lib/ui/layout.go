package ui

import (
//  "fmt"
//  "github.com/askovpen/goated/lib/msgapi"
  "github.com/jroimartin/gocui"
//  "strconv"
  "log"
)

func setCurrentViewOnTop(name string) (*gocui.View, error) {
  if _, err := App.SetCurrentView(name); err != nil {
    return nil, err
  }
  return App.SetViewOnTop(name)
}
/*
func getAreaNew(m msgapi.AreaPrimitive) string {
  if m.GetCount()-m.GetLast()>0 {
    return "\033[37;1m+\033[0m"
  } else {
    return " "
  }
}*/
func Layout(g *gocui.Gui) error {
  maxX, maxY := g.Size()
  Status, err:= g.SetView("status", -1, maxY-2, maxX, maxY);
  if err!=nil && err!=gocui.ErrUnknownView { 
    return err
  }
  Status.Frame=false
  Status.Wrap=false
  Status.BgColor=gocui.ColorBlue
  Status.FgColor=gocui.ColorWhite|gocui.AttrBold
  Status.Clear()
//  fmt.Fprintf(Status," Loading...")
  err=CreateAreaList()
  if err!=nil {
    log.Print(err)
  }
/*  AreaList, err:= g.SetView("AreaList", 0, 0, maxX-1, maxY-2);
  if err!=nil && err!=gocui.ErrUnknownView { 
    return err
  }
  AreaList.Wrap=false
  AreaList.Highlight = true
  AreaList.SelBgColor = gocui.ColorBlue
  AreaList.SelFgColor = gocui.ColorWhite | gocui.AttrBold
  AreaList.Clear()
    fmt.Fprintf(AreaList, "\033[33;1m Area %-"+strconv.FormatInt(int64(maxX-23),10)+"s %6s %6s \033[0m\n",
    "EchoID","Msgs","New")
  for i, a := range msgapi.Areas {
    fmt.Fprintf(AreaList, "%4d%s %-"+strconv.FormatInt(int64(maxX-23),10)+"s %6d %6d \n",
      i+1,
      getAreaNew(a),
      a.GetName(),
      a.GetCount(),
      a.GetCount()-a.GetLast())
  }
  AreaList.SetCursor(0,1)
  if _, err = setCurrentViewOnTop("AreaList"); err != nil {
    return err
  }*/
  return nil
}
