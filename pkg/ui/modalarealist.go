package ui

import (
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strconv"
)

// ModalAreaList is a centered message window used to inform the user or prompt them
type ModalAreaList struct {
	*tview.Box
	table     *tview.Table
	frame     *tview.Frame
	textColor tcell.Color
	title     string
	done      func(buttonIndex int)
}

// NewModalAreaList returns a new modal message window.
func NewModalAreaList() *ModalAreaList {
	m := &ModalAreaList{
		Box:       tview.NewBox(),
		textColor: tview.Styles.PrimaryTextColor,
	}
	m.table = tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.ColorWhite, tcell.ColorNavy, tcell.AttrBold).
		SetSelectedFunc(func(row int, column int) {
			m.done(row)
		})
	m.frame = tview.NewFrame(m.table).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBorder(true).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorBlue).SetBorderAttributes(tcell.AttrBold).SetTitleColor(tcell.ColorYellow).SetTitleAlign(tview.AlignLeft)
	m.table.SetCell(
		0, 0, tview.NewTableCell(" Area").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false))
	m.table.SetCell(
		0, 1, tview.NewTableCell("EchoID").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1).
			SetSelectable(false))
	m.table.SetCell(
		0, 2, tview.NewTableCell("Msgs").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	m.table.SetCell(
		0, 3, tview.NewTableCell("   New").
			SetTextColor(tcell.ColorYellow).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	for i, ar := range msgapi.Areas {
		m.table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i), 10)+" ").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		m.table.SetCell(i+1, 1, tview.NewTableCell(ar.GetName()).SetTextColor(tcell.ColorSilver))
		m.table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()), 10)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
		m.table.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()-ar.GetLast()), 10)).SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSilver))
	}
	return m
}

// SetTextColor sets the color of the message text.
func (m *ModalAreaList) SetTextColor(color tcell.Color) *ModalAreaList {
	m.textColor = color
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an emptry string.
func (m *ModalAreaList) SetDoneFunc(handler func(buttonIndex int)) *ModalAreaList {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks. Note that words are wrapped, too, based on the final size of the
// window.
func (m *ModalAreaList) SetText(text string) *ModalAreaList {
	m.title = text
	m.frame.SetTitle(text)
	return m
}

// AddButtons adds buttons to the window. There must be at least one button and
// a "done" handler so the window can be closed again.

// Focus is called when this primitive receives focus.
func (m *ModalAreaList) Focus(delegate func(p tview.Primitive)) {
	//delegate(m.form)
	delegate(m.table)
}

// HasFocus returns whether or not this primitive has focus.
func (m *ModalAreaList) HasFocus() bool {
	//return m.form.HasFocus()
	return m.table.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *ModalAreaList) Draw(screen tcell.Screen) {
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
