package ui

import(
  "github.com/rivo/tview"
)

var (
  App  *tview.Application
  Pages *tview.Pages
  Status *tview.TextView
  StatusTime *tview.TextView
  showKludges bool
)
