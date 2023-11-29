package display

import (
	"github.com/rivo/tview"
)

// Source: https://github.com/rivo/tview/wiki/Modal
func center(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}
