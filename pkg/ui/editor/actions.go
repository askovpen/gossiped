package editor

import (
	"strings"
	"unicode/utf8"
)

func (v *View) deselect(index int) bool {
	if v.Cursor.HasSelection() {
		v.Cursor.Loc = v.Cursor.CurSelection[index]
		v.Cursor.ResetSelection()
		v.Cursor.StoreVisualX()
		return true
	}
	return false
}

// Center centers the view on the cursor
func (v *View) Center() bool {
	v.Topline = v.Cursor.Y - v.height/2
	if v.Topline+v.height > v.Buf.NumLines {
		v.Topline = v.Buf.NumLines - v.height
	}
	if v.Topline < 0 {
		v.Topline = 0
	}
	return true
}

// CursorUp moves the cursor up
func (v *View) CursorUp() bool {
	if v.Readonly {
		v.ScrollUp(1)
		return false
	}
	v.deselect(0)
	v.Cursor.Up()
	return true
}

// CursorDown moves the cursor down
func (v *View) CursorDown() bool {
	if v.Readonly {
		v.ScrollDown(1)
		return false
	}
	v.deselect(1)
	v.Cursor.Down()
	return true
}

// CursorLeft moves the cursor left
func (v *View) CursorLeft() bool {
	if v.Cursor.HasSelection() {
		v.Cursor.Loc = v.Cursor.CurSelection[0]
		v.Cursor.ResetSelection()
		v.Cursor.StoreVisualX()
	} else {
		tabstospaces := v.Buf.Settings["tabstospaces"].(bool)
		tabmovement := v.Buf.Settings["tabmovement"].(bool)
		if tabstospaces && tabmovement {
			tabsize := int(v.Buf.Settings["tabsize"].(float64))
			line := v.Buf.Line(v.Cursor.Y)
			if v.Cursor.X-tabsize >= 0 && line[v.Cursor.X-tabsize:v.Cursor.X] == Spaces(tabsize) && IsStrWhitespace(line[0:v.Cursor.X-tabsize]) {
				for i := 0; i < tabsize; i++ {
					v.Cursor.Left()
				}
			} else {
				v.Cursor.Left()
			}
		} else {
			v.Cursor.Left()
		}
	}
	return true
}

// CursorRight moves the cursor right
func (v *View) CursorRight() bool {
	if v.Cursor.HasSelection() {
		v.Cursor.Loc = v.Cursor.CurSelection[1]
		v.Cursor.ResetSelection()
		v.Cursor.StoreVisualX()
	} else {
		tabstospaces := v.Buf.Settings["tabstospaces"].(bool)
		tabmovement := v.Buf.Settings["tabmovement"].(bool)
		if tabstospaces && tabmovement {
			tabsize := int(v.Buf.Settings["tabsize"].(float64))
			line := v.Buf.Line(v.Cursor.Y)
			if v.Cursor.X+tabsize < Count(line) && line[v.Cursor.X:v.Cursor.X+tabsize] == Spaces(tabsize) && IsStrWhitespace(line[0:v.Cursor.X]) {
				for i := 0; i < tabsize; i++ {
					v.Cursor.Right()
				}
			} else {
				v.Cursor.Right()
			}
		} else {
			v.Cursor.Right()
		}
	}
	return true
}

// StartOfLine moves the cursor to the start of the line
func (v *View) StartOfLine() bool {
	v.deselect(0)

	if v.Cursor.X != 0 {
		v.Cursor.Start()
	} else {
		v.Cursor.StartOfText()
	}
	return true
}

// EndOfLine moves the cursor to the end of the line
func (v *View) EndOfLine() bool {
	v.deselect(0)
	v.Cursor.End()
	return true
}

// Retab changes all tabs to spaces or all spaces to tabs depending
// on the user's settings
func (v *View) Retab() bool {
	toSpaces := v.Buf.Settings["tabstospaces"].(bool)
	tabsize := int(v.Buf.Settings["tabsize"].(float64))
	dirty := false

	for i := 0; i < v.Buf.NumLines; i++ {
		l := v.Buf.Line(i)

		ws := GetLeadingWhitespace(l)
		if ws != "" {
			if toSpaces {
				ws = strings.Replace(ws, "\t", Spaces(tabsize), -1)
			} else {
				ws = strings.Replace(ws, Spaces(tabsize), "\t", -1)
			}
		}

		l = strings.TrimLeft(l, " \t")
		v.Buf.lines[i].data = []byte(ws + l)
		dirty = true
	}

	v.Buf.IsModified = dirty
	return true
}

// CursorStart moves the cursor to the start of the buffer
func (v *View) CursorStart() bool {
	v.deselect(0)

	v.Cursor.X = 0
	v.Cursor.Y = 0

	return true
}

// CursorEnd moves the cursor to the end of the buffer
func (v *View) CursorEnd() bool {
	v.deselect(0)

	v.Cursor.Loc = v.Buf.End()
	v.Cursor.StoreVisualX()

	return true
}

