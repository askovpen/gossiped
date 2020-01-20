package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// SearchString struct
type SearchString struct {
	*tview.Box
	txt string
}

func NewSearchString() *SearchString {
	return &SearchString{
		Box: tview.NewBox().SetBackgroundColor(tcell.ColorDefault),
	}
}

// Draw searchString
func (e *SearchString) Draw(screen tcell.Screen) {
	e.Box.Draw(screen)
	x, y, _, _ := e.GetInnerRect()
	tview.Print(screen, ">>Pick New Area: ", x, y, 17, 0, tcell.ColorSilver)
	tview.Print(screen, e.txt, x+17, y, len(e.txt), 0, tcell.ColorSilver)
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
