package input

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/controller"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
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

		switch key := event.Key(); key {
		// Arrow down -> move contextual list down
		case tcell.KeyDown, tcell.KeyCtrlN:
			field.controller.OnPostingAccountListAcction(listaction.NEXT)
			return nil
		// Arrow up -> move contextual list up
		case tcell.KeyUp, tcell.KeyCtrlP:
			field.controller.OnPostingAccountListAcction(listaction.PREV)
			return nil
		// If user hit enter, select from context (list)
		case tcell.KeyEnter:
			field.controller.OnPostingAccountDone(input.Context)
			return nil
		// if Ctrl+J, use input as it is
		case tcell.KeyCtrlJ:
			field.controller.OnPostingAccountChanged(field.GetText())
			field.controller.OnPostingAccountDone(input.Input)
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
