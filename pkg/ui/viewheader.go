package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/utils"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	//"github.com/mattn/go-runewidth"
	//"log"
)

// ViewHeader widget
type ViewHeader struct {
	*tview.Box
	sInputs   [8][]rune
	sPosition int
	sCoords   [8]coords
	done      func(string)
	msg       *msgapi.Message
}

// NewViewHeader create new ViewHeader
func NewViewHeader(msg *msgapi.Message) *ViewHeader {
	var si [8][]rune
	if msg == nil {
		si = [8][]rune{
			[]rune("0"),
			[]rune("0"),
			[]rune(""),
			[]rune(""),
			[]rune(""),
			[]rune(""),
			[]rune(""),
			[]rune(""),
		}
	} else {
		repl := ""
		if msg.ReplyTo > 0 {
			repl = fmt.Sprintf("-%d ", msg.ReplyTo)
		}
		for _, rn := range msg.Replies {
			repl += fmt.Sprintf("+%d ", rn)
		}
		si = [8][]rune{
			[]rune(fmt.Sprintf("%d", msg.MsgNum)),
			[]rune(fmt.Sprintf("%d", msgapi.Areas[msgapi.Lookup(msg.Area)].GetCount())),
			[]rune(repl),
			[]rune(msg.From),
			[]rune(msg.FromAddr.String()),
			[]rune(msg.To),
			[]rune(msg.ToAddr.String()),
			[]rune(msg.Subject),
		}
	}
	eh := &ViewHeader{
		Box: tview.NewBox().SetBackgroundColor(tcell.ColorDefault),
		sCoords: [8]coords{
			{f: 8, t: 13, y: 0},
			{f: 17, t: 22, y: 0},
			{f: 23, t: 67, y: 0},
			{f: 8, t: 42, y: 1},
			{f: 43, t: 58, y: 1},
			{f: 8, t: 42, y: 2},
			{f: 43, t: 58, y: 2},
			{f: 8, t: 67, y: 3},
		},
		sInputs:   si,
		sPosition: 0,
		msg:       msg,
	}
	return eh
}

// Draw header
func (e *ViewHeader) Draw(screen tcell.Screen) {
	e.Box.Draw(screen)
	x, y, _, _ := e.GetInnerRect()
	tview.Print(screen, "of", x+14, y, 2, 0, tcell.ColorSilver)
	tview.Print(screen, "Msg  :", x+1, y, 6, 0, tcell.ColorSilver)
	tview.Print(screen, "From :", x+1, y+1, 6, 0, tcell.ColorSilver)
	tview.Print(screen, "To   :", x+1, y+2, 6, 0, tcell.ColorSilver)
	tview.Print(screen, "Subj :", x+1, y+3, 6, 0, tcell.ColorSilver)
	if e.HasFocus() {
		for i := e.sCoords[0].f; i < e.sCoords[0].t; i++ {
			screen.SetContent(x+i, y+e.sCoords[0].y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorNavy))
		}
	}
	for i := 0; i < 8; i++ {
		str := string(e.sInputs[i])
		if utils.NamesEqual(config.Config.Username, str) {
			str = "[::b]" + str
		}
		tview.Print(screen, str, x+e.sCoords[i].f, y+e.sCoords[i].y, len(e.sInputs[i]), 0, tcell.ColorSilver)
	}
	if e.HasFocus() {
		screen.ShowCursor(x+e.sCoords[0].f+len(e.sInputs[0][:e.sPosition]), y+e.sCoords[0].y)
	}
}

// InputHandler event handler
func (e *ViewHeader) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return e.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		add := func(r rune) {
			e.sInputs[0] = append(e.sInputs[0], ' ')
			copy(e.sInputs[0][e.sPosition+1:], e.sInputs[0][e.sPosition:])
			e.sInputs[0][e.sPosition] = r
			e.sPosition++
		}
		switch key := event.Key(); key {
		case tcell.KeyRight:
			if e.sPosition < len(e.sInputs[0]) {
				e.sPosition++
			}
		case tcell.KeyLeft:
			if e.sPosition > 0 {
				e.sPosition--
			}
		case tcell.KeyEnter:
			if e.done != nil {
				if len(e.sInputs[0]) > 0 {
					e.done(string(e.sInputs[0]))
				}
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if e.sPosition > 0 {
				if e.sPosition < len(e.sInputs[0]) {
					e.sInputs[0] = append(e.sInputs[0][:(e.sPosition-1)], e.sInputs[0][e.sPosition:]...)
				} else {
					e.sInputs[0] = e.sInputs[0][:(e.sPosition - 1)]
				}
				e.sPosition--
			}
		case tcell.KeyRune:
			if event.Rune() >= '0' && event.Rune() <= '9' && len(e.sInputs[0]) < (e.sCoords[0].t-e.sCoords[0].f) {
				add(event.Rune())
			}
		}
	})
}

// SetDoneFunc callback
func (e *ViewHeader) SetDoneFunc(handler func(string)) *ViewHeader {
	e.done = handler
	return e
}
