package display

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/controller"
	contextmod "github.com/vitorqb/addledger/internal/display/context"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	Layout struct {
		*tview.Flex
		state            *state.State
		View             *View
		Input            *Input
		Context          *Context
		statementDisplay *StatementDisplay
	}
)

func NewLayout(
	controller controller.IInputController,
	state *state.State,
	eventBus eventbus.IEventBus,
) (*Layout, error) {
	view := NewView(state)
	input := NewInput(controller, state, eventBus)

	// Creates a Context
	accountList, err := NewAccountList(state, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create account list: %w", err)
	}
	descriptionPicker, err := contextmod.NewDescriptionPicker(state, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create description picker: %w", err)
	}
	ammountGuesser, err := contextmod.NewAmmountGuesser(state)
	if err != nil {
		return nil, fmt.Errorf("failed to create ammount guesser: %w", err)
	}
	dateGuesser, err := NewDateGuesser(state)
	if err != nil {
		return nil, fmt.Errorf("failed to create date guesser: %w", err)
	}
	tagsPicker, err := NewTagsPicker(state, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create tags picker: %w", err)
	}
	contextWidgets := []ContextWidget{
		{PageName: "accountList", Widget: accountList},
		{PageName: "descriptionPicker", Widget: descriptionPicker},
		{PageName: "ammountGuesser", Widget: ammountGuesser},
		{PageName: "dateGuesser", Widget: dateGuesser},
		{PageName: "tagsPicker", Widget: tagsPicker},
		{PageName: "empty", Widget: tview.NewBox()},
	}

	context, err := NewContext(state, contextWidgets)
	if err != nil {
		return nil, fmt.Errorf("failed to instatiate context: %w", err)
	}

	// Create a StatementDisplay, which displays info about the current statement
	// (if any). It starts hidden and is only shown when there are statements.
	statementDisplay := NewStatementDisplay(state)

	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.GetContent(), 0, 5, false).
		AddItem(statementDisplay, 0, 0, false).
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
	layout := &Layout{
		Flex:             flex,
		state:            state,
		View:             view,
		Input:            input,
		Context:          context,
		statementDisplay: statementDisplay,
	}
	layout.Refresh()
	state.AddOnChangeHook(layout.Refresh)
	return layout, nil
}

func (l *Layout) Refresh() {
	// If there are statements, display the statement display
	if len(l.state.StatementEntries) > 0 {
		l.Flex.ResizeItem(l.statementDisplay, 0, 1)
		return
	}
	l.Flex.ResizeItem(l.statementDisplay, 0, 0)
}
