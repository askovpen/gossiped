package ui

import (
	"fmt"
	"github.com/askovpen/gossiped/pkg/msgapi"
	"github.com/rivo/tview"
)

func (a *App) EditMsg(areaId int, msgNum uint32) (string, tview.Primitive, bool, bool) {
	layout := tview.NewFlex()
	return fmt.Sprintf("EditMsg", msgapi.Areas[areaId].GetName(), msgNum), layout, true, true
}
