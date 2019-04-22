package ui

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

var (
	openColorRegex = regexp.MustCompile(`\[([a-zA-Z]*|#[0-9a-zA-Z]*)$`)
	newLineRegex   = regexp.MustCompile(`\r?\n`)

	TabSize = 4
)

type cursor struct {
	X int
	Y int
}

type textViewIndex struct {
	Line            int    // The index into the "buffer" variable.
	Pos             int    // The index into the "buffer" string (byte position).
	NextPos         int    // The (byte) index of the next character in this buffer line.
	Width           int    // The screen width of this line.
	ForegroundColor string // The starting foreground color ("" = don't change, "-" = reset).
	BackgroundColor string // The starting background color ("" = don't change, "-" = reset).
	Attributes      string // The starting attributes ("" = don't change, "-" = reset).
}

type EditBody struct {
	sync.Mutex
	*tview.Box
	buffer      []string
	recentBytes []byte
	index       []*textViewIndex
	align       int

	lastWidth int

	longestLine int

	lineOffset int

	trackEnd bool

	columnOffset int

	pageSize int

	wrap bool

	wordWrap bool

	textColor tcell.Color

	dynamicColors bool

	changed func()

	done func(tcell.Key)

	cur cursor
}

func NewEditBody() *EditBody {
	return &EditBody{
		Box:           tview.NewBox(),
		lineOffset:    -1,
		align:         tview.AlignLeft,
		wrap:          true,
		wordWrap:      true,
		textColor:     tview.Styles.PrimaryTextColor,
		dynamicColors: false,
	}
}
func (t *EditBody) getCursor(p int) cursor {
	log.Printf("p: %d", p)
	_, _, width, _ := t.GetInnerRect()
	c := cursor{}
	t.reindexBuffer(width)
	//pt:=0
	for line := 0; line < len(t.index); line++ {
		index := t.index[line]
		if index.Line == p {
			c.Y = line
			break
		}
	}

	return c
}

func (t *EditBody) SetText(text string, p int) *EditBody {
	t.Clear()
	fmt.Fprint(t, text)
	t.cur = t.getCursor(p)
	t.ScrollTo(t.cur.Y, 0)
	return t
}

func (t *EditBody) GetText(stripTags bool) string {
	buffer := t.buffer
	if !stripTags {
		buffer = append(buffer, string(t.recentBytes))
	}

	// Add newlines again.
	text := strings.Join(buffer, "\n")

	// Strip from tags if required.
	if stripTags {
		if t.dynamicColors {
			text = colorPattern.ReplaceAllString(text, "")
		}
	}

	return text
}

func (t *EditBody) SetChangedFunc(handler func()) *EditBody {
	t.changed = handler
	return t
}

func (t *EditBody) SetDoneFunc(handler func(key tcell.Key)) *EditBody {
	t.done = handler
	return t
}

func (t *EditBody) ScrollTo(row, column int) *EditBody {
	t.lineOffset = row
	t.columnOffset = column
	return t
}

func (t *EditBody) ScrollToBeginning() *EditBody {
	t.trackEnd = false
	t.lineOffset = 0
	t.columnOffset = 0
	return t
}

func (t *EditBody) ScrollToEnd() *EditBody {
	t.trackEnd = true
	t.columnOffset = 0
	return t
}

func (t *EditBody) GetScrollOffset() (row, column int) {
	return t.lineOffset, t.columnOffset
}

func (t *EditBody) Clear() *EditBody {
	t.buffer = nil
	t.recentBytes = nil
	t.index = nil
	return t
}

