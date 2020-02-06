package ui

import (
	"github.com/rivo/tview"
)

// App ui struct
type App struct {
	App         *tview.Application
	Layout      *tview.Flex
	Pages       *tview.Pages
	sb          *StatusBar
	al          *tview.Table
	im          IM
	showKludges bool
}

// NewApp return new App
func NewApp() *App {
	a := &App{}
	a.App = tview.NewApplication()

	a.Pages = tview.NewPages()
	a.Pages.AddPage(a.AreaList())
	a.Pages.AddPage(a.AreaListQuit())
	a.Pages.AddPage(a.AreaListHelp())
	//a.Pages.AddPage(a.ViewMsgHelp())

	a.sb = NewStatusBar(a)
	a.sb.Run()
	a.Layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.Pages, 0, 1, true).
		AddItem(a.sb.SB, 1, 1, false)
	return a
}

// Run run App
func (a *App) Run() error {
	return a.App.SetRoot(a.Layout, true).Run()
}
