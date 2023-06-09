package display

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	Layout struct {
		state   *state.State
		View    *View
		Input   *Input
		Context *Context
		flex    *tview.Flex
	}
)

func NewLayout(
	controller controller.IInputController,
	state *state.State,
	eventBus eventbus.IEventBus,
) (*Layout, error) {
	view := NewView(state)
	input := NewInput(controller, state, eventBus)
	context, err := NewContext(state, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to instatiate context: %w", err)
	}
	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.GetContent(), 0, 5, false).
		AddItem(input.GetContent(), 0, 1, false).
		AddItem(context.GetContent(), 0, 10, false)
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyCtrlZ:
			controller.OnUndo()
			return nil
		}
		return event
	})
	return &Layout{
		state:   state,
		View:    view,
		Input:   input,
		Context: context,
		flex:    flex,
	}, nil
}

func (l *Layout) GetContent() tview.Primitive {
	return l.flex
}
