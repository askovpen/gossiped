package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Modal is a centered message window used to inform the user or prompt them
// for an immediate decision. It needs to have at least one button (added via
// AddButtons()) or it will never disappear.
//
// See https://github.com/rivo/tview/wiki/Modal for an example.
type ModalMenu struct {
	*tview.Box
	table     *tview.Table
	frame     *tview.Frame
	textColor tcell.Color
	title     string
	done      func(buttonIndex int)
	y         int
	width     int
}

// NewModal returns a new modal message window.
func NewModalMenu() *ModalMenu {
	m := &ModalMenu{
		Box:       tview.NewBox(),
		textColor: tview.Styles.PrimaryTextColor,
		y:         1,
		width:     0,
	}
	m.table = tview.NewTable().
		SetSelectable(true, false).
		SetSelectedStyle(tcell.ColorWhite, tcell.ColorBlue, tcell.AttrBold).
		SetSelectedFunc(func(row int, column int) {
			m.done(row)
		})
	m.frame = tview.NewFrame(m.table).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBorder(true).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorderPadding(0, 0, 1, 1).SetBorderColor(tcell.ColorRed).SetBorderAttributes(tcell.AttrBold).SetTitleColor(tcell.ColorYellow)
	return m
}

// SetTextColor sets the color of the message text.
func (m *ModalMenu) SetTextColor(color tcell.Color) *ModalMenu {
	m.textColor = color
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an emptry string.
func (m *ModalMenu) SetDoneFunc(handler func(buttonIndex int)) *ModalMenu {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks. Note that words are wrapped, too, based on the final size of the
// window.
func (m *ModalMenu) SetText(text string) *ModalMenu {
	m.title = text
	m.frame.SetTitle(text)
	return m
}

// AddButtons adds buttons to the window. There must be at least one button and
// a "done" handler so the window can be closed again.
func (m *ModalMenu) SetY(y int) *ModalMenu {
	m.y = y
	return m
}

func (m *ModalMenu) AddButtons(labels []string) *ModalMenu {
	for index, label := range labels {
		func(i int, l string) {
			//m.list.AddItem(label,"",0,func() {m.done(i,l)})
			m.table.SetCell(i, 0, tview.NewTableCell(label))
			if m.width < len(label) {
				m.width = len(label)
			}
		}(index, label)
	}
	return m
}

// Focus is called when this primitive receives focus.
func (m *ModalMenu) Focus(delegate func(p tview.Primitive)) {
	//delegate(m.form)
	delegate(m.table)
}

// HasFocus returns whether or not this primitive has focus.
func (m *ModalMenu) HasFocus() bool {
	//return m.form.HasFocus()
	return m.table.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *ModalMenu) Draw(screen tcell.Screen) {
	height := m.table.GetRowCount() + 2
	width := m.width + 4
	if len(m.title) > width-2 {
		width = len(m.title) + 2
	}
	m.frame.Clear()
	x := 1
	y := m.y
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}
