package input

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/listaction"
)

type PostingAccountField struct {
	*tview.InputField
	controller controller.IInputController
}

func NewPostingAccount(controller controller.IInputController) *PostingAccountField {
	field := &PostingAccountField{tview.NewInputField(), controller}
	field.SetLabel("Account: ")

	// When done, send info to controller
	field.SetDoneFunc(func(_ tcell.Key) {
		text := field.GetText()
		field.controller.OnPostingAccountDone(text)
	})

	// When receive a key, maybe dispatch an action to the context for
	// autocompletion.
	field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		listAction := eventKeyToListAction(event)
		if listAction != listaction.NONE {
			controller.OnPostingAccountListAcction(listAction)
			return nil
		}
		return event
	})

	// When current text changes make sure controller is aware.
	field.SetChangedFunc(func(text string) {
		field.controller.OnPostingAccountChanged(text)
	})
	return field
}

func eventKeyToListAction(event *tcell.EventKey) listaction.ListAction {
	switch key := event.Key(); key {
	case tcell.KeyDown, tcell.KeyCtrlN:
		return listaction.NEXT
	case tcell.KeyUp, tcell.KeyCtrlP:
		return listaction.PREV
	}
	return listaction.NONE
}
