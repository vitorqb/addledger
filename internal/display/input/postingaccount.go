package input

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/controller"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
)

type PostingAccountField struct {
	*tview.InputField
	controller controller.IInputController
}

func NewPostingAccount(
	controller controller.IInputController,
	eventbus eventbusmod.IEventBus,
) *PostingAccountField {
	field := &PostingAccountField{tview.NewInputField(), controller}
	field.SetLabel("Account: ")

	// Custom handling of user input
	field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		logrus.WithField("key", event.Key()).
			WithField("rune", event.Rune()).
			Debug("Received event in PostingAccount input")

		// If key is linked to account on Context's List of Accounts, dispatch
		// and action to it.
		listAction := eventKeyToListAction(event)
		if listAction != listaction.NONE {
			logrus.WithField("action", listAction).Debug("Dispatching listAction")
			controller.OnPostingAccountListAcction(listAction)
			return nil
		}

		switch key := event.Key(); key {
		// If user hit enter...
		case tcell.KeyEnter:
			if field.GetText() == "" {
				// ...with no text, he's done entering stuff!
				field.controller.OnPostingAccountDone("")
			} else {
				// ...with some text written, he's selecting from context
				field.controller.OnPostingAccountSelectedFromContext()
			}
			return nil
		// if Ctrl+J, use input as it is
		case tcell.KeyCtrlJ:
			text := field.GetText()
			field.controller.OnPostingAccountDone(text)
			return nil
		// if Tab then autocompletes
		case tcell.KeyTab:
			field.controller.OnPostingAccountInsertFromContext()
			return nil
		}
		// Else delegates to default
		return event
	})

	// When current text changes make sure controller is aware.
	field.SetChangedFunc(func(text string) {
		logrus.WithField("text", text).Debug("PostingAccount input change")
		field.controller.OnPostingAccountChanged(text)
	})

	// Subscribes to eventbus
	err := eventbus.Subscribe(eventbusmod.Subscription{
		Topic: "input.postingaccount.settext",
		Handler: func(e eventbusmod.Event) {
			text, ok := e.Data.(string)
			if !ok {
				logrus.WithField("event", e).Error("Received invalid event")
				return
			}
			field.SetText(text)
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("could not subscribe to eventbus")
	}

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
