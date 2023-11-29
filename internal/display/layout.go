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
		*tview.Pages
		mainView         *tview.Flex
		state            *state.State
		View             *View
		Input            *Input
		Context          *Context
		statementDisplay *StatementDisplay
		setFocus         func(p tview.Primitive) *tview.Application
	}
)

const (
	modalWith   = 50
	modalHeight = 10
)

func NewLayout(
	controller controller.IInputController,
	state *state.State,
	eventBus eventbus.IEventBus,
	setFocus func(p tview.Primitive) *tview.Application,
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

	mainView := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.GetContent(), 0, 5, false).
		AddItem(statementDisplay, 0, 0, false).
		AddItem(input.GetContent(), 0, 1, false).
		AddItem(context.GetContent(), 0, 10, false)

	shortcutModal := center(NewShortcutModal(controller), modalWith, modalHeight)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage("main", mainView, true)
	pages.AddPage("shortcutModal", shortcutModal, true, false)
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return HandleGlobalShortcuts(controller, event)
	})

	layout := &Layout{
		Pages:            pages,
		mainView:         mainView,
		state:            state,
		View:             view,
		Input:            input,
		Context:          context,
		statementDisplay: statementDisplay,
		setFocus:         setFocus,
	}
	layout.Refresh()
	state.AddOnChangeHook(layout.Refresh)
	return layout, nil
}

// Expose ResizeItem from Main View
func (l *Layout) ResizeItem(item tview.Primitive, width, height int) {
	l.mainView.ResizeItem(item, width, height)
}

// Expose GetItem from Main View
func (l *Layout) GetItem(index int) tview.Primitive {
	return l.mainView.GetItem(index)
}

func (l *Layout) Refresh() {
	l.refreshStatementDisplay()
	l.refreshShortcutModalDisplay()
}

func (l *Layout) refreshStatementDisplay() {
	if len(l.state.StatementEntries) > 0 {
		l.ResizeItem(l.statementDisplay, 0, 1)
		return
	}
	l.ResizeItem(l.statementDisplay, 0, 0)
}

func (l *Layout) refreshShortcutModalDisplay() {
	if l.state.ShortcutModalDisplayed {
		l.ShowPage("shortcutModal")
		return
	}
	l.HidePage("shortcutModal")
	// Make sure that once the modal is hidden we focus back on the input field.
	l.setFocus(l.Input.GetContent())
}

func HandleGlobalShortcuts(controller controller.IInputController, event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return nil
	}
	switch key := event.Key(); key {
	case tcell.KeyCtrlZ:
		controller.OnUndo()
		return nil
	case tcell.KeyCtrlQ:
		controller.OnDisplayShortcutModal()
		return nil
	}
	return event
}
