package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
)

// StatusBar struct
type StatusBar struct {
	SB         *tview.Flex
	status     *tview.TextView
	statusTime *tview.TextView
	app        *App
}

// NewStatusBar func
func NewStatusBar(app *App) *StatusBar {
	sb := &StatusBar{}

	sb.app = app

	sb.status = tview.NewTextView().SetWrap(false)
	sb.status.SetBackgroundColor(tcell.ColorNavy)
	sb.status.SetTextColor(tcell.ColorYellow)
	sb.status.SetDynamicColors(true)
	sb.status.SetChangedFunc(func() {
		sb.app.App.Draw()
	})

	sb.statusTime = tview.NewTextView().SetWrap(false)
	sb.statusTime.SetBackgroundColor(tcell.ColorNavy)
	sb.statusTime.SetTextColor(tcell.ColorYellow)
	sb.statusTime.SetDynamicColors(true)
	sb.statusTime.SetChangedFunc(func() {
		sb.app.App.Draw()
	})

	sb.SB = tview.NewFlex().
		AddItem(sb.status, 0, 1, false).
		AddItem(sb.statusTime, 10, 1, false)
	return sb
}

// SetStatus set status
func (sb StatusBar) SetStatus(s string) {
	sb.status.SetText(" [::b][white]" + s)
}

// Run update timers
func (sb StatusBar) Run() {
	sb.statusTime.SetText("[::b][white]" + time.Now().Format("15:04:05"))
	clock := time.NewTicker(1 * time.Second)
	go func() {
		for t := range clock.C {
			//fmt.Fprintf(ui.StatusTime,"%s",t.Format("15:04:05"))
			sb.statusTime.SetText("[::b][white]" + t.Format("15:04:05"))
		}
	}()
}
