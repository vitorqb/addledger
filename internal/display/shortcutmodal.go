package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:generate $MOCKGEN --source=shortcutmodal.go --destination=../../mocks/display/shortcutmodal_mock.go

type ShortcutModalController interface {
	OnHideShortcutModal()
	OnDiscardStatement()
}

type ShortcutModal struct {
	*tview.TextView
	controller ShortcutModalController
}

func NewShortcutModal(controller ShortcutModalController) *ShortcutModal {
	modal := &ShortcutModal{tview.NewTextView(), controller}
	modal.SetBorder(true)
	modal.SetTitle("Shortcuts")
	modal.SetText("d - Discard statement\nq - Quit")
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'd':
				modal.controller.OnDiscardStatement()
				modal.controller.OnHideShortcutModal()
				return nil
			case 'q':
				modal.controller.OnHideShortcutModal()
				return nil
			}
		case tcell.KeyEscape, tcell.KeyCtrlQ:
			modal.controller.OnHideShortcutModal()
			return nil
		}
		return event
	})
	return modal
}
