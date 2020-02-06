package editor

import (
	"strings"
	//"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// The View struct stores information about a view into a buffer.
// It stores information about the cursor, and the viewport
// that the user sees the buffer from.
type View struct {
	*tview.Box

	// A pointer to the buffer's cursor for ease of access
	Cursor *Cursor

	// The topmost line, used for vertical scrolling
	Topline int
	// The leftmost column, used for horizontal scrolling
	leftCol int

	// Specifies whether or not this view is readonly
	Readonly bool

	// Actual width and height
	width  int
	height int

	// Where this view is located
	x, y int

	// How much to offset because of line numbers
	lineNumOffset int

	// The buffer
	Buf *Buffer

	// We need to keep track of insert key press toggle
	isOverwriteMode bool
	// lastLoc         Loc

	// lastCutTime stores when the last ctrl+k was issued.
	// It is used for clearing the clipboard to replace it with fresh cut lines.
	// lastCutTime time.Time

	// The cellview used for displaying and syntax highlighting
	cellview *CellView

	// The scrollbar
	scrollbar *ScrollBar

	// The keybindings
	bindings KeyBindings

	// The colorscheme
	colorscheme Colorscheme

	// The runtime files
	done func()
}

// NewView returns a new view with the specified buffer.
func NewView(buf *Buffer) *View {
	v := new(View)

	v.Box = tview.NewBox()
	v.x, v.y, v.width, v.height = 0, 0, 0, 0

	v.cellview = new(CellView)

	v.OpenBuffer(buf)

	v.scrollbar = &ScrollBar{
		view: v,
	}
	v.bindings = DefaultKeyBindings

	return v
}

// SetRect sets a new position for the view.
func (v *View) SetRect(x, y, width, height int) {
	v.Box.SetRect(x, y, width, height)
	v.x, v.y, v.width, v.height = v.Box.GetInnerRect()
}

// InputHandler returns a handler which received key events when this view has focus,
func (v *View) InputHandler() func(event *tcell.EventKey, _ func(p tview.Primitive)) {
	return v.WrapInputHandler(func(event *tcell.EventKey, _ func(p tview.Primitive)) {
		v.HandleEvent(event)
	})
}

// GetKeybindings gets the keybindings for this view.
func (v *View) GetKeybindings() KeyBindings {
	return v.bindings
}

// SetKeybindings sets the keybindings for this view.
func (v *View) SetKeybindings(bindings KeyBindings) {
	v.bindings = bindings
}

// SetColorscheme sets the colorscheme for this view.
func (v *View) SetColorscheme(colorscheme Colorscheme) {
	v.colorscheme = colorscheme
	v.Buf.updateRules()
}

// ScrollUp scrolls the view up n lines (if possible)
func (v *View) ScrollUp(n int) {
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline-n >= 0 {
		v.Topline -= n
	} else if v.Topline > 0 {
		v.Topline--
	}
}

// ScrollDown scrolls the view down n lines (if possible)
func (v *View) ScrollDown(n int) {
	// Try to scroll by n but if it would overflow, scroll by 1
	if v.Topline+n <= v.Buf.NumLines {
		v.Topline += n
	} else if v.Topline < v.Buf.NumLines-1 {
		v.Topline++
	}
}

// OpenBuffer opens a new buffer in this view.
// This resets the topline, event handler and cursor.
func (v *View) OpenBuffer(buf *Buffer) {
	v.Buf = buf
	v.Cursor = &buf.Cursor
	v.Topline = 0
	v.leftCol = 0
	v.Cursor.ResetSelection()
	v.Relocate()
	v.Center()

	// Set isOverwriteMode to false, because we assume we are in the default mode when editor
	// is opened
	v.isOverwriteMode = false
	v.Buf.updateRules()
	v.SetColorscheme(ParseColorscheme(`
	color-link comment "bold yellow"
	color-link icomment "bold white"
	color-link origin "bold white"
	color-link tearline "bold white"
	color-link tagline "bold white"
	color-link kludge "bold black"
	`))
}

// Bottomline returns the line number of the lowest line in the view
// You might think that this is obviously just v.Topline + v.height
// but if softwrap is enabled things get complicated since one buffer
// line can take up multiple lines in the view
func (v *View) Bottomline() int {
	screenX, screenY := 0, 0
	numLines := 0
	for lineN := v.Topline; lineN < v.Topline+v.height; lineN++ {
		line := v.Buf.Line(lineN)

		colN := 0
		for _, ch := range line {
			if screenX >= v.width-v.lineNumOffset {
				screenX = 0
				screenY++
			}

			if ch == '\t' {
				screenX += int(v.Buf.Settings["tabsize"].(float64)) - 1
			}

			screenX++
			colN++
		}
		screenX = 0
		screenY++
		numLines++

		if screenY >= v.height {
			break
		}
	}
	return numLines + v.Topline
}

// Relocate moves the view window so that the cursor is in view
// This is useful if the user has scrolled far away, and then starts typing
func (v *View) Relocate() bool {
	height := v.Bottomline() - v.Topline
	ret := false
	cy := v.Cursor.Y
	scrollmargin := int(v.Buf.Settings["scrollmargin"].(float64))
	if cy < v.Topline+scrollmargin && cy > scrollmargin-1 {
		v.Topline = cy - scrollmargin
		ret = true
	} else if cy < v.Topline {
		v.Topline = cy
		ret = true
	}
	if cy > v.Topline+height-1-scrollmargin && cy < v.Buf.NumLines-scrollmargin {
		v.Topline = cy - height + 1 + scrollmargin
		ret = true
	} else if cy >= v.Buf.NumLines-scrollmargin && cy >= height {
		v.Topline = v.Buf.NumLines - height
		ret = true
	}

	return ret
}

// ExecuteActions executes the supplied actions
func (v *View) ExecuteActions(actions []func(*View) bool) bool {
	relocate := false
	readonlyBindingsList := []string{"Delete", "Insert", "Backspace", "Cut", "Play", "Paste", "Move", "Add", "DuplicateLine", "Macro"}
	for _, action := range actions {
		readonlyBindingsResult := false
		funcName := ShortFuncName(action)
		if v.Readonly {
			// check for readonly and if true only let key bindings get called if they do not change the contents.
			for _, readonlyBindings := range readonlyBindingsList {
				if strings.Contains(funcName, readonlyBindings) {
					readonlyBindingsResult = true
				}
			}
		}
		if !readonlyBindingsResult {
			// call the key binding
			relocate = action(v) || relocate
		}
	}

	return relocate
}

// SetCursor sets the view's and buffer's cursor
func (v *View) SetCursor(c *Cursor) bool {
	if c == nil {
		return false
	}
	v.Cursor = c
	v.Buf.curCursor = c.Num

	return true
}

// HandleEvent handles an event passed by the main loop
func (v *View) HandleEvent(event tcell.Event) {
	if !v.HasFocus() {return}
	// This bool determines whether the view is relocated at the end of the function
	// By default it's true because most events should cause a relocate
	relocate := true

	switch e := event.(type) {
	case *tcell.EventKey:
		// Check first if input is a key binding, if it is we 'eat' the input and don't insert a rune
		isBinding := false
		for key, actions := range v.bindings {
			if e.Key() == key.keyCode {
				if e.Key() == tcell.KeyRune {
					if e.Rune() != key.r {
						continue
					}
				}
				if e.Modifiers() == key.modifiers {
					for _, c := range v.Buf.cursors {
						ok := v.SetCursor(c)
						if !ok {
							break
						}
						relocate = false
						isBinding = true
						relocate = v.ExecuteActions(actions) || relocate
					}
					v.SetCursor(&v.Buf.Cursor)
					v.Buf.MergeCursors()
					break
				}
			}
		}

		if !isBinding && e.Key() == tcell.KeyRune {
			// Check viewtype if readonly don't insert a rune (readonly help and log view etc.)
			if !v.Readonly {
				for _, c := range v.Buf.cursors {
					v.SetCursor(c)

					// Insert a character
					if v.Cursor.HasSelection() {
						v.Cursor.DeleteSelection()
						v.Cursor.ResetSelection()
					}

					if v.isOverwriteMode {
						next := v.Cursor.Loc
						next.X++
						v.Buf.Replace(v.Cursor.Loc, next, string(e.Rune()))
					} else {
						v.Buf.Insert(v.Cursor.Loc, string(e.Rune()))
					}
				}
				v.SetCursor(&v.Buf.Cursor)
			}
		}
	}

	if relocate {
		v.Relocate()
		// We run relocate again because there's a bug with relocating with softwrap
		// when for example you jump to the bottom of the buffer and it tries to
		// calculate where to put the topline so that the bottom line is at the bottom
		// of the terminal and it runs into problems with visual lines vs real lines.
		// This is (hopefully) a temporary solution
		v.Relocate()
	}
}

func (v *View) mainCursor() bool {
	return v.Buf.curCursor == len(v.Buf.cursors)-1
}

// displayView draws the view to the screen
func (v *View) displayView(screen tcell.Screen) {
	if v.leftCol != 0 {
		v.leftCol = 0
	}

	v.lineNumOffset = 0

	xOffset := v.x + v.lineNumOffset
	yOffset := v.y

	height := v.height
	width := v.width
	left := v.leftCol
	top := v.Topline

	v.cellview.Draw(v.Buf, v.colorscheme, top, height, left, width-v.lineNumOffset)

	realLineN := top - 1
	visualLineN := 0
	var line []*Char
	for visualLineN, line = range v.cellview.lines {
		var firstChar *Char
		if len(line) > 0 {
			firstChar = line[0]
		}

		if firstChar != nil {
			realLineN = firstChar.realLoc.Y
		} else {
			realLineN++
		}

		var lastChar *Char
		cursorSet := false
		for _, char := range line {
			if char != nil {
				lineStyle := char.style

				for _, c := range v.Buf.cursors {
					v.SetCursor(c)
				}
				v.SetCursor(&v.Buf.Cursor)

				screen.SetContent(xOffset+char.visualLoc.X, yOffset+char.visualLoc.Y, char.drawChar, nil, lineStyle)

				for i, c := range v.Buf.cursors {
					v.SetCursor(c)
					if !v.Cursor.HasSelection() &&
						v.Cursor.Y == char.realLoc.Y && v.Cursor.X == char.realLoc.X && (!cursorSet || i != 0) {
						ShowMultiCursor(screen, xOffset+char.visualLoc.X, yOffset+char.visualLoc.Y, i)
						cursorSet = true
					}
				}
				v.SetCursor(&v.Buf.Cursor)

				lastChar = char
			}
		}

		lastX := 0
		if lastChar != nil {
			lastX = xOffset + lastChar.visualLoc.X + lastChar.width
			for i, c := range v.Buf.cursors {
				v.SetCursor(c)
				if !v.Cursor.HasSelection() &&
					v.Cursor.Y == lastChar.realLoc.Y && v.Cursor.X == lastChar.realLoc.X+1 {
					ShowMultiCursor(screen, lastX, yOffset+lastChar.visualLoc.Y, i)
				}
			}
			v.SetCursor(&v.Buf.Cursor)
		} else if len(line) == 0 {
			for i, c := range v.Buf.cursors {
				v.SetCursor(c)
				if !v.Cursor.HasSelection() &&
					v.Cursor.Y == realLineN {
					ShowMultiCursor(screen, xOffset, yOffset+visualLineN, i)
				}
			}
			v.SetCursor(&v.Buf.Cursor)
			lastX = xOffset
		}
	}
}

// ShowMultiCursor will display a cursor at a location
// If i == 0 then the terminal cursor will be used
// Otherwise a fake cursor will be drawn at the position
func ShowMultiCursor(screen tcell.Screen, x, y, i int) {
	if i == 0 {
		screen.ShowCursor(x, y)
	} else {
		r, _, _, _ := screen.GetContent(x, y)
		screen.SetContent(x, y, r, nil, defStyle.Reverse(true))
	}
}

// Draw renders the view and the cursor
func (v *View) Draw(screen tcell.Screen) {
	v.Box.Draw(screen)

	v.x, v.y, v.width, v.height = v.Box.GetInnerRect()

	// TODO(pdg): just clear from the last line down.
	for y := v.y; y < v.y+v.height; y++ {
		for x := v.x; x < v.x+v.width; x++ {
			screen.SetContent(x, y, ' ', nil, defStyle)
		}
	}

	v.displayView(screen)

	// Don't draw the cursor if it is out of the viewport or if it has a selection
	if v.Cursor.Y-v.Topline < 0 || v.Cursor.Y-v.Topline > v.height-1 || v.Cursor.HasSelection() || v.Readonly {
		screen.HideCursor()
	}

	if v.Buf.Settings["scrollbar"].(bool) {
		v.scrollbar.Display(screen)
	}
}

// SetDoneFunc callback
func (v *View) SetDoneFunc(handler func()) *View {
	v.done = handler
	return v
}
