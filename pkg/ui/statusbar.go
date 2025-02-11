package ui

import (
	"time"

	"github.com/askovpen/gossiped/pkg/config"
	"github.com/rivo/tview"
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
	styleText := config.GetElementStyle(config.ColorAreaStatusBar, config.ColorElementText)
	sb.status = tview.NewTextView().SetWrap(false)
	sb.status.SetDynamicColors(true)
	sb.status.SetTextStyle(styleText)
	sb.status.SetChangedFunc(func() {
		sb.app.App.Draw()
	})

	sb.statusTime = tview.NewTextView().SetWrap(false)
	sb.statusTime.SetTextStyle(styleText)
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
	styleText := config.GetElementStyle(config.ColorAreaStatusBar, config.ColorElementText)

	//sb.status.SetText(" [::b][white]" + s)
	sb.status.SetTextStyle(styleText)
	sb.status.SetText(" " + s)
}

// Run update timers
func (sb StatusBar) Run() {
	if config.Config.Statusbar.Clock {
		styleText := config.GetElementStyle(config.ColorAreaStatusBar, config.ColorElementText)
		sb.statusTime.SetTextStyle(styleText)
		sb.statusTime.SetText(time.Now().Format("15:04:05"))
		clock := time.NewTicker(1 * time.Second)
		go func() {
			for t := range clock.C {
				//fmt.Fprintf(ui.StatusTime,"%s",t.Format("15:04:05"))
				sb.statusTime.SetText(t.Format("15:04:05"))
			}
		}()
	}
}
