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
	field.SetDoneFunc(func(_ tcell.Key) {
		text := field.GetText()
		field.controller.OnPostingAccountDone(text)
	})
	field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		listAction := eventKeyToListAction(event)
		if listAction != listaction.NONE {
			controller.OnPostingAccountListAcction(listAction)
			return nil
		}
		return event
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