// InsertSpace inserts a space
func (v *View) InsertSpace() bool {
	if v.Cursor.HasSelection() {
		v.Cursor.DeleteSelection()
		v.Cursor.ResetSelection()
	}
	v.Buf.Insert(v.Cursor.Loc, " ")
	// v.Cursor.Right()
	return true
}

// InsertNewline inserts a newline plus possible some whitespace if autoindent is on
func (v *View) InsertNewline() bool {
	// Insert a newline
	if v.Cursor.HasSelection() {
		v.Cursor.DeleteSelection()
		v.Cursor.ResetSelection()
	}

	ws := GetLeadingWhitespace(v.Buf.Line(v.Cursor.Y))
	cx := v.Cursor.X
	v.Buf.Insert(v.Cursor.Loc, "\n")
	// v.Cursor.Right()

	if v.Buf.Settings["autoindent"].(bool) {
		if cx < len(ws) {
			ws = ws[0:cx]
		}
		v.Buf.Insert(v.Cursor.Loc, ws)
		// for i := 0; i < len(ws); i++ {
		// 	v.Cursor.Right()
		// }

		// Remove the whitespaces if keepautoindent setting is off
		if IsSpacesOrTabs(v.Buf.Line(v.Cursor.Y-1)) && !v.Buf.Settings["keepautoindent"].(bool) {
			line := v.Buf.Line(v.Cursor.Y - 1)
			v.Buf.Remove(Loc{0, v.Cursor.Y - 1}, Loc{Count(line), v.Cursor.Y - 1})
		}
	}
	v.Cursor.LastVisualX = v.Cursor.GetVisualX()

	return true
}

// Backspace deletes the previous character
func (v *View) Backspace() bool {
	// Delete a character
	if v.Cursor.HasSelection() {
		v.Cursor.DeleteSelection()
		v.Cursor.ResetSelection()
	} else if v.Cursor.Loc.GreaterThan(v.Buf.Start()) {
		// We have to do something a bit hacky here because we want to
		// delete the line by first moving left and then deleting backwards
		// but the undo redo would place the cursor in the wrong place
		// So instead we move left, save the position, move back, delete
		// and restore the position

		// If the user is using spaces instead of tabs and they are deleting
		// whitespace at the start of the line, we should delete as if it's a
		// tab (tabSize number of spaces)
		lineStart := sliceEnd(v.Buf.LineBytes(v.Cursor.Y), v.Cursor.X)
		tabSize := int(v.Buf.Settings["tabsize"].(float64))
		if v.Buf.Settings["tabstospaces"].(bool) && IsSpaces(lineStart) && utf8.RuneCount(lineStart) != 0 && utf8.RuneCount(lineStart)%tabSize == 0 {
			loc := v.Cursor.Loc
			v.Buf.Remove(loc.Move(-tabSize, v.Buf), loc)
		} else {
			loc := v.Cursor.Loc
			v.Buf.Remove(loc.Move(-1, v.Buf), loc)
		}
	}
	v.Cursor.LastVisualX = v.Cursor.GetVisualX()

	return true
}

// Delete deletes the next character
func (v *View) Delete() bool {
	if v.Cursor.HasSelection() {
		v.Cursor.DeleteSelection()
		v.Cursor.ResetSelection()
	} else {
		loc := v.Cursor.Loc
		if loc.LessThan(v.Buf.End()) {
			v.Buf.Remove(loc, loc.Move(1, v.Buf))
		}
	}
	return true
}

// IndentSelection indents the current selection
func (v *View) IndentSelection() bool {
	if v.Cursor.HasSelection() {
		start := v.Cursor.CurSelection[0]
		end := v.Cursor.CurSelection[1]
		if end.Y < start.Y {
			start, end = end, start
			v.Cursor.SetSelectionStart(start)
			v.Cursor.SetSelectionEnd(end)
		}

		startY := start.Y
		endY := end.Move(-1, v.Buf).Y
		endX := end.Move(-1, v.Buf).X
		tabsize := len(v.Buf.IndentString())
		for y := startY; y <= endY; y++ {
			v.Buf.Insert(Loc{0, y}, v.Buf.IndentString())
			if y == startY && start.X > 0 {
				v.Cursor.SetSelectionStart(start.Move(tabsize, v.Buf))
			}
			if y == endY {
				v.Cursor.SetSelectionEnd(Loc{endX + tabsize + 1, endY})
			}
		}
		v.Cursor.Relocate()

		return true
	}
	return false
}

// OutdentLine moves the current line back one indentation
func (v *View) OutdentLine() bool {
	if v.Cursor.HasSelection() {
		return false
	}

	for x := 0; x < len(v.Buf.IndentString()); x++ {
		if len(GetLeadingWhitespace(v.Buf.Line(v.Cursor.Y))) == 0 {
			break
		}
		v.Buf.Remove(Loc{0, v.Cursor.Y}, Loc{1, v.Cursor.Y})
	}
	v.Cursor.Relocate()
	return true
}

