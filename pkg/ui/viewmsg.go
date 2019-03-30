package ui

import (
	"fmt"
	"github.com/askovpen/goated/pkg/config"
	"github.com/askovpen/goated/pkg/msgapi"
	"github.com/askovpen/gocui"
	"log"
	"strconv"
	"strings"
)

func viewMsg(areaID int, msgNum uint32) error {
	if msgNum == 0 && msgapi.Areas[areaID].GetCount() != 0 {
		msgNum = 1
	}
	msgEmpty := false
	curAreaID = areaID
	curMsgNum = msgNum
	maxX, maxY := App.Size()
	msg, err := msgapi.Areas[areaID].GetMsg(msgNum)
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
		msgapi.Areas[areaID].GetCount(),
		msgapi.Areas[areaID].GetCount()-msgNum)
	if msgEmpty {
		MsgHeader, _ := App.SetView("MsgHeader", 0, 0, maxX-1, 5)
		MsgHeader.Clear()
		MsgHeader.Title = msgapi.Areas[areaID].GetName()
		MsgHeader.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
		MsgHeader.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
		MsgHeader.FrameBgColor = gocui.ColorBlack
		fmt.Fprintf(MsgHeader, " Msg  : %-34s %-36s\n",
			fmt.Sprintf("%d of %d", msgNum, msgapi.Areas[areaID].GetCount()),
			"")
		fmt.Fprintf(MsgHeader, " From :\n")
		fmt.Fprintf(MsgHeader, " To   : \n")
		fmt.Fprintf(MsgHeader, " Subj : ")
		MsgBody, _ := App.SetView("MsgBody", -1, 5, maxX, maxY-1)
		MsgBody.Frame = false
		MsgBody.Wrap = true
		MsgBody.Clear()
	} else {
		msgapi.Areas[areaID].SetLast(msgNum)
		MsgHeader, _ := App.SetView("MsgHeader", 0, 0, maxX-1, 5)
		MsgHeader.Wrap = false
		MsgHeader.Clear()
		MsgHeader.Title = msgapi.Areas[areaID].GetName()
		MsgHeader.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
		MsgHeader.TitleBgColor = gocui.ColorBlack
		MsgHeader.FrameFgColor = gocui.ColorBlue | gocui.AttrBold
		repl := ""
		if msg.ReplyTo > 0 {
			repl = fmt.Sprintf("-%d ", msg.ReplyTo)
		}
		for _, rn := range msg.Replies {
			repl += fmt.Sprintf("+%d ", rn)
		}
		fmt.Fprintf(MsgHeader, " Msg  : %-34s %-36s\n",
			fmt.Sprintf("%d of %d %s", msgNum, msgapi.Areas[areaID].GetCount(), repl),
			strings.Join(msg.Attrs, " "))
		fmt.Fprintf(MsgHeader, " From : %-34s %-15s %-18s %-40s\n",
			msg.From,
			msg.FromAddr.String(),
			msg.DateWritten.Format("02 Jan 06 15:04:05"),
			config.GetCity(msg.FromAddr.ShortString()))
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
		err := viewMsg(curAreaID, curMsgNum-1)
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
	if curMsgNum < msgapi.Areas[curAreaID].GetCount() {
		err := viewMsg(curAreaID, curMsgNum+1)
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
	viewMsg(curAreaID, 1)
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
	if n > 0 && uint32(n) <= msgapi.Areas[curAreaID].GetCount() {
		err := viewMsg(curAreaID, uint32(n))
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
	fmt.Fprintf(editableNumberTitle, " of %d", msgapi.Areas[curAreaID].GetCount())
	App.SetCurrentView("editNumber")
	ActiveWindow = "editNumber"
	return nil
}

func lastMsg(g *gocui.Gui, v *gocui.View) error {
	quitMsgView(g, v)
	viewMsg(curAreaID, msgapi.Areas[curAreaID].GetCount())
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
func scrollPgDn(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	_, sy := v.Size()
	log.Printf("vbl: %d, sy: %d, oy: %d", len(v.ViewBufferLines()), sy, oy)
	if oy >= len(v.ViewBufferLines())-sy-1 {
		return nil
	}
	if oy+sy <= len(v.ViewBufferLines())-sy-1 {
		if err := v.SetOrigin(ox, oy+sy); err != nil {
			return err
		}
	} else {
		if err := v.SetOrigin(ox, len(v.ViewBufferLines())-sy-1); err != nil {
			return err
		}
	}
	return nil
}
func scrollPgUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	_, sy := v.Size()
	if oy == 0 {
		return nil
	}
	if oy-sy >= 0 {
		if err := v.SetOrigin(ox, oy-sy); err != nil {
			return err
		}
	} else {
		if err := v.SetOrigin(ox, 0); err != nil {
			return err
		}
	}
	return nil
}
func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, sy := v.Size()
		if oy >= len(v.ViewBufferLines())-sy-1 {
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
	viewMsg(curAreaID, curMsgNum)
	ActiveWindow = "MsgBody"
	return nil
}

func listMsgs(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := App.Size()
	_, sy := v.Size()
	ml := msgapi.Areas[curAreaID].GetMessages()
	if len(*ml) == 0 {
		return nil
	}
	v, _ = App.SetView("listMsgs", 0, 1, maxX-1, maxY-2)
	v.Title = "List Messages"
	v.Highlight = true
	v.TitleFgColor = gocui.ColorYellow | gocui.AttrBold
	v.FrameFgColor = gocui.ColorRed | gocui.AttrBold
	v.FrameBgColor = gocui.ColorBlack
	v.SelBgColor = gocui.ColorBlue
	v.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	fmt.Fprintf(v, "\033[33;1m%5s  %-19s %-19s %-"+strconv.FormatInt(int64(maxX-60), 10)+"s %-10s\033[0m\n",
		"Msg", "From", "To", "Subj", "Written")
	for i, mh := range *ml {
		ch := " "
		if i == int(curMsgNum-1) {
			ch = "\033[37;1m,\033[0m"
		}
		fmt.Fprintf(v, "%5d%s %-19.19s %-19.19s %-"+strconv.FormatInt(int64(maxX-60), 10)+"."+strconv.FormatInt(int64(maxX-60), 10)+"s %-10.10s\n",
			mh.MsgNum,
			ch,
			mh.From,
			mh.To,
			mh.Subject,
			mh.DateWritten.Format("02 Jan 06"))
	}
	if int(curMsgNum) < sy+2 {
		v.SetCursor(0, int(curMsgNum))
	} else {
		v.SetOrigin(0, int(curMsgNum)-sy-2)
		v.SetCursor(0, sy+2)
	}
	ActiveWindow = "listMsgs"
	return nil
}
func upSelectMessage(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if cy > 1 {
		v.SetCursor(0, cy-1)
	} else if oy > 0 {
		v.SetOrigin(0, oy-1)
	}
	return nil
}

func pgUpSelectMessage(g *gocui.Gui, v *gocui.View) error {
	_, sy := v.Size()
	_, oy := v.Origin()
	if oy < sy-1 {
		if oy > 0 {
			v.SetOrigin(0, 0)
		} else {
			v.SetCursor(0, 1)
		}
	} else {
		v.SetOrigin(0, oy-sy+1)
	}
	return nil
}

func downSelectMessage(g *gocui.Gui, v *gocui.View) error {
	_, sy := v.Size()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if cy < sy-1 && uint32(cy) < msgapi.Areas[curAreaID].GetCount() && uint32(cy+oy) < msgapi.Areas[curAreaID].GetCount() {
		v.SetCursor(0, cy+1)
	} else if uint32(cy+oy) < msgapi.Areas[curAreaID].GetCount() {
		v.SetOrigin(0, oy+1)
	}
	return nil
}

func selectMessage(g *gocui.Gui, v *gocui.View) error {
	_, oy := v.Origin()
	_, cy := v.Cursor()
	g.DeleteView("listMsgs")
	err := viewMsg(curAreaID, uint32(oy+cy))
	if err != nil {
		errorMsg(err.Error(), "AreaList")
		return nil
	}
	ActiveWindow = "MsgBody"
	return nil
}

func pgDnSelectMessage(g *gocui.Gui, v *gocui.View) error {
	_, sy := v.Size()
	_, oy := v.Origin()
	_, cy := v.Cursor()
	if int(msgapi.Areas[curAreaID].GetCount())-oy < sy-1 && cy != sy {
		v.SetCursor(0, int(msgapi.Areas[curAreaID].GetCount())-oy)
	} else if cy < sy-1 {
		v.SetCursor(0, sy-1)
	} else if int(msgapi.Areas[curAreaID].GetCount())-oy > sy-1 {
		v.SetOrigin(0, oy+sy-1)
		if int(msgapi.Areas[curAreaID].GetCount())-(oy+sy-1) < sy-1 {
			v.SetCursor(0, int(msgapi.Areas[curAreaID].GetCount())-(oy+sy-1))
		}
	}
	return nil
}

func cancelSelectMessage(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("listMsgs")
	ActiveWindow = "MsgBody"
	return nil
}
