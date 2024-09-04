package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/askovpen/gossiped/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	//"github.com/mattn/go-runewidth"
	"strings"
)

// ViewHeader widget
type ViewHeader struct {
	*tview.Box
	sInputs   [10][]rune
	sPosition int
	sCoords   [10]coords
	done      func(string)
	msg       *msgapi.Message
}

// NewViewHeader create new ViewHeader
func NewViewHeader(msg *msgapi.Message) *ViewHeader {
	var si [10][]rune
	if msg == nil {
		si = [10][]rune{[]rune("0"), []rune("0"), []rune(""), []rune(""), []rune(""), []rune(""), []rune(""), []rune(""), []rune(""), []rune("")}
	} else {
		repl := ""
		if msg.ReplyTo > 0 {
			repl = fmt.Sprintf("-%d ", msg.ReplyTo)
		}
		for _, rn := range msg.Replies {
			repl += fmt.Sprintf("+%d ", rn)
		}
		if msg.Corrupted {
			msg.Attrs = append(msg.Attrs, "[red]Corrupted")
		}
		if len(msg.Attrs) > 0 {
			repl += "[" + strings.Join(msg.Attrs, " ") + "]"
		}
		si = [10][]rune{
			[]rune(fmt.Sprintf("%d", msg.MsgNum)),
			[]rune(fmt.Sprintf("%d", msgapi.Areas[msgapi.Lookup(msg.Area)].GetCount())),
			[]rune(repl),
			[]rune(msg.From),
			[]rune(msg.FromAddr.String()),
			[]rune(msg.DateWritten.Format("02 Jan 06 15:04:05")),
			[]rune(msg.To),
			[]rune(msg.ToAddr.String()),
			[]rune(msg.DateArrived.Format("02 Jan 06 15:04:05")),
			[]rune(msg.Subject),
		}
	}
	eh := &ViewHeader{
		Box: tview.NewBox().SetBackgroundColor(tcell.ColorDefault),
		sCoords: [10]coords{
			{f: 8, t: 13, y: 0},
			{f: 17, t: 22, y: 0},
			{f: 23, t: 67, y: 0},
			{f: 8, t: 42, y: 1},
			{f: 43, t: 58, y: 1},
			{f: 60, t: 78, y: 1},
			{f: 8, t: 42, y: 2},
			{f: 43, t: 58, y: 2},
			{f: 60, t: 78, y: 2},
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
	defStyle := config.StyleDefault
	boxFg, boxBg, _ := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementWindow).Decompose()
	borderStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementBorder)
	e.Box.SetBackgroundColor(boxBg)
	e.Box.SetBorderStyle(borderStyle)
	x, y, _, _ := e.GetInnerRect()
	itemStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementItem)
	highlightStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementHighlight)
	headerStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementHeader)
	_, bgSel, _ := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementSelection).Decompose()
	tview.Print(screen, config.FormatTextWithStyle("of", itemStyle), x+14, y, 2, 0, boxFg)
	tview.Print(screen, config.FormatTextWithStyle("Msg  :", headerStyle), x+1, y, 6, 0, boxFg)
	tview.Print(screen, config.FormatTextWithStyle("From :", headerStyle), x+1, y+1, 6, 0, boxFg)
	tview.Print(screen, config.FormatTextWithStyle("To   :", headerStyle), x+1, y+2, 6, 0, boxFg)
	tview.Print(screen, config.FormatTextWithStyle("Subj :", headerStyle), x+1, y+3, 6, 0, boxFg)
	if e.HasFocus() {
		for i := e.sCoords[0].f; i < e.sCoords[0].t; i++ {
			screen.SetContent(x+i, y+e.sCoords[0].y, ' ', nil, defStyle.Background(bgSel))
		}
	}
	for i := 0; i < len(e.sCoords); i++ {
		str := string(e.sInputs[i])
		style := itemStyle
		if utils.NamesEqual(config.Config.Username, str) {
			style = highlightStyle
		} else {
			style = itemStyle
		}
		tview.Print(screen, config.FormatTextWithStyle(str, style), x+e.sCoords[i].f, y+e.sCoords[i].y, len(e.sInputs[i]), 0, boxFg)
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
