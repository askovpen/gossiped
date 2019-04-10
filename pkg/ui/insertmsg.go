package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/types"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func (a *App) InsertMsg(areaId int) (string, tview.Primitive, bool, bool) {
	newMsg := &msgapi.Message{From: config.Config.Username, FromAddr: config.Config.Address, AreaID: areaId}
	eh := NewEditHeader(config.Config.Username, config.Config.Address.String(), "", "", "")
	eh.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorBlue).
		SetTitle(" " + msgapi.Areas[areaId].GetName() + " ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorYellow)
	eb := NewEditBody()
	eh.SetDoneFunc(func(r [5][]rune) {
		newMsg.From = string(r[0])
		newMsg.FromAddr = types.AddrFromString(string(r[1]))
		newMsg.To = string(r[2])
		newMsg.ToAddr = types.AddrFromString(string(r[3]))
		newMsg.Subject = string(r[4])
		mv, p := newMsg.ToEditNewView()
		eb.SetText(mv, p)
		a.App.SetFocus(eb)
	})
	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(eh, 6, 1, true).
		AddItem(eb, 0, 1, false)
	return fmt.Sprintf("InsertMsg-%s", msgapi.Areas[areaId].GetName()), layout, true, true
}
