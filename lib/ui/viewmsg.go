package ui

import (
	"fmt"
	"github.com/askovpen/goated/lib/msgapi"
	"github.com/askovpen/gocui"
	"log"
	"strconv"
	"strings"
)

func viewMsg(areaId int, msgNum uint32) error {
	if msgNum == 0 && msgapi.Areas[areaId].GetCount() != 0 {
		msgNum = 1
	}
	msgEmpty := false
	curAreaId = areaId
	curMsgNum = msgNum
	maxX, maxY := App.Size()
	msg, err := msgapi.Areas[areaId].GetMsg(msgNum)
	if err != nil {
		if err.Error() == "Empty Area" {
			msgEmpty = true
		} else {
			log.Printf("vM err: %s", err.Error())
			return err
		}
	}
	StatusLine = fmt.Sprintf("Msg %d of %d (%d left)",
		msgNum,
		msgapi.Areas[areaId].GetCount(),
		msgapi.Areas[areaId].GetCount()-msgNum)
	if msgEmpty {
		MsgHeader, _ := App.SetView("MsgHeader", 0, 0, maxX-1, 5)
		MsgHeader.Clear()
		MsgHeader.Title = msgapi.Areas[areaId].GetName()
		MsgHeader.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
		MsgHeader.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
		fmt.Fprintf(MsgHeader, " Msg  : %-34s %-36s\n",
			fmt.Sprintf("%d of %d", msgNum, msgapi.Areas[areaId].GetCount()),
			"")
		fmt.Fprintf(MsgHeader, " From :\n")
		fmt.Fprintf(MsgHeader, " To   : \n")
		fmt.Fprintf(MsgHeader, " Subj : ")
		MsgBody, _ := App.SetView("MsgBody", -1, 5, maxX, maxY-1)
		MsgBody.Frame = false
		MsgBody.Wrap = true
		MsgBody.Clear()
	} else {
		msgapi.Areas[areaId].SetLast(msgNum)
		MsgHeader, _ := App.SetView("MsgHeader", 0, 0, maxX-1, 5)
		MsgHeader.Wrap = false
		MsgHeader.Clear()
		MsgHeader.Title = msgapi.Areas[areaId].GetName()
		MsgHeader.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
		MsgHeader.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
		fmt.Fprintf(MsgHeader, " Msg  : %-34s %-36s\n",
			fmt.Sprintf("%d of %d", msgNum, msgapi.Areas[areaId].GetCount()),
			strings.Join(msg.Attrs, " "))
		fmt.Fprintf(MsgHeader, " From : %-34s %-15s %-18s\n",
			msg.From,
			msg.FromAddr.String(),
			msg.DateWritten.Format("02 Jan 06 15:04:05"))
		fmt.Fprintf(MsgHeader, " To   : %-34s %-15s %-18s\n",
			msg.To,
			msg.ToAddr.String(),
			msg.DateArrived.Format("02 Jan 06 15:04:05"))
		corrupted := ""
		if msg.Corrupted {
			corrupted = "\033[31;1mCorrupted\033[0m"
		}
		fmt.Fprintf(MsgHeader, " Subj : %-59s %9s",
			msg.Subject, corrupted)
		MsgBody, _ := App.SetView("MsgBody", -1, 5, maxX, maxY-1)
		MsgBody.Frame = false
		MsgBody.Wrap = true
		MsgBody.Clear()
		fmt.Fprintf(MsgBody, "%s", msg.ToView(showKludges))
	}
	return nil
}
func prevMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	if curMsgNum > 1 {
		err := viewMsg(curAreaId, curMsgNum-1)
		if err != nil {
			errorMsg(err.Error(), "AreaList")
			return nil
		}
		ActiveWindow = "MsgBody"
	}
	return nil
}

func nextMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	if curMsgNum < msgapi.Areas[curAreaId].GetCount() {
		err := viewMsg(curAreaId, curMsgNum+1)
		if err != nil {
			errorMsg(err.Error(), "AreaList")
			return nil
		}
		ActiveWindow = "MsgBody"
	}
	return nil
}

func firstMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	viewMsg(curAreaId, 1)
	ActiveWindow = "MsgBody"
	return nil
}
func editMsgNumEnter(g *gocui.Gui, v *gocui.View) error {
	g.Cursor = false
	ActiveWindow = "MsgBody"
	en, err := g.View("editNumber")
	if err != nil {
		return err
	}

	n, err := strconv.ParseInt(strings.Trim(en.Buffer(), "\n"), 10, 32)
	if err != nil {
		n = -1
	}
	if err := g.DeleteView("editNumber"); err != nil {
		return err
	}
	if err := g.DeleteView("editNumberTitle"); err != nil {
		return err
	}
	if n > 0 && uint32(n) <= msgapi.Areas[curAreaId].GetCount() {
		err := viewMsg(curAreaId, uint32(n))
		if err != nil {
			errorMsg(err.Error(), "AreaList")
		}
	}
	return nil
}

func editMsgNum(g *gocui.Gui, v *gocui.View) error {
	editableNumber, err := App.SetView("editNumber", 8, 0, 15, 2)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	editableNumber.Frame = false
	editableNumber.BgColor = gocui.ColorWhite
	editableNumber.FgColor = gocui.ColorBlue
	editableNumber.Editable = true
	editableNumberTitle, _ := App.SetView("editNumberTitle", 14, 0, 24, 2)
	editableNumberTitle.Frame = false
	g.Cursor = true
	fmt.Fprintf(editableNumberTitle, " of %d", msgapi.Areas[curAreaId].GetCount())
	App.SetCurrentView("editNumber")
	ActiveWindow = "editNumber"
	return nil
}

func lastMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	viewMsg(curAreaId, msgapi.Areas[curAreaId].GetCount())
	ActiveWindow = "MsgBody"
	return nil
}

func quitMsgView(g *gocui.Gui, v *gocui.View) error {
	ActiveWindow = "AreaList"
	if err := g.DeleteView("MsgHeader"); err != nil {
		return err
	}
	if err := g.DeleteView("MsgBody"); err != nil {
		return err
	}
	return nil
}
func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, sy := v.Size()
		if oy >= len(v.BufferLines())-sy-1 {
			return nil
		}
		if err := v.SetOrigin(ox, oy+1); err != nil {
			return err
		}
	}
	return nil
}

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if oy == 0 {
			return nil
		}
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}
func toggleKludges(g *gocui.Gui, v *gocui.View) error {
	showKludges = !showKludges
	quitMsgView(g, v)
	viewMsg(curAreaId, curMsgNum)
	ActiveWindow = "MsgBody"
	return nil
}
