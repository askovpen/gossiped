package editor

import (
	"time"
)

const (
	// TextEventInsert represents an insertion event
	TextEventInsert = 1
	// TextEventRemove represents a deletion event
	TextEventRemove = -1
	// TextEventReplace represents a replace event
	TextEventReplace = 0
)

// TextEvent holds data for a manipulation on some text that can be undone
type TextEvent struct {
	C Cursor

	EventType int
	Deltas    []Delta
	Time      time.Time
}

// A Delta is a change to the buffer
type Delta struct {
	Text  string
	Start Loc
	End   Loc
}

// ExecuteTextEvent runs a text event
func ExecuteTextEvent(t *TextEvent, buf *Buffer) {
	if t.EventType == TextEventInsert {
		for _, d := range t.Deltas {
			buf.insert(d.Start, []byte(d.Text))
		}
	} else if t.EventType == TextEventRemove {
		for i, d := range t.Deltas {
			t.Deltas[i].Text = buf.remove(d.Start, d.End)
		}
	} else if t.EventType == TextEventReplace {
		for i, d := range t.Deltas {
			t.Deltas[i].Text = buf.remove(d.Start, d.End)
			buf.insert(d.Start, []byte(d.Text))
			t.Deltas[i].Start = d.Start
			t.Deltas[i].End = Loc{d.Start.X + Count(d.Text), d.Start.Y}
		}
		for i, j := 0, len(t.Deltas)-1; i < j; i, j = i+1, j-1 {
			t.Deltas[i], t.Deltas[j] = t.Deltas[j], t.Deltas[i]
		}
	}
}

// EventHandler executes text manipulations and allows undoing and redoing
type EventHandler struct {
	buf *Buffer
}

// NewEventHandler returns a new EventHandler
func NewEventHandler(buf *Buffer) *EventHandler {
	eh := new(EventHandler)
	eh.buf = buf
	return eh
}

// Insert creates an insert text event and executes it
func (eh *EventHandler) Insert(start Loc, text string) {
	e := &TextEvent{
		C:         *eh.buf.cursors[eh.buf.curCursor],
		EventType: TextEventInsert,
		Deltas:    []Delta{{text, start, Loc{0, 0}}},
		Time:      time.Now(),
	}
	eh.Execute(e)
	charCount := Count(text)
	e.Deltas[0].End = start.Move(charCount, eh.buf)
	end := e.Deltas[0].End

	for _, c := range eh.buf.cursors {
		move := func(loc Loc) Loc {
			if start.Y != end.Y && loc.GreaterThan(start) {
				loc.Y += end.Y - start.Y
			} else if loc.Y == start.Y && loc.GreaterEqual(start) {
				loc = loc.Move(charCount, eh.buf)
			}
			return loc
		}
		c.Loc = move(c.Loc)
		c.CurSelection[0] = move(c.CurSelection[0])
		c.CurSelection[1] = move(c.CurSelection[1])
		c.OrigSelection[0] = move(c.OrigSelection[0])
		c.OrigSelection[1] = move(c.OrigSelection[1])
		c.LastVisualX = c.GetVisualX()
	}
}

// Remove creates a remove text event and executes it
func (eh *EventHandler) Remove(start, end Loc) {
	e := &TextEvent{
		C:         *eh.buf.cursors[eh.buf.curCursor],
		EventType: TextEventRemove,
		Deltas:    []Delta{{"", start, end}},
		Time:      time.Now(),
	}
	eh.Execute(e)

	for _, c := range eh.buf.cursors {
		move := func(loc Loc) Loc {
			if start.Y != end.Y && loc.GreaterThan(end) {
				loc.Y -= end.Y - start.Y
			} else if loc.Y == end.Y && loc.GreaterEqual(end) {
				loc = loc.Move(-Diff(start, end, eh.buf), eh.buf)
			}
			return loc
		}
		c.Loc = move(c.Loc)
		c.CurSelection[0] = move(c.CurSelection[0])
		c.CurSelection[1] = move(c.CurSelection[1])
		c.OrigSelection[0] = move(c.OrigSelection[0])
		c.OrigSelection[1] = move(c.OrigSelection[1])
		c.LastVisualX = c.GetVisualX()
	}
}

// MultipleReplace creates an multiple insertions executes them
func (eh *EventHandler) MultipleReplace(deltas []Delta) {
	e := &TextEvent{
		C:         *eh.buf.cursors[eh.buf.curCursor],
		EventType: TextEventReplace,
		Deltas:    deltas,
		Time:      time.Now(),
	}
	eh.Execute(e)
}

// Replace deletes from start to end and replaces it with the given string
func (eh *EventHandler) Replace(start, end Loc, replace string) {
	eh.Remove(start, end)
	eh.Insert(start, replace)
}

// Execute a textevent and add it to the undo stack
func (eh *EventHandler) Execute(t *TextEvent) {
	ExecuteTextEvent(t, eh.buf)
}
