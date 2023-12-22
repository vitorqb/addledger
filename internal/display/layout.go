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
	// MainView represents the main view of the application, which contains
	// 3 parts:
	// - View: the main view, which displays the current journal entry
	// - Input: the input field, where the user can type new entries
	// - Context: a context box, which displays information about the current
	//   input and helps the user to fill it.
	// The MainView may also contain a StatementDisplay, which displays
	// information about the current statement (if any).
	MainView struct {
		*tview.Flex
		view             *View
		input            *Input
		context          *Context
		statementDisplay *StatementDisplay
		state            *state.State
	}

	// Layout represents the layout visible to the user.
	// The layout inherits from tview.Pages, which allows to switch between
	// the main view of the application (containing View, Input and Context)
	// and other modal views (such as the ShortcutModal).
	Layout struct {
		*tview.Pages
		mainView *MainView
		state    *state.State
		setFocus func(p tview.Primitive) *tview.Application
	}
)

type LayoutPage string

var (
	MainPage               LayoutPage = "main"
	ShortcutModalPage      LayoutPage = "shortcutModal"
	LoadStatementModalPage LayoutPage = "loadStatementModal"
)

const (
	modalWith   = 50
	modalHeight = 10
)

func NewMainView(view *View, input *Input, context *Context, statementDisplay *StatementDisplay, state *state.State) *MainView {
	mainView := &MainView{
		Flex:             tview.NewFlex(),
		view:             view,
		input:            input,
		context:          context,
		statementDisplay: statementDisplay,
		state:            state,
	}
	mainView.SetDirection(tview.FlexRow)
	mainView.AddItem(view.GetContent(), 0, 5, false)
	mainView.AddItem(statementDisplay, 0, 0, false)
	mainView.AddItem(input.GetContent(), 0, 1, false)
	mainView.AddItem(context.GetContent(), 0, 10, false)
	state.AddOnChangeHook(mainView.Refresh)
	return mainView
}

func (m *MainView) Focus(delegate func(p tview.Primitive)) {
	// Focusing on the MainView should always focus on the input field.
	delegate(m.input.GetContent())
}

func (m *MainView) InputHasFocus() bool {
	return m.input.GetContent().HasFocus()
}

func (m *MainView) RefreshStatementDisplay() {
	if len(m.state.StatementEntries) > 0 {
		m.ResizeItem(m.statementDisplay, 0, 1)
		return
	}
	m.ResizeItem(m.statementDisplay, 0, 0)
}

func (m *MainView) Refresh() {
	m.RefreshStatementDisplay()
}

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

	mainView := NewMainView(view, input, context, statementDisplay, state)
	shortcutModal := center(NewShortcutModal(controller), modalWith, modalHeight)
	loadStatementModal := center(NewLoadStatementModal(controller), modalWith*2, modalHeight*2)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage(string(MainPage), mainView, true)
	pages.AddPage(string(ShortcutModalPage), shortcutModal, true, false)
	pages.AddPage(string(LoadStatementModalPage), loadStatementModal, true, false)
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return HandleGlobalShortcuts(controller, event)
	})

	layout := &Layout{
		Pages:    pages,
		mainView: mainView,
		state:    state,
		setFocus: setFocus,
	}
	layout.Refresh()
	state.AddOnChangeHook(layout.Refresh)
	return layout, nil
}

func (l *Layout) InputHasFocus() bool {
	return l.mainView.InputHasFocus()
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
	l.refreshShortcutModalDisplay()
	l.refreshLoadStatementModalDisplay()
	// Make sure that once the modal is hidden we focus back on the main view.
	if frontPage, _ := l.GetFrontPage(); frontPage == string(MainPage) {
		l.setFocus(l.mainView)
	}
}

func (l *Layout) refreshShortcutModalDisplay() {
	if l.state.Display.ShortcutModal() {
		l.ShowPage(string(ShortcutModalPage))
		return
	}
	l.HidePage(string(ShortcutModalPage))
}

func (l *Layout) refreshLoadStatementModalDisplay() {
	if l.state.Display.LoadStatementModal() {
		l.ShowPage(string(LoadStatementModalPage))
		return
	}
	l.HidePage(string(LoadStatementModalPage))
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
