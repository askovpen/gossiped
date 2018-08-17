package ui

import (
	"fmt"
	"github.com/askovpen/goated/lib/config"
	"github.com/askovpen/goated/lib/msgapi"
	"github.com/askovpen/goated/lib/types"
	"github.com/jroimartin/gocui"
	"log"
)

func editMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	if newMsg == nil {
		newMsg = &msgapi.Message{From: config.Config.Username, FromAddr: config.Config.Address}
	}
	maxX, maxY := g.Size()
	msgHeader, _ := g.SetView("MsgHeader", 0, 0, maxX-1, 5)
	msgHeader.Clear()
	fmt.Fprintf(msgHeader, " Msg  : %-34s Pvt\n",
		fmt.Sprintf("%d of %d", curMsgNum+1, msgapi.Areas[curAreaId].GetCount()+1))
	fmt.Fprintf(msgHeader, " From :\n")
	fmt.Fprintf(msgHeader, " To   : \n")
	fmt.Fprintf(msgHeader, " Subj : ")
	msgBody, _ := g.SetView("editMsgBody", -1, 5, maxX, maxY-1)
	msgBody.Frame = false
	msgBody.Wrap = true
	msgBody.Editable = true
	msgBody.Clear()
	msgFromName, _ := g.SetView("editFromName", 8, 1, 42, 3)
	msgFromName.Clear()
	msgFromName.Frame = false
	msgFromName.Editable = true
	fmt.Fprintf(msgFromName, "%s", config.Config.Username)
	msgFromAddr, _ := g.SetView("editFromAddr", 43, 1, 65, 3)
	msgFromAddr.Clear()
	msgFromAddr.Frame = false
	msgFromAddr.Editable = true
	fmt.Fprintf(msgFromAddr, "%s", config.Config.Address)
	msgToName, _ := g.SetView("editToName", 8, 2, 42, 4)
	msgToName.Clear()
	msgToName.Frame = false
	msgToName.Editable = true
	fmt.Fprintf(msgToName, "All")
	msgToAddr, _ := g.SetView("editToAddr", 43, 2, 57, 4)
	msgToAddr.Clear()
	msgToAddr.Frame = false
	msgToAddr.Editable = true
	g.Cursor = true
	App.SetCurrentView("editFromName")
	msgSubj, _ := g.SetView("editSubj", 8, 3, 60, 5)
	msgSubj.Clear()
	msgSubj.Frame = false
	msgSubj.Editable = true
	ActiveWindow = "editToName"
	return nil
}

func editToNameNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editToName")
	newMsg.To = vn.Buffer()
	ActiveWindow = "editToAddr"
	return nil
}

func editFromNameNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editFromName")
	newMsg.From = vn.Buffer()
	ActiveWindow = "editFromAddr"
	return nil
}

func editToAddrNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editToAddr")
	newMsg.ToAddr = types.AddrFromString(vn.Buffer())
	ActiveWindow = "editSubj"
	return nil
}

func editFromAddrNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editFromAddr")
	newMsg.FromAddr = types.AddrFromString(vn.Buffer())
	ActiveWindow = "editToName"
	return nil
}

func editToSubjNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editSubj")
	newMsg.Subject = vn.Buffer()
	ActiveWindow = "editFromName"
	return nil
}

func editToSubjBody(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editSubj")
	newMsg.Subject = string(vn.Buffer())
	vn, _ = g.View("editMsgBody")
	fmt.Fprintf(vn, "Hello, %s\n\n\n--- %s\n * Origin: %s (%s)", newMsg.From, config.LongPID, config.Config.Origin, config.Config.Address)
	ActiveWindow = "editMsgBody"
	return nil
}

func editMsgBodyMenu(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	vn, _ := g.View("editMsgBody")
	newMsg.Body = string(vn.Buffer())
	v, _ = App.SetView("editMenuMsg", 0, 6, 17, 11)
	v.Title = "Save?"
	v.Highlight = true
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	fmt.Fprintf(v, "Yes!\nNo, Drop\nContinue Writing\nEdit Header")
	ActiveWindow = "editMenuMsg"
	return nil
}
func saveMessage(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	log.Printf("cy %d", cy)
	if cy == 0 {
		g.DeleteView("MsgHeader")
		g.DeleteView("editMsgBody")
		g.DeleteView("editFromName")
		g.DeleteView("editFromAddr")
		g.DeleteView("editToName")
		g.DeleteView("editToAddr")
		g.DeleteView("editSubj")
		g.DeleteView("editMenuMsg")
		err := msgapi.Areas[curAreaId].SaveMsg(newMsg)
		if err != nil {
			errorMsg(err.Error(), "AreaList")
		}
	}
	return nil
}
