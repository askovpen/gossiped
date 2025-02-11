package ui

import (
	"fmt"
	"strconv"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
func NewModalMessageList(area *msgapi.AreaPrimitive) *ModalMessageList {
	_, defBg, _ := config.StyleDefault.Decompose()
	m := &ModalMessageList{
		Box:       tview.NewBox().SetBackgroundColor(defBg),
		textColor: tview.Styles.PrimaryTextColor,
	}
	styleBorder := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementBorder)
	styleSelection := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementSelection)
	fgHeader, bgHeader, attrHeader := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementHeader).Decompose()
	fgItem, bgItem, attrItem := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementItem).Decompose()
	fgTitle, bgTitle, attrTitle := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementTitle).Decompose()
	fgHigh, bgHigh, attrHigh := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementHighlight).Decompose()
	//fgCur, bgCur, attrCur := config.GetElementStyle(config.ColorAreaMessageList, config.ColorElementCurrent).Decompose()
	m.table = tview.NewTable().
		SetFixed(1, 0).
		SetSelectable(true, false).
		SetSelectedStyle(styleSelection).
		SetSelectedFunc(func(row int, column int) {
			m.done(uint32(row))
		})
	m.frame = tview.NewFrame(m.table).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBackgroundColor(defBg)
	m.table.SetBackgroundColor(defBg)
	m.frame.SetTitle(fmt.Sprintf("[%s:%s:%s] List Messages ", fgTitle.String(), bgTitle.String(), config.MaskToStringStyle(attrTitle)))
	m.frame.SetBorder(true).
		SetBorderStyle(styleBorder).
		SetBorderPadding(0, 0, 1, 1).
		SetTitleAlign(tview.AlignLeft)
	m.table.SetCell(0, 0, tview.NewTableCell(" Msg ").
		SetSelectable(false).
		SetAlign(tview.AlignRight).
		SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader))
	m.table.SetCell(
		0, 1, tview.NewTableCell("From").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false))
	m.table.SetCell(
		0, 2, tview.NewTableCell("To").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false))
	m.table.SetCell(
		0, 3, tview.NewTableCell("Subj").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetExpansion(1).
			SetSelectable(false))
	m.table.SetCell(
		0, 4, tview.NewTableCell("Written").
			SetTextColor(fgHeader).SetBackgroundColor(bgHeader).SetAttributes(attrHeader).
			SetSelectable(false).
			SetAlign(tview.AlignRight))
	for i, mh := range *(*area).GetMessages() {
		ch := " "
		fg, bg, attr := fgItem, bgItem, attrItem
		if i == int((*area).GetLast()-1) {
			//fg, bg, attr = fgCur, bgCur, attrCur
			fg, bg, attr = fgHigh, bgHigh, attrHigh
			ch = "*"
		}
		fromCondition := utils.NamesEqual(mh.From, config.Config.Username)
		toCondition := utils.NamesEqual(mh.To, config.Config.Username)
		m.table.SetCell(i+1, 0, tview.NewTableCell(strconv.FormatInt(int64(mh.MsgNum), 10)+ch).
			SetAlign(tview.AlignRight).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		//mh.From, mh.To, mh.Subject, mh.DateWritten.Format("02 Jan 06"))
		if fromCondition {
			m.table.SetCell(i+1, 1, tview.NewTableCell(mh.From).
				SetTextColor(fgHigh).SetBackgroundColor(bgHigh).SetAttributes(attrHigh))
		} else {
			m.table.SetCell(i+1, 1, tview.NewTableCell(mh.From).
				SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		}
		if toCondition {
			m.table.SetCell(i+1, 2, tview.NewTableCell(mh.To).
				SetTextColor(fgHigh).SetBackgroundColor(bgHigh).SetAttributes(attrHigh))
		} else {
			m.table.SetCell(i+1, 2, tview.NewTableCell(mh.To).
				SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		}
		m.table.SetCell(i+1, 3, tview.NewTableCell(mh.Subject).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
		m.table.SetCell(i+1, 4, tview.NewTableCell(mh.DateWritten.Format("02 Jan 06")).
			SetTextColor(fg).SetBackgroundColor(bg).SetAttributes(attr))
	}
	m.table.Select(int((*area).GetLast()), 0)
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

// InputHandler handle input
func (m *ModalMessageList) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.HasFocus() {
			if handler := m.table.InputHandler(); handler != nil {
				handler(event, setFocus)
			}
			return
		}
	})
}
