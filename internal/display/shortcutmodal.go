package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:generate $MOCKGEN --source=shortcutmodal.go --destination=../../mocks/display/shortcutmodal_mock.go

type ShortcutModalController interface {
	OnHideShortcutModal()
}

type ShortcutModal struct {
	*tview.Box
	controller ShortcutModalController
}

func NewShortcutModal(controller ShortcutModalController) *ShortcutModal {
	modal := &ShortcutModal{tview.NewBox(), controller}
	modal.SetBorder(true)
	modal.SetTitle("Shortcuts")
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'q' {
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
