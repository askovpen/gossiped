package ui

import (
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SearchString struct
type SearchString struct {
	*tview.Box
	txt string
}

// NewSearchString create SearchString
func NewSearchString() *SearchString {
	return &SearchString{
		Box: tview.NewBox().SetBackgroundColor(tcell.ColorDefault),
	}
}

// Draw searchString
func (e *SearchString) Draw(screen tcell.Screen) {
	stylePrompt := config.GetElementStyle(config.ColorAreaAreaList, config.ColorElementPrompt)
	//styleBorder := styles.GetElementStyle(styles.ColorAreaAreaList, styles.ColorElementBorder)
	fg, bg, _ := stylePrompt.Decompose()
	e.Box.Draw(screen)
	e.Box.SetBackgroundColor(bg)
	//e.Box.SetBorderStyle(styleBorder)
	x, y, _, _ := e.GetInnerRect()
	tview.Print(screen, config.FormatTextWithStyle(">> Pick New Area: ", stylePrompt), x, y, 18, 0, fg)
	tview.Print(screen, config.FormatTextWithStyle(e.txt, stylePrompt), x+18, y, len(e.txt), 0, fg)
}

// AddChar to searchString
func (e *SearchString) AddChar(ch rune) {
	e.txt += string(ch)
}

// Clear searchString
func (e *SearchString) Clear() {
	e.txt = ""
}

// GetText return searchString text
func (e *SearchString) GetText() string {
	return e.txt
}
