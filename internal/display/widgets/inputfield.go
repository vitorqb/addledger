package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/listaction"
)

// InputField is a wrapper of a tview.InputField that adds a few extra
// features.
type InputField struct {
	*tview.InputField
}

// NewInputField creates a new InputField.
func NewInputField() *InputField { return &InputField{InputField: tview.NewInputField()} }

// ContextualListHooks is a set of hooks that connects the input field to a
// contextual list.
type ContextualListLinkOpts struct {
	// InputName is the name of the input field. It is used to identify the
	// input field when sending events to the eventbus.
	InputName string
	// OnListAction is called when the user presses a key that should trigger
	// an action on the contextual list.
	OnListAction func(listaction.ListAction)
	// OnDone is called when the user is done entering input.
	OnDone func(input.DoneSource)
	// OnInsertFromContext is called when the user presses a key that should
	// insert the currently selected item from the contextual list.
	OnInsertFromContext func()
}

// ConnectContextualList connects the input field to a contextual list.
func (i *InputField) LinkContextualList(eventbus eventbusmod.IEventBus, options ContextualListLinkOpts) {
	// Handle input and dispatches to proper handlers
	i.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyCtrlN:
			options.OnListAction(listaction.NEXT)
			return nil
		case tcell.KeyUp, tcell.KeyCtrlP:
			options.OnListAction(listaction.PREV)
			return nil
		case tcell.KeyEnter:
			options.OnDone(input.Context)
			return nil
		case tcell.KeyCtrlJ:
			options.OnDone(input.Input)
			return nil
		case tcell.KeyTab:
			options.OnInsertFromContext()
			return nil
		}
		return event
	})

	// Subscribes to eventbus
	err := eventbus.Subscribe(eventbusmod.Subscription{
		Topic: "input." + options.InputName + ".settext",
		Handler: func(e eventbusmod.Event) {
			text, ok := e.Data.(string)
			if !ok {
				return
			}
			i.SetText(text)
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("could not subscribe to eventbus")
	}
}
