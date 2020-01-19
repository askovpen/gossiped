package editor

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/askovpen/gossiped/pkg/highlight"
)

// LargeFileThreshold const
const LargeFileThreshold = 50000

var (
	// 0 - no line type detected
	// 1 - lf detected
	// 2 - crlf detected
	fileformat = 0
)

// Buffer stores the text for files that are loaded into the text editor
// It uses a rope to efficiently store the string and contains some
// simple functions for saving and wrapper functions for modifying the rope
type Buffer struct {
	// The eventhandler for undo/redo
	*EventHandler
	// This stores all the text in the buffer as an array of lines
	*LineArray

	// The path to the loaded file, if any
	Path string

	Cursor    Cursor
	cursors   []*Cursor // for multiple cursors
	curCursor int       // the current cursor

	// Name of the buffer on the status line
	name string

	// Whether or not the buffer has been modified since it was opened
	IsModified bool

	// NumLines is the number of lines in the buffer
	NumLines int

	syntaxDef   *highlight.Def
	highlighter *highlight.Highlighter

	// Buffer local settings
	Settings map[string]interface{}
}

// NewBufferFromString creates a new buffer containing the given string
func NewBufferFromString(text string) *Buffer {
	return NewBuffer(strings.NewReader(text), int64(len(text)), "1.msg", nil)
}

// NewBuffer creates a new buffer from a given reader
func NewBuffer(reader io.Reader, size int64, path string, cursorPosition []string) *Buffer {
	b := new(Buffer)
	b.LineArray = NewLineArray(size, reader)

	b.Settings = DefaultLocalSettings()
	//	for k, v := range globalSettings {
	//		if _, ok := b.Settings[k]; ok {
	//			b.Settings[k] = v
	//		}
	//	}

	if fileformat == 1 {
		b.Settings["fileformat"] = "unix"
	} else if fileformat == 2 {
		b.Settings["fileformat"] = "dos"
	}

	b.Path = path

	b.EventHandler = NewEventHandler(b)

	b.update()

	b.Cursor = Cursor{
		Loc: Loc{0, 0},
		buf: b,
	}

	//InitLocalSettings(b)

	b.cursors = []*Cursor{&b.Cursor}

	return b
}

// GetName returns the name that should be displayed in the statusline
// for this buffer
func (b *Buffer) GetName() string {
	return b.name
}

// updateRules updates the syntax rules and filetype for this buffer
// This is called when the colorscheme changes
func (b *Buffer) updateRules() {
	ryaml := `
filetype: msg
detect:
  filename: "\\.msg$"
rules:
- comment: ".*\\>+.*$"
- icomment: ".*[^>](>>)+[^>].*$"
- tagline: "^\\.\\.\\..*$"
- origin: "^ \\* Origin:.*$"
- tearline: "^--- .*$"
- kludge: "^@.*$"
- kludge: "^SEEN-BY: .*$"
`
	rehighlight := false
	file, err := highlight.ParseFile([]byte(ryaml))
	if err != nil {
		return
	}
	ftdetect, err := highlight.ParseFtDetect(file)
	if err != nil {
		return
	}
	header := new(highlight.Header)
	header.FileType = file.FileType
	header.FtDetect = ftdetect
	b.syntaxDef, err = highlight.ParseDef(file, header)
	if err != nil {
		return
	}
	rehighlight = true
	if b.highlighter == nil || rehighlight {
		if b.syntaxDef != nil {
			b.Settings["filetype"] = b.syntaxDef.FileType
			b.highlighter = highlight.NewHighlighter(b.syntaxDef)
			if b.Settings["syntax"].(bool) {
				b.highlighter.HighlightStates(b)
			}
		}
	}

}

// FileType returns the buffer's filetype
func (b *Buffer) FileType() string {
	return b.Settings["filetype"].(string)
}

// IndentString returns a string representing one level of indentation
func (b *Buffer) IndentString() string {
	if b.Settings["tabstospaces"].(bool) {
		return Spaces(int(b.Settings["tabsize"].(float64)))
	}
	return "\t"
}

// Update fetches the string from the rope and updates the `text` and `lines` in the buffer
func (b *Buffer) update() {
	b.NumLines = len(b.lines)
}

// MergeCursors merges any cursors that are at the same position
// into one cursor
func (b *Buffer) MergeCursors() {
	var cursors []*Cursor
	for i := 0; i < len(b.cursors); i++ {
		c1 := b.cursors[i]
		if c1 != nil {
			for j := 0; j < len(b.cursors); j++ {
				c2 := b.cursors[j]
				if c2 != nil && i != j && c1.Loc == c2.Loc {
					b.cursors[j] = nil
				}
			}
			cursors = append(cursors, c1)
		}
	}

	b.cursors = cursors

	for i := range b.cursors {
		b.cursors[i].Num = i
	}

	if b.curCursor >= len(b.cursors) {
		b.curCursor = len(b.cursors) - 1
	}
}

