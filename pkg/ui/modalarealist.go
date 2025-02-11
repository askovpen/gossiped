package ui

import (
	"strconv"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	defFg, defBg, _ := config.StyleDefault.Decompose()
	m := &ModalAreaList{
		Box:       tview.NewBox().SetBackgroundColor(defBg),
		textColor: defFg,
	}
	borderFg, _, borderAttr := config.GetElementStyle(config.ColorAreaAreaListModal, config.ColorElementBorder).Decompose()
	headerStyle := config.GetElementStyle(config.ColorAreaAreaListModal, config.ColorElementHeader)
	selectionStyle := config.GetElementStyle(config.ColorAreaAreaListModal, config.ColorElementSelection)
	itemStyle := config.GetElementStyle(config.ColorAreaAreaListModal, config.ColorElementItem)
	fgItem, bgItem, attrItem := itemStyle.Decompose()
	fgHeader, bgHeader, attrHeader := headerStyle.Decompose()
	m.table = tview.NewTable().
		SetFixed(1, 0).
		SetBordersColor(borderFg).
		SetSelectable(true, false).
		SetSelectedStyle(selectionStyle).
		SetSelectedFunc(func(row int, column int) {
			m.done(row)
		})
	m.frame = tview.NewFrame(m.table).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBackgroundColor(defBg)
	m.table.SetBackgroundColor(defBg)
	m.frame.SetBorder(true).
		SetTitleAlign(tview.AlignLeft).
		SetBorderAttributes(borderAttr).
		SetBorderColor(borderFg).
		SetBorderPadding(0, 0, 1, 1)
	m.table.SetCell(
		0, 0, tview.NewTableCell(" Area").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false))
	m.table.SetCell(
		0, 1, tview.NewTableCell("EchoID").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetExpansion(1).
			SetSelectable(false))
	m.table.SetCell(
		0, 2, tview.NewTableCell("Msgs").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	m.table.SetCell(
		0, 3, tview.NewTableCell("   New").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	for i, ar := range msgapi.Areas {
		m.table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(i), 10)+" ").
			SetAlign(tview.AlignRight).SetTextColor(fgItem).SetBackgroundColor(bgItem).SetAttributes(attrItem))
		m.table.SetCell(i+1, 1, tview.NewTableCell(ar.GetName()).
			SetTextColor(fgItem).SetBackgroundColor(bgItem).SetAttributes(attrItem))
		m.table.SetCell(i+1, 2, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()), 10)).
			SetAlign(tview.AlignRight).
			SetTextColor(fgItem).SetBackgroundColor(bgItem).SetAttributes(attrItem))
		m.table.SetCell(i+1, 3, tview.NewTableCell(strconv.FormatInt(int64(ar.GetCount()-ar.GetLast()), 10)).
			SetAlign(tview.AlignRight).
			SetTextColor(fgItem).SetBackgroundColor(bgItem).SetAttributes(attrItem))
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
	style := config.GetElementStyle(config.ColorAreaAreaListModal, config.ColorElementTitle)
	m.frame.SetTitle(config.FormatTextWithStyle(text, style))
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

// InputHandler handle input
func (m *ModalAreaList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.HasFocus() {
			if handler := m.table.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		}
	})
}
