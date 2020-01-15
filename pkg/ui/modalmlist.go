package ui

import (
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/utils"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
)

// ModalMessageList is a centered message window used to inform the user or prompt them
// for an immediate decision. It needs to have at least one button (added via
// AddButtons()) or it will never disappear.
//
// See https://github.com/rivo/tview/wiki/Modal for an example.
type ModalMessageList struct {
	*tview.Box
	table     *tview.Table
	frame     *tview.Frame
	textColor tcell.Color
	done      func(msgNum uint32)
}

// NewModalMessageList returns a new modal message window.
func NewModalMessageList(areaID int) *ModalMessageList {
	m := &ModalMessageList{
		Box:       tview.NewBox(),
		textColor: tview.Styles.PrimaryTextColor,
	}
	m.table = tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.ColorWhite, tcell.ColorNavy, tcell.AttrBold).
		SetSelectedFunc(func(row int, column int) {
			m.done(uint32(row))
		})
	m.frame = tview.NewFrame(m.table).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetTitle("List Messages")
	m.frame.SetBorder(true).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorRed).SetBorderAttributes(tcell.AttrBold).SetTitleColor(tcell.ColorYellow).SetTitleAlign(tview.AlignLeft)
	m.table.SetCell(
		0, 0, tview.NewTableCell(" Msg ").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	m.table.SetCell(
		0, 1, tview.NewTableCell("From").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false))
	m.table.SetCell(
		0, 2, tview.NewTableCell("To").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false))
	m.table.SetCell(
		0, 3, tview.NewTableCell("Subj").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1).
			SetSelectable(false))
	m.table.SetCell(
		0, 4, tview.NewTableCell("Written").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	for i, mh := range *msgapi.Areas[areaID].GetMessages() {
		ch := " "
		if i == int(msgapi.Areas[areaID].GetLast()-1) {
			ch = "[::b],"
		}
		//mh.From, mh.To, mh.Subject, mh.DateWritten.Format("02 Jan 06"))
		m.table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(mh.MsgNum), 10)+ch).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		if utils.NamesEqual(mh.From, config.Config.Username) {
			m.table.SetCell(i+1, 1, tview.NewTableCell(mh.From).SetTextColor(tcell.ColorSilver).SetAttributes(tcell.AttrBold))
		} else {
			m.table.SetCell(i+1, 1, tview.NewTableCell(mh.From).SetTextColor(tcell.ColorSilver))
		}
		if utils.NamesEqual(mh.To, config.Config.Username) {
			m.table.SetCell(i+1, 2, tview.NewTableCell(mh.To).SetTextColor(tcell.ColorSilver).SetAttributes(tcell.AttrBold))
		} else {
			m.table.SetCell(i+1, 2, tview.NewTableCell(mh.To).SetTextColor(tcell.ColorSilver))
		}
		m.table.SetCell(i+1, 3, tview.NewTableCell(mh.Subject).SetTextColor(tcell.ColorSilver))
		m.table.SetCell(i+1, 4, tview.NewTableCell(mh.DateWritten.Format("02 Jan 06")).SetTextColor(tcell.ColorSilver))
	}
	m.table.Select(int(msgapi.Areas[areaID].GetLast()), 0)
	return m
}

// SetTextColor sets the color of the message text.
func (m *ModalMessageList) SetTextColor(color tcell.Color) *ModalMessageList {
	m.textColor = color
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an emptry string.
func (m *ModalMessageList) SetDoneFunc(handler func(msgNum uint32)) *ModalMessageList {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks. Note that words are wrapped, too, based on the final size of the
// window.
//func (m *ModalMessageList) SetText(text string) *ModalMessageList {
//	m.title = text
//	m.frame.SetTitle(text)
//	return m
//}

// AddButtons adds buttons to the window. There must be at least one button and
// a "done" handler so the window can be closed again.

// Focus is called when this primitive receives focus.
func (m *ModalMessageList) Focus(delegate func(p tview.Primitive)) {
	//delegate(m.form)
	delegate(m.table)
}

// HasFocus returns whether or not this primitive has focus.
func (m *ModalMessageList) HasFocus() bool {
	//return m.form.HasFocus()
	return m.table.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *ModalMessageList) Draw(screen tcell.Screen) {
	width, height := screen.Size()
	height -= 7
	m.frame.Clear()
	x := 0
	y := 6
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