func (t *EditBody) Write(p []byte) (n int, err error) {
	// Notify at the end.
	t.Lock()
	changed := t.changed
	t.Unlock()
	if changed != nil {
		defer changed() // Deadlocks may occur if we lock here.
	}

	t.Lock()
	defer t.Unlock()

	// Copy data over.
	newBytes := append(t.recentBytes, p...)
	t.recentBytes = nil

	// If we have a trailing invalid UTF-8 byte, we'll wait.
	if r, _ := utf8.DecodeLastRune(p); r == utf8.RuneError {
		t.recentBytes = newBytes
		return len(p), nil
	}

	// If we have a trailing open dynamic color, exclude it.
	if t.dynamicColors {
		location := openColorRegex.FindIndex(newBytes)
		if location != nil {
			t.recentBytes = newBytes[location[0]:]
			newBytes = newBytes[:location[0]]
		}
	}

	// If we have a trailing open region, exclude it.
	// Transform the new bytes into strings.
	newBytes = bytes.Replace(newBytes, []byte{'\t'}, bytes.Repeat([]byte{' '}, TabSize), -1)
	for index, line := range newLineRegex.Split(string(newBytes), -1) {
		if index == 0 {
			if len(t.buffer) == 0 {
				t.buffer = []string{line}
			} else {
				t.buffer[len(t.buffer)-1] += line
			}
		} else {
			t.buffer = append(t.buffer, line)
		}
	}

	// Reset the index.
	t.index = nil

	return len(p), nil
}

// reindexBuffer re-indexes the buffer such that we can use it to easily draw
// the buffer onto the screen. Each line in the index will contain a pointer
// into the buffer from which on we will print text. It will also contain the
// color with which the line starts.
func (t *EditBody) reindexBuffer(width int) {
	if t.index != nil {
		return // Nothing has changed. We can still use the current index.
	}
	t.index = nil

	// If there's no space, there's no index.
	if width < 1 {
		return
	}

	// Initial states.

	// Go through each line in the buffer.
	for bufferIndex, str := range t.buffer {
		_, _, _, _, _, strippedStr, _ := decomposeString(str, t.dynamicColors)

		// Split the line if required.
		var splitLines []string
		str = strippedStr
		if t.wrap && len(str) > 0 {
			for len(str) > 0 {
				extract := runewidth.Truncate(str, width, "")
				if t.wordWrap && len(extract) < len(str) {
					// Add any spaces from the next line.
					if spaces := spacePattern.FindStringIndex(str[len(extract):]); spaces != nil && spaces[0] == 0 {
						extract = str[:len(extract)+spaces[1]]
					}

					// Can we split before the mandatory end?
					matches := boundaryPattern.FindAllStringIndex(extract, -1)
					if len(matches) > 0 {
						// Yes. Let's split there.
						extract = extract[:matches[len(matches)-1][1]]
					}
				}
				splitLines = append(splitLines, extract)
				str = str[len(extract):]
			}
		} else {
			// No need to split the line.
			splitLines = []string{str}
		}

		// Create index from split lines.
		var (
			originalPos                                  int
			foregroundColor, backgroundColor, attributes string
		)
		for _, splitLine := range splitLines {
			line := &textViewIndex{
				Line:            bufferIndex,
				Pos:             originalPos,
				ForegroundColor: foregroundColor,
				BackgroundColor: backgroundColor,
				Attributes:      attributes,
			}

			// Shift original position with tags.
			lineLength := len(splitLine)
			//remainingLength := lineLength
			//tagEnd := originalPos
			totalTagLength := 0

			// Advance to next line.
			originalPos += lineLength + totalTagLength

			// Append this line.
			line.NextPos = originalPos
			line.Width = stringWidth(splitLine)
			t.index = append(t.index, line)
		}

		// Word-wrapped lines may have trailing whitespace. Remove it.
		if t.wrap && t.wordWrap {
			for _, line := range t.index {
				str := t.buffer[line.Line][line.Pos:line.NextPos]
				spaces := spacePattern.FindAllStringIndex(str, -1)
				if spaces != nil && spaces[len(spaces)-1][1] == len(str) {
					oldNextPos := line.NextPos
					line.NextPos -= spaces[len(spaces)-1][1] - spaces[len(spaces)-1][0]
					line.Width -= stringWidth(t.buffer[line.Line][line.NextPos:oldNextPos])
				}
			}
		}
	}

	// Calculate longest line.
	t.longestLine = 0
	for _, line := range t.index {
		if line.Width > t.longestLine {
			t.longestLine = line.Width
		}
	}
}