// OutdentSelection takes the current selection and moves it back one indent level
func (v *View) OutdentSelection() bool {
	if v.Cursor.HasSelection() {
		start := v.Cursor.CurSelection[0]
		end := v.Cursor.CurSelection[1]
		if end.Y < start.Y {
			start, end = end, start
			v.Cursor.SetSelectionStart(start)
			v.Cursor.SetSelectionEnd(end)
		}

		startY := start.Y
		endY := end.Move(-1, v.Buf).Y
		for y := startY; y <= endY; y++ {
			for x := 0; x < len(v.Buf.IndentString()); x++ {
				if len(GetLeadingWhitespace(v.Buf.Line(y))) == 0 {
					break
				}
				v.Buf.Remove(Loc{0, y}, Loc{1, y})
			}
		}
		v.Cursor.Relocate()

		return true
	}
	return false
}

// InsertTab inserts a tab or spaces
func (v *View) InsertTab() bool {
	if v.Cursor.HasSelection() {
		return false
	}

	tabBytes := len(v.Buf.IndentString())
	bytesUntilIndent := tabBytes - (v.Cursor.GetVisualX() % tabBytes)
	v.Buf.Insert(v.Cursor.Loc, v.Buf.IndentString()[:bytesUntilIndent])
	// for i := 0; i < bytesUntilIndent; i++ {
	// 	v.Cursor.Right()
	// }

	return true
}

// DeleteLine deletes the current line
func (v *View) DeleteLine() bool {
	v.Cursor.SelectLine()
	if !v.Cursor.HasSelection() {
		return false
	}
	v.Cursor.DeleteSelection()
	v.Cursor.ResetSelection()
	return true
}

// DeleteToEnd deletes the current line
func (v *View) DeleteToEnd() bool {
	x, y := runeToByteIndex(v.Cursor.Loc.X, v.Buf.LineBytes(v.Cursor.Loc.Y)), v.Cursor.Loc.Y
	v.Buf.DeleteToEnd(Loc{x, y})
	return true
}

// Start moves the viewport to the start of the buffer
func (v *View) Start() bool {
	if v.mainCursor() {
		v.Topline = 0
	}
	return false
}

// End moves the viewport to the end of the buffer
func (v *View) End() bool {
	if v.mainCursor() {
		if v.height > v.Buf.NumLines {
			v.Topline = 0
		} else {
			v.Topline = v.Buf.NumLines - v.height
		}

	}
	return false
}

// PageUp scrolls the view up a page
func (v *View) PageUp() bool {
	if v.mainCursor() {
		if v.Topline > v.height {
			v.ScrollUp(v.height)
		} else {
			v.Topline = 0
		}
	}
	return false
}

// PageDown scrolls the view down a page
func (v *View) PageDown() bool {
	if v.mainCursor() {
		if v.Buf.NumLines-(v.Topline+v.height) > v.height {
			v.ScrollDown(v.height)
		} else if v.Buf.NumLines >= v.height {
			v.Topline = v.Buf.NumLines - v.height
		}
	}
	return false
}

// CursorPageUp places the cursor a page up
func (v *View) CursorPageUp() bool {
	v.deselect(0)

	if v.Cursor.HasSelection() {
		v.Cursor.Loc = v.Cursor.CurSelection[0]
		v.Cursor.ResetSelection()
		v.Cursor.StoreVisualX()
	}
	v.Cursor.UpN(v.height)

	return true
}

// CursorPageDown places the cursor a page up
func (v *View) CursorPageDown() bool {
	v.deselect(0)

	if v.Cursor.HasSelection() {
		v.Cursor.Loc = v.Cursor.CurSelection[1]
		v.Cursor.ResetSelection()
		v.Cursor.StoreVisualX()
	}
	v.Cursor.DownN(v.height)

	return true
}

// HalfPageUp scrolls the view up half a page
func (v *View) HalfPageUp() bool {
	if v.mainCursor() {
		if v.Topline > v.height/2 {
			v.ScrollUp(v.height / 2)
		} else {
			v.Topline = 0
		}
	}
	return false
}

// HalfPageDown scrolls the view down half a page
func (v *View) HalfPageDown() bool {
	if v.mainCursor() {
		if v.Buf.NumLines-(v.Topline+v.height) > v.height/2 {
			v.ScrollDown(v.height / 2)
		} else if v.Buf.NumLines >= v.height {
			v.Topline = v.Buf.NumLines - v.height

		}
	}
	return false
}

// ToggleOverwriteMode lets the user toggle the text overwrite mode
func (v *View) ToggleOverwriteMode() bool {
	if v.mainCursor() {
		v.isOverwriteMode = !v.isOverwriteMode
	}
	return false
}

// Escape leaves current mode
func (v *View) Escape() bool {
	v.done()
	return false
}
