package ui

import (
	"fmt"
	"github.com/askovpen/goated/pkg/config"
	"github.com/askovpen/goated/pkg/msgapi"
	"github.com/askovpen/goated/pkg/types"
	"github.com/askovpen/gocui"
	//"log"
	"strings"
)

func answerMsgNewArea(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	newMsgAreaID = oy + cy - 1
	g.DeleteView("iAreaList")
	editMsg(g, v)
	return nil
}

func answerMsg(g *gocui.Gui, v *gocui.View) error {
	newMsgType = newMsgTypeAnswer
	newMsgAreaID = curAreaID
	err := editMsg(g, v)
	if err != nil {
		return err
	}
	return nil
}

func editMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	var origMessage *msgapi.Message
	if newMsgType == 0 {
		newMsgAreaID = curAreaID
	}
	if newMsg == nil {
		newMsg = &msgapi.Message{From: config.Config.Username, FromAddr: config.Config.Address, AreaID: newMsgAreaID}
		newMsg.Kludges = make(map[string]string)
		newMsg.Kludges["PID:"] = config.PID
		newMsg.Kludges["CHRS:"] = config.Config.Chrs.Default
		if msgapi.Areas[newMsgAreaID].GetChrs() != "" {
			newMsg.Kludges["CHRS:"] = msgapi.Areas[newMsgAreaID].GetChrs()
		}
	}
	if (newMsgType & newMsgTypeAnswer) != 0 {
		origMessage, _ = msgapi.Areas[curAreaID].GetMsg(curMsgNum)
		newMsg.To = origMessage.From
		newMsg.ToAddr = origMessage.FromAddr
		newMsg.Kludges["REPLY:"] = origMessage.Kludges["MSGID:"]
	} else if (newMsgType & newMsgTypeForward) != 0 {
		origMessage, _ = msgapi.Areas[curAreaID].GetMsg(curMsgNum)
	}
	maxX, maxY := g.Size()
	msgHeader, _ := g.SetView("MsgHeader", 0, 0, maxX-1, 5)
	msgHeader.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
	msgHeader.FrameBgColor = gocui.ColorBlack
	msgHeader.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
	msgHeader.Title = msgapi.Areas[newMsgAreaID].GetName()
	msgHeader.Clear()
	fmt.Fprintf(msgHeader, " Msg  : %-34s Pvt\n",
		fmt.Sprintf("%d of %d", msgapi.Areas[newMsgAreaID].GetCount()+1, msgapi.Areas[newMsgAreaID].GetCount()+1))
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
	if (newMsgType & newMsgTypeAnswer) != 0 {
		fmt.Fprintf(msgToName, "%s", origMessage.From)
	} else if msgapi.Areas[newMsgAreaID].GetType() != msgapi.EchoAreaTypeNetmail {
		fmt.Fprintf(msgToName, "All")
	}
	msgToAddr, _ := g.SetView("editToAddr", 43, 2, 57, 4)
	msgToAddr.Clear()
	msgToAddr.Frame = false
	msgToAddr.Editable = true
	if (newMsgType&newMsgTypeAnswer) != 0 && msgapi.Areas[newMsgAreaID].GetType() == msgapi.EchoAreaTypeNetmail {
		fmt.Fprintf(msgToAddr, "%s", origMessage.FromAddr)
	}
	g.Cursor = true
	App.SetCurrentView("editFromName")
	msgSubj, _ := g.SetView("editSubj", 8, 3, 60, 5)
	msgSubj.Clear()
	msgSubj.Frame = false
	msgSubj.Editable = true
	if (newMsgType & newMsgTypeAnswer) != 0 {
		fmt.Fprintf(msgSubj, "%s", origMessage.Subject)
		ActiveWindow = "editSubj"
	} else if (newMsgType & newMsgTypeForward) != 0 {
		fmt.Fprintf(msgSubj, "%s", origMessage.Subject)
		ActiveWindow = "editToName"
	} else {
		ActiveWindow = "editToName"
	}
	return nil
}

func editToNameNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editToName")
	newMsg.To = strings.Trim(vn.Buffer(), "\n")
	if msgapi.Areas[curAreaID].GetType() == msgapi.EchoAreaTypeNetmail {
		ActiveWindow = "editToAddr"
	} else {
		ActiveWindow = "editSubj"
	}
	return nil
}

func editFromNameNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editFromName")
	newMsg.From = strings.Trim(vn.Buffer(), "\n")
	ActiveWindow = "editFromAddr"
	return nil
}

func editToAddrNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editToAddr")
	newMsg.ToAddr = types.AddrFromString(strings.Trim(vn.Buffer(), "\n"))
	ActiveWindow = "editSubj"
	return nil
}

func editFromAddrNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editFromAddr")
	newMsg.FromAddr = types.AddrFromString(strings.Trim(vn.Buffer(), "\n"))
	ActiveWindow = "editToName"
	return nil
}

func editToSubjNext(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editSubj")
	newMsg.Subject = strings.Trim(vn.Buffer(), "\n")
	ActiveWindow = "editFromName"
	return nil
}

func editToSubjBody(g *gocui.Gui, v *gocui.View) error {
	vn, _ := g.View("editSubj")
	newMsg.Subject = strings.Trim(vn.Buffer(), "\n")
	var origMessage *msgapi.Message
	var p int
	var mv string
	if (newMsgType&newMsgTypeAnswer) != 0 || (newMsgType&newMsgTypeForward) != 0 {
		origMessage, _ = msgapi.Areas[curAreaID].GetMsg(curMsgNum)
	}
	vn, _ = g.View("editMsgBody")
	if (newMsgType & newMsgTypeAnswer) != 0 {
		mv, p = newMsg.ToEditAnswerView(origMessage)
	} else if (newMsgType & newMsgTypeForward) != 0 {
		mv, p = newMsg.ToEditForwardView(origMessage)
	} else {
		mv, p = newMsg.ToEditNewView()
	}
	_, maxY := vn.Size()
	if p > maxY-1 {
		vn.SetCursor(0, maxY-1)
		vn.SetOrigin(0, p-maxY-1)
	} else {
		vn.SetCursor(0, p)
	}
	if vn.Buffer() == "" {
		fmt.Fprintf(vn, mv)
	}
	ActiveWindow = "editMsgBody"
	return nil
}

func editMsgBodyMenu(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	vn, _ := g.View("editMsgBody")
	newMsg.Body = string(vn.Buffer())
	v, _ = App.SetView("editMenuMsg", 0, 6, 19, 11)
	v.Title = "Save?"
	v.Highlight = true
	v.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
	v.FrameFgColor = gocui.ColorRed | gocui.AttrBold
	v.FrameBgColor = gocui.ColorBlack
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	fmt.Fprintf(v, " Yes!             \n No, Drop         \n Continue Writing \n Edit Header      ")
	ActiveWindow = "editMenuMsg"
	return nil
}
func editMsgBodyMenuUp(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy == 0 {
		cy = 3
	} else {
		cy--
	}
	v.SetCursor(cx, cy)
	return nil
}
func editMsgBodyMenuDown(g *gocui.Gui, v *gocui.View) error {
	cx, cy := v.Cursor()
	if cy == 3 {
		cy = 0
	} else {
		cy++
	}
	v.SetCursor(cx, cy)
	return nil
}
func saveMessage(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	g.Cursor = false
	switch cy {
	case 0:
		g.DeleteView("MsgHeader")
		g.DeleteView("editMsgBody")
		g.DeleteView("editFromName")
		g.DeleteView("editFromAddr")
		g.DeleteView("editToName")
		g.DeleteView("editToAddr")
		g.DeleteView("editSubj")
		g.DeleteView("editMenuMsg")
		err := msgapi.Areas[newMsgAreaID].SaveMsg(newMsg.MakeBody())
		newMsgType = 0
		newMsg = nil
		if err != nil {
			errorMsg(err.Error(), "AreaList")
		} else {
			viewMsg(curAreaID, curMsgNum)
			ActiveWindow = "MsgBody"
		}
	case 1:
		newMsg = nil
		newMsgType = 0
		g.DeleteView("MsgHeader")
		g.DeleteView("editMsgBody")
		g.DeleteView("editFromName")
		g.DeleteView("editFromAddr")
		g.DeleteView("editToName")
		g.DeleteView("editToAddr")
		g.DeleteView("editSubj")
		g.DeleteView("editMenuMsg")
		viewMsg(curAreaID, curMsgNum)
		ActiveWindow = "MsgBody"
	case 2:
		g.DeleteView("editMenuMsg")
		g.Cursor = true
		ActiveWindow = "editMsgBody"
	case 3:
		g.DeleteView("editMenuMsg")
		g.Cursor = true
		ActiveWindow = "editFromName"
	}
	return nil
}