// Draw draws this primitive onto the screen.
func (t *EditBody) Draw(screen tcell.Screen) {
	t.Lock()
	defer t.Unlock()
	t.Box.Draw(screen)

	// Get the available size.
	x, y, width, height := t.GetInnerRect()
	t.pageSize = height

	// If the width has changed, we need to reindex.
	if width != t.lastWidth && t.wrap {
		t.index = nil
	}
	t.lastWidth = width

	// Re-index.
	t.reindexBuffer(width)

	// If we don't have an index, there's nothing to draw.
	if t.index == nil {
		return
	}

	// Move to highlighted regions.

	// Adjust line offset.
	if t.lineOffset+height > len(t.index) {
		t.trackEnd = true
	}
	if t.trackEnd {
		t.lineOffset = len(t.index) - height
	}
	if t.lineOffset < 0 {
		t.lineOffset = 0
	}

	// Adjust column offset.
	if t.align == tview.AlignLeft {
		if t.columnOffset+width > t.longestLine {
			t.columnOffset = t.longestLine - width
		}
		if t.columnOffset < 0 {
			t.columnOffset = 0
		}
	}

	// Draw the buffer.
	defaultStyle := tcell.StyleDefault.Foreground(t.textColor)
	for line := t.lineOffset; line < len(t.index); line++ {
		// Are we done?
		if line-t.lineOffset >= height {
			break
		}

		// Get the text for this line.
		index := t.index[line]
		text := t.buffer[index.Line][index.Pos:index.NextPos]
		foregroundColor := index.ForegroundColor
		backgroundColor := index.BackgroundColor
		attributes := index.Attributes
		re := regexp.MustCompile(">+")
		if len(text) > 10 && text[0:11] == " * Origin: " {
			attributes = "b"
		} else if len(text) > 3 && text[0:4] == "--- " {
			attributes = "b"
		} else if len(text) > 3 && text[0:4] == "... " {
			attributes = "b"
		} else if ind := re.FindStringIndex(text); ind != nil {
			ind2 := strings.Index(text, "<")
			if (ind2 == -1 || ind2 > ind[1]) && ind[0] < 6 {
				attributes = "b"
				if (ind[1]-ind[0])%2 != 0 {
					foregroundColor = "yellow"
				}
			}
		}

		// Process tags.
		colorTagIndices, colorTags, _, _, escapeIndices, strippedText, _ := decomposeString(text, t.dynamicColors)

		// Calculate the position of the line.
		var skip, posX int
		if t.align == tview.AlignLeft {
			posX = -t.columnOffset
		}
		if posX < 0 {
			skip = -posX
			posX = 0
		}

		// Print the line.
		var colorPos, escapePos, tagOffset, skipped int
		iterateString(strippedText, func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
			// Process tags.
			for {
				if colorPos < len(colorTags) && textPos+tagOffset >= colorTagIndices[colorPos][0] && textPos+tagOffset < colorTagIndices[colorPos][1] {
					// Get the color.
					foregroundColor, backgroundColor, attributes = styleFromTag(foregroundColor, backgroundColor, attributes, colorTags[colorPos])
					tagOffset += colorTagIndices[colorPos][1] - colorTagIndices[colorPos][0]
					colorPos++
				} else {
					break
				}
			}

			// Skip the second-to-last character of an escape tag.
			if escapePos < len(escapeIndices) && textPos+tagOffset == escapeIndices[escapePos][1]-2 {
				tagOffset++
				escapePos++
			}

			// Mix the existing style with the new style.
			_, _, existingStyle, _ := screen.GetContent(x+posX, y+line-t.lineOffset)
			_, background, _ := existingStyle.Decompose()
			style := overlayStyle(background, defaultStyle, foregroundColor, backgroundColor, attributes)

			// Skip to the right.
			if !t.wrap && skipped < skip {
				skipped += screenWidth
				return false
			}

			// Stop at the right border.
			if posX+screenWidth > width {
				return true
			}

			// Draw the character.
			for offset := screenWidth - 1; offset >= 0; offset-- {
				if offset == 0 {
					screen.SetContent(x+posX+offset, y+line-t.lineOffset, main, comb, style)
				} else {
					screen.SetContent(x+posX+offset, y+line-t.lineOffset, ' ', nil, style)
				}
			}

			// Advance.
			posX += screenWidth
			return false
		})
	}

	screen.ShowCursor(t.cur.X, y+t.cur.Y-t.lineOffset)
}
func (t *EditBody) getRealPos(str string, tp int) int {
	p := 0
	pt := 0
	if tp == 0 {
		return p
	}
	iterateString(str, func(main rune, comb []rune, textPos, textWidth, screenPos, screenWidth int) bool {
		p += len(string(main))
		pt += 1
		if tp == pt {
			return true
		}
		return false
	})
	return p
}