// UpdateCursors updates all the cursors indicies
func (b *Buffer) UpdateCursors() {
	for i, c := range b.cursors {
		c.Num = i
	}
}

// Modified returns if this buffer has been modified since
// being opened
func (b *Buffer) Modified() bool {
	return b.IsModified
}

func (b *Buffer) insert(pos Loc, value []byte) {
	b.IsModified = true
	b.LineArray.insert(pos, value)
	b.update()
}
func (b *Buffer) remove(start, end Loc) string {
	b.IsModified = true
	sub := b.LineArray.remove(start, end)
	b.update()
	return sub
}
//func (b *Buffer) deleteToEnd(start Loc) {
//	b.IsModified = true
//	b.LineArray.DeleteToEnd(start)
//	b.update()
//}

// Start returns the location of the first character in the buffer
func (b *Buffer) Start() Loc {
	return Loc{0, 0}
}

// End returns the location of the last character in the buffer
func (b *Buffer) End() Loc {
	return Loc{utf8.RuneCount(b.lines[b.NumLines-1].data), b.NumLines - 1}
}

// RuneAt returns the rune at a given location in the buffer
func (b *Buffer) RuneAt(loc Loc) rune {
	line := b.LineRunes(loc.Y)
	if len(line) > 0 {
		return line[loc.X]
	}
	return '\n'
}

// LineBytes returns a single line as an array of runes
func (b *Buffer) LineBytes(n int) []byte {
	if n >= len(b.lines) {
		return []byte{}
	}
	return b.lines[n].data
}

// LineRunes returns a single line as an array of runes
func (b *Buffer) LineRunes(n int) []rune {
	if n >= len(b.lines) {
		return []rune{}
	}
	return toRunes(b.lines[n].data)
}

// Line returns a single line
func (b *Buffer) Line(n int) string {
	if n >= len(b.lines) {
		return ""
	}
	return string(b.lines[n].data)
}

// LinesNum returns the number of lines in the buffer
func (b *Buffer) LinesNum() int {
	return len(b.lines)
}

// Lines returns an array of strings containing the lines from start to end
func (b *Buffer) Lines(start, end int) []string {
	lines := b.lines[start:end]
	var slice []string
	for _, line := range lines {
		slice = append(slice, string(line.data))
	}
	return slice
}

// Len gives the length of the buffer
func (b *Buffer) Len() (n int) {
	for _, l := range b.lines {
		n += utf8.RuneCount(l.data)
	}

	if len(b.lines) > 1 {
		n += len(b.lines) - 1 // account for newlines
	}

	return
}

// MoveLinesUp moves the range of lines up one row
func (b *Buffer) MoveLinesUp(start int, end int) {
	// 0 < start < end <= len(b.lines)
	if start < 1 || start >= end || end > len(b.lines) {
		return // what to do? FIXME
	}
	if end == len(b.lines) {
		b.Insert(
			Loc{
				utf8.RuneCount(b.lines[end-1].data),
				end - 1,
			},
			"\n"+b.Line(start-1),
		)
	} else {
		b.Insert(
			Loc{0, end},
			b.Line(start-1)+"\n",
		)
	}
	b.Remove(
		Loc{0, start - 1},
		Loc{0, start},
	)
}

// MoveLinesDown moves the range of lines down one row
func (b *Buffer) MoveLinesDown(start int, end int) {
	// 0 <= start < end < len(b.lines)
	// if end == len(b.lines), we can't do anything here because the
	// last line is unaccessible, FIXME
	if start < 0 || start >= end || end >= len(b.lines)-1 {
		return // what to do? FIXME
	}
	b.Insert(
		Loc{0, start},
		b.Line(end)+"\n",
	)
	end++
	b.Remove(
		Loc{0, end},
		Loc{0, end + 1},
	)
}

// ClearMatches clears all of the syntax highlighting for this buffer
func (b *Buffer) ClearMatches() {
	for i := range b.lines {
		b.SetMatch(i, nil)
		b.SetState(i, nil)
	}
}

func (b *Buffer) clearCursors() {
	for i := 1; i < len(b.cursors); i++ {
		b.cursors[i] = nil
	}
	b.cursors = b.cursors[:1]
	b.UpdateCursors()
	b.Cursor.ResetSelection()
}
