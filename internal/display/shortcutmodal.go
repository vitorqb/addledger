package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:generate $MOCKGEN --source=shortcutmodal.go --destination=../../mocks/display/shortcutmodal_mock.go

type ShortcutModalController interface {
	OnHideShortcutModal()
	// Aciton to discard the current loaded statement.
	OnDiscardStatement()
	// Action to load a new statement.
	OnLoadStatement()
}

type ShortcutModal struct {
	*tview.TextView
	controller ShortcutModalController
}

func getBodyText() string {
	return strings.Trim(
		"d - Discard statement\n"+
			"l - Load statement\n"+
			"q - Quit\n",
		"\n",
	)
}

func NewShortcutModal(controller ShortcutModalController) *ShortcutModal {
	modal := &ShortcutModal{tview.NewTextView(), controller}
	modal.SetBorder(true)
	modal.SetTitle("Shortcuts")
	modal.SetText(getBodyText())
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'l':
				modal.controller.OnLoadStatement()
				modal.controller.OnHideShortcutModal()
				return nil
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
