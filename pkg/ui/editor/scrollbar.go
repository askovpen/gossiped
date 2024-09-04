package editor

import (
	"github.com/askovpen/gossiped/pkg/config"
	"github.com/gdamore/tcell/v2"
)

// ScrollBar represents an optional scrollbar that can be used
type ScrollBar struct {
	view *View
}

// Display shows the scrollbar
func (sb *ScrollBar) Display(screen tcell.Screen) {
	style := config.StyleDefault.Reverse(true)
	screen.SetContent(sb.view.x+sb.view.width-1, sb.view.y+sb.pos(), ' ', nil, style)
}

func (sb *ScrollBar) pos() int {
	numlines := sb.view.Buf.NumLines
	h := sb.view.height
	filepercent := float32(sb.view.Topline) / float32(numlines)

	return int(filepercent * float32(h))
}
