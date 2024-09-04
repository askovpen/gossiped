package ui

import (
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	//"github.com/mattn/go-runewidth"
	//"log"
)

type coords struct {
	f int
	t int
	y int
}

// EditHeader widget
type EditHeader struct {
	*tview.Box
	sIndex    int
	sInputs   [5][]rune
	sPosition [5]int
	sCoords   [5]coords
	done      func([5][]rune)
	msg       *msgapi.Message
}

// NewEditHeader create new EditHeader
func NewEditHeader(msg *msgapi.Message) *EditHeader {
	eh := &EditHeader{
		Box: tview.NewBox().SetBackgroundColor(tcell.ColorDefault),
		sCoords: [5]coords{
			{f: 8, t: 42, y: 1},
			{f: 43, t: 58, y: 1},
			{f: 8, t: 42, y: 2},
			{f: 43, t: 58, y: 2},
			{f: 8, t: 67, y: 3},
		},
		sInputs: [5][]rune{
			[]rune(msg.From),
			[]rune(msg.FromAddr.String()),
			[]rune(msg.To),
			[]rune(msg.ToAddr.String()),
			[]rune(msg.Subject),
		},
		sPosition: [5]int{stringWidth(msg.From), stringWidth(msg.FromAddr.String()), stringWidth(msg.To), stringWidth(msg.ToAddr.String()), stringWidth(msg.Subject)},
		sIndex:    0,
		msg:       msg,
	}
	return eh
}

// Draw header
func (e *EditHeader) Draw(screen tcell.Screen) {
	e.Box.Draw(screen)

	boxFg, boxBg, _ := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementWindow).Decompose()
	//_, defBg, _ := config.StyleDefault.Decompose()
	e.Box.SetBackgroundColor(boxBg)
	x, y, _, _ := e.GetInnerRect()
	itemStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementItem)
	itemStyle = itemStyle.Attributes(tcell.AttrNone)
	headerStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementHeader)
	selectionStyle := config.GetElementStyle(config.ColorAreaMessageHeader, config.ColorElementSelection)

	//defaultStyle := config.StyleDefault
	//defFg, _, _ := defaultStyle.Decompose()
	//_, bgSel, _ := styles.GetElementStyleDecomposed(styles.ColorAreaMessageHeader, styles.ColorElementSelection)

	tview.Print(screen, config.FormatTextWithStyle("Msg  :", headerStyle), x+1, y, 6, 0, boxBg)
	tview.Print(screen, config.FormatTextWithStyle("From :", headerStyle), x+1, y+1, 6, 0, boxBg)
	tview.Print(screen, config.FormatTextWithStyle("To   :", headerStyle), x+1, y+2, 6, 0, boxBg)
	tview.Print(screen, config.FormatTextWithStyle("Subj :", headerStyle), x+1, y+3, 6, 0, boxBg)

	if e.HasFocus() {
		for i := e.sCoords[e.sIndex].f; i < e.sCoords[e.sIndex].t; i++ {
			screen.SetContent(x+i, y+e.sCoords[e.sIndex].y, ' ', nil, selectionStyle)
		}
	}
	for i := 0; i < 5; i++ {
		tview.Print(screen, config.FormatTextWithStyle(string(e.sInputs[i]), itemStyle), x+e.sCoords[i].f, y+e.sCoords[i].y, len(e.sInputs[i]), 0, boxFg)
	}
	if e.HasFocus() {
		screen.ShowCursor(x+e.sCoords[e.sIndex].f+len(e.sInputs[e.sIndex][:e.sPosition[e.sIndex]]), y+e.sCoords[e.sIndex].y)
	}
}

// InputHandler event handler
func (e *EditHeader) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return e.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		add := func(r rune) {
			e.sInputs[e.sIndex] = append(e.sInputs[e.sIndex], ' ')
			copy(e.sInputs[e.sIndex][e.sPosition[e.sIndex]+1:], e.sInputs[e.sIndex][e.sPosition[e.sIndex]:])
			e.sInputs[e.sIndex][e.sPosition[e.sIndex]] = r
			e.sPosition[e.sIndex]++
		}
		switch key := event.Key(); key {
		case tcell.KeyTab:
			e.sIndex++
			if e.sIndex == 5 {
				e.sIndex = 0
			} else if (*e.msg.AreaObject).GetType() != msgapi.EchoAreaTypeNetmail && e.sIndex == 3 {
				e.sIndex = 4
			}
		case tcell.KeyRight:
			if e.sPosition[e.sIndex] < len(e.sInputs[e.sIndex]) {
				e.sPosition[e.sIndex]++
			}
		case tcell.KeyLeft:
			if e.sPosition[e.sIndex] > 0 {
				e.sPosition[e.sIndex]--
			}
		case tcell.KeyEnter:
			if e.sIndex == 4 {
				if e.done != nil {
					if len(e.sInputs[0]) > 0 && len(e.sInputs[1]) > 0 && len(e.sInputs[2]) > 0 {
						e.done(e.sInputs)
					}
				}
			} else {
				e.sIndex++
				if (*e.msg.AreaObject).GetType() != msgapi.EchoAreaTypeNetmail && e.sIndex == 3 {
					e.sIndex = 4
				}
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if e.sPosition[e.sIndex] > 0 {
				if e.sPosition[e.sIndex] < len(e.sInputs[e.sIndex]) {
					e.sInputs[e.sIndex] = append(e.sInputs[e.sIndex][:(e.sPosition[e.sIndex]-1)], e.sInputs[e.sIndex][e.sPosition[e.sIndex]:]...)
				} else {
					e.sInputs[e.sIndex] = e.sInputs[e.sIndex][:(e.sPosition[e.sIndex] - 1)]
				}
				e.sPosition[e.sIndex]--
			}
		case tcell.KeyRune:
			add(event.Rune())
		}
	})
}

// SetDoneFunc callback
func (e *EditHeader) SetDoneFunc(handler func([5][]rune)) *EditHeader {
	e.done = handler
	return e
}