// InputHandler returns the handler for this primitive.
func (t *EditBody) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		_, _, width, _ := t.GetInnerRect()
		add := func(r rune) {
			//newText := i.text[:i.cursorPos] + string(r) + i.text[i.cursorPos:]
			line := t.index[t.cur.Y]
			pos := t.getRealPos(t.buffer[line.Line][line.Pos:line.NextPos], t.cur.X)
			t.buffer[line.Line] = t.buffer[line.Line][:line.Pos+pos] + string(r) + t.buffer[line.Line][line.Pos+pos:]
			line.NextPos += len(string(r))
			if t.cur.X < stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
				t.cur.X++
				if t.cur.X == width {
					t.cur.Y++
					t.cur.X = 0
				}
			}
			//t.index=nil
		}
		ent := func() {
			line := t.index[t.cur.Y]
			pos := t.getRealPos(t.buffer[line.Line][line.Pos:line.NextPos], t.cur.X)
			if t.cur.X == stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
				t.buffer = append(t.buffer, "")
				copy(t.buffer[line.Line+2:], t.buffer[line.Line+1:])
				t.buffer[line.Line+1] = ""
				t.cur.X = 0
				t.cur.Y++
				t.index = nil
			} else if t.cur.X > 0 {
				tmp := t.buffer[line.Line][line.Pos+pos : line.NextPos]
				t.buffer = append(t.buffer, "")
				copy(t.buffer[line.Line+2:], t.buffer[line.Line+1:])
				t.buffer[line.Line+1] = tmp
				t.buffer[line.Line] = t.buffer[line.Line][line.Pos:pos]
				t.cur.X = 0
				t.cur.Y++
				t.index = nil

			} else {
				t.buffer = append(t.buffer, "")
				copy(t.buffer[line.Line+1:], t.buffer[line.Line:])
				t.buffer[line.Line] = ""
				t.cur.Y++
				t.index = nil
			}
			if t.cur.Y-t.lineOffset >= t.pageSize {
				t.lineOffset++
			}

		}
		bck := func() {
			line := t.index[t.cur.Y]
			pos := t.getRealPos(t.buffer[line.Line][line.Pos:line.NextPos], t.cur.X)
			if t.cur.X > 0 {
				//log.Print(charWidth(t.buffer[line.Line], t.cur.X-1))
				t.buffer[line.Line] = t.buffer[line.Line][:line.Pos+pos-charWidth(t.buffer[line.Line], t.cur.X-1)] + t.buffer[line.Line][line.Pos+pos:]
				t.cur.X--
				t.index[line.Line].NextPos--
				//t.index = nil
			} else if t.cur.Y > 0 {
				ln := stringWidth(t.buffer[line.Line])
				t.buffer[line.Line-1] = t.buffer[line.Line-1] + t.buffer[line.Line]
				t.buffer = append(t.buffer[:line.Line], t.buffer[line.Line+1:]...)
				t.cur.Y--
				t.cur.X = stringWidth(t.buffer[line.Line-1]) - ln
				t.index = nil
			}
			if t.cur.Y-t.lineOffset == 0 {
				t.lineOffset--
			}

		}
		del := func() {
			line := t.index[t.cur.Y]
			pos := t.getRealPos(t.buffer[line.Line][line.Pos:line.NextPos], t.cur.X)
			if t.cur.X < stringWidth(t.buffer[line.Line]) {
				t.buffer[line.Line] = t.buffer[line.Line][:line.Pos+pos] + t.buffer[line.Line][line.Pos+pos+charWidth(t.buffer[line.Line], t.cur.X):]
				//t.index[line.Line].NextPos--
				t.index = nil
			} else if t.cur.Y < len(t.index)-1 {
				t.buffer[line.Line] = t.buffer[line.Line] + t.buffer[line.Line+1]
				t.buffer = append(t.buffer[:line.Line+1], t.buffer[line.Line+2:]...)
				t.index = nil
			}
		}
		key := event.Key()
		switch key {
		case tcell.KeyRune:
			add(event.Rune())
		case tcell.KeyHome:
			t.cur.X = 0
		case tcell.KeyEnd:
			line := t.index[t.cur.Y]
			t.cur.X = stringWidth(t.buffer[line.Line])
		case tcell.KeyUp:
			t.trackEnd = false
			if t.cur.Y > 0 {
				t.cur.Y--
				line := t.index[t.cur.Y]
				if t.cur.X >= stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
					t.cur.X = stringWidth(t.buffer[line.Line][line.Pos:line.NextPos])
				}
			}
			if t.cur.Y-t.lineOffset+1 <= 0 {
				t.lineOffset--
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			bck()
		case tcell.KeyDelete:
			del()
		case tcell.KeyEnter:
			ent()
		case tcell.KeyDown:
			if t.cur.Y < len(t.index)-1 {
				t.cur.Y++
				line := t.index[t.cur.Y]
				if t.cur.X >= stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
					t.cur.X = stringWidth(t.buffer[line.Line][line.Pos:line.NextPos])
				}
			}
			if t.cur.Y-t.lineOffset >= t.pageSize {
				t.lineOffset++
			}
		case tcell.KeyLeft:
			if t.cur.X > 0 {
				t.cur.X--
			}
			//t.columnOffset--
		case tcell.KeyRight:
			line := t.index[t.cur.Y]
			if t.cur.X < stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
				t.cur.X++
			}
			//t.columnOffset++

		case tcell.KeyPgDn, tcell.KeyCtrlF:
			if t.cur.Y < len(t.index)-t.pageSize-1 {
				t.cur.Y += t.pageSize
				t.lineOffset += t.pageSize
			} else {
				t.cur.Y = len(t.index) - 1
			}
			line := t.index[t.cur.Y]
			if t.cur.X >= stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
				t.cur.X = stringWidth(t.buffer[line.Line][line.Pos:line.NextPos])
			}

		case tcell.KeyPgUp, tcell.KeyCtrlB:
			if t.cur.Y > t.pageSize {
				t.cur.Y -= t.pageSize
				t.lineOffset -= t.pageSize
				t.trackEnd = false
			} else {
				t.lineOffset = 0
				t.cur.Y = 0
			}
			line := t.index[t.cur.Y]
			if t.cur.X >= stringWidth(t.buffer[line.Line][line.Pos:line.NextPos]) {
				t.cur.X = stringWidth(t.buffer[line.Line][line.Pos:line.NextPos])
			}

		}
	})
}
