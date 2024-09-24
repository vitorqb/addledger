package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:generate $MOCKGEN --source=shortcutmodal.go --destination=../../mocks/display/shortcutmodal_mock.go

type ShortcutModalController interface {
	OnHideShortcutModal()
	// Displays the statement modal
	OnShowStatementModal()
}

type ShortcutModal struct {
	*tview.TextView
	controller ShortcutModalController
}

func getBodyText() string {
	return strings.Trim(
		"s - Statement modal\n"+
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
			case 's':
				modal.controller.OnShowStatementModal()
				modal.controller.OnHideShortcutModal()
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
