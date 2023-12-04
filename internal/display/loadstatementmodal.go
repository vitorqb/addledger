package display

import "github.com/rivo/tview"

type (
	LoadStatementModal struct {
		*tview.Box
	}
)

func NewLoadStatementModal() *LoadStatementModal {
	box := tview.NewBox()
	box.SetBorder(true)
	box.SetTitle("Load Statement")
	return &LoadStatementModal{box}
}
