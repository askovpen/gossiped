package ui

import (
	"fmt"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ModalHelp widget
type ModalHelp struct {
	*tview.Box
	txt       *tview.TextView
	frame     *tview.Frame
	textColor tcell.Color
	title     string
	done      func()
}

// NewModalHelp return new ModalHelp
func NewModalHelp() *ModalHelp {
	m := &ModalHelp{
		Box:       tview.NewBox(),
		textColor: tview.Styles.PrimaryTextColor,
	}
	m.txt = tview.NewTextView()
	m.frame = tview.NewFrame(m.txt).SetBorders(0, 0, 1, 0, 0, 0)
	fgBorder, bgBorder, styleBorder := config.GetElementStyle(config.ColorAreaHelp, "border").Decompose()
	fgTitle, bgTitle, styleTitle := config.GetElementStyle(config.ColorAreaHelp, "title").Decompose()
	title := fmt.Sprintf("[:%s:%s] Keys ", bgTitle.String(), config.MaskToStringStyle(styleTitle))
	m.frame.SetBorder(true).
		SetBackgroundColor(bgBorder).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(fgBorder).
		SetBorderAttributes(styleBorder).
		SetTitleColor(fgTitle).
		SetTitleAlign(tview.AlignLeft).
		SetTitle(title)
	return m
}

// InputHandler Input Handler
func (m *ModalHelp) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if event.Key() == tcell.KeyEscape {
			m.done()
		}
	})
}

// SetText Set Text
func (m *ModalHelp) SetText(txt string) *ModalHelp {
	style := config.GetElementStyle(config.ColorAreaHelp, "text")
	m.txt.SetTextStyle(style)
	m.txt.SetText(txt)
	return m
}

// Draw draw
func (m *ModalHelp) Draw(screen tcell.Screen) {
	width, height := screen.Size()
	height--
	m.frame.Clear()
	x := 0
	y := 0
	m.SetRect(x, y, width, height)

	// Draw the frame.
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}

// SetDoneFunc Set Done Function
func (m *ModalHelp) SetDoneFunc(handler func()) *ModalHelp {
	m.done = handler
	return m
}

// AreaListHelp Area List Help
func (a *App) AreaListHelp() (string, tview.Primitive, bool, bool) {
	modal := NewModalHelp().
		SetText(`
Home         Move selection bar to first area
End          Move selection bar to last area
Down         Move selection bar to next area
Up           Move selection bar to previous area
Enter, Right Enter the Reader for the selected area
ESC          Exit gossipEd, prompt for final decision
Ctrl-C       Exit immediately, no questions asked
<xyz>        Search for areas containing the string xyz`).
		SetDoneFunc(func() {
			a.Pages.HidePage("AreaListHelp")
		})
	return "AreaListHelp", modal, false, false
}

// ViewMsgHelp View Msg Help
func (a *App) ViewMsgHelp() (string, tview.Primitive, bool, bool) {
	modal := NewModalHelp().
		SetText(`
Ins, Ctrl-I    Enter a new message
Del            Delete current/marked message(s), ask first
Right/Left     Next/Previous message
Home/End       Display first/last part of current message
</>            Go to First/Last message
Ctrl-G         Go to message number
F3, Ctrl-Q     Quote-Reply to message. (Reply to FROM name)
Ctrl-N         Quote-Reply in another area
Ctrl-L         Enter the Message Lister
Ctrl-F         Forward message to another area
Alt-K          Show Kludges
`).
		SetDoneFunc(func() {
			a.Pages.HidePage("ViewMsgHelp")
			a.Pages.RemovePage("ViewMsgHelp")
		})
	return "ViewMsgHelp", modal, true, true
}
