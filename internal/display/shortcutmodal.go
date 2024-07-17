package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

//go:generate $MOCKGEN --source=shortcutmodal.go --destination=../../mocks/display/shortcutmodal_mock.go

type ShortcutModalController interface {
	OnHideShortcutModal()
	// !!!! TODO DELETE
	// Aciton to discard the current loaded statement.
	OnPopStatement()
	// !!!! TODO DELETE
	// Action to load a new statement.
	OnLoadStatementRequest()
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
			case 's':
				modal.controller.OnShowStatementModal()
				modal.controller.OnHideShortcutModal()
			case 'l':
				modal.controller.OnLoadStatementRequest()
				modal.controller.OnHideShortcutModal()
				return nil
			case 'd':
				modal.controller.OnPopStatement()
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
