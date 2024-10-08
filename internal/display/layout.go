package display

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/controller"
	contextmod "github.com/vitorqb/addledger/internal/display/context"
	"github.com/vitorqb/addledger/internal/display/statement"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/state"
)

//go:generate $MOCKGEN --source=layout.go --destination=../../mocks/display/layout_mock.go

type StatementControllerAdapter struct{ controller.IInputController }

var _ statement.Controller = &StatementControllerAdapter{}

func (s *StatementControllerAdapter) HideModal()   { s.OnHideStatementModal() }
func (s *StatementControllerAdapter) LoadRequest() { s.OnLoadStatementRequest() }
func (s *StatementControllerAdapter) DiscardStatementEntry(index int) {
	s.OnDiscardStatementEntry(index)
}

type (
	// MainView represents the main view of the application, which contains
	// 4 main parts:
	// - View: the main view, which displays the current journal entry
	// - Input: the input field, where the user can type new entries
	// - Context: a context box, which displays information about the current
	//   input and helps the user to fill it.
	// - MessageBox: a box that displays messages to the user.
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

	// MessageBox is a tview.TextView that displays messages to the user.
	MessageBox struct {
		*tview.TextView
	}

	// tview App abstraction with only the methods we need.
	TviewApp interface {
		QueueUpdateDraw(func()) *tview.Application
		SetFocus(p tview.Primitive) *tview.Application
	}
)

type LayoutPage string

var (
	MainPage               LayoutPage = "main"
	ShortcutModalPage      LayoutPage = "shortcutModal"
	StatementModalPage     LayoutPage = "statementModal"
	LoadStatementModalPage LayoutPage = "loadStatementModal"
)

const (
	modalWith   = 50
	modalHeight = 10
)

func NewMainView(
	view *View,
	input *Input,
	context *Context,
	statementDisplay *StatementDisplay,
	messageBox *MessageBox,
	state *state.State,
) *MainView {
	mainView := &MainView{
		Flex:             tview.NewFlex(),
		view:             view,
		input:            input,
		context:          context,
		statementDisplay: statementDisplay,
		state:            state,
	}
	mainView.SetDirection(tview.FlexRow)
	mainView.AddItem(view, 0, 3, false)
	mainView.AddItem(statementDisplay, 0, 0, false)
	mainView.AddItem(input, 0, 1, false)
	mainView.AddItem(context, 0, 6, false)
	mainView.AddItem(messageBox, 0, 1, false)
	state.AddOnChangeHook(mainView.Refresh)
	return mainView
}

func (m *MainView) Focus(delegate func(p tview.Primitive)) {
	// Focusing on the MainView should always focus on the input field.
	delegate(m.input)
}

func (m *MainView) InputHasFocus() bool {
	return m.input.HasFocus()
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
	app TviewApp,
) (*Layout, error) {
	view := NewView(state)
	input := NewInput(controller, state, eventBus)
	messageBox := NewMessageBox(state)

	// Creates a Context
	accountList, err := NewAccountList(state, eventBus)
	if err != nil {
		return nil, fmt.Errorf("failed to create account list: %w", err)
	}
	descriptionPicker, err := contextmod.NewDescriptionPicker(state, eventBus, app)
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

	mainView := NewMainView(view, input, context, statementDisplay, messageBox, state)
	shortcutModal := center(NewShortcutModal(controller), modalWith, modalHeight)
	loadStatementModal := center(NewLoadStatementModal(controller, state.Display.StatementModal), modalWith*2, modalHeight*2)
	statementModal := center(statement.NewModal(&StatementControllerAdapter{controller}, state), modalWith*3, modalHeight*3)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage(string(MainPage), mainView, true)
	pages.AddPage(string(ShortcutModalPage), shortcutModal, true, false)
	pages.AddPage(string(LoadStatementModalPage), loadStatementModal, true, false)
	pages.AddPage(string(StatementModalPage), statementModal, true, false)
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return HandleGlobalShortcuts(controller, event)
	})

	layout := &Layout{
		Pages:    pages,
		mainView: mainView,
		state:    state,
		setFocus: app.SetFocus,
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
	l.refreshStatementModalDisplay()
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

func (l *Layout) refreshStatementModalDisplay() {
	if l.state.Display.StatementModal.Visible() {
		l.ShowPage(string(StatementModalPage))
		return
	}
	l.HidePage(string(StatementModalPage))
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

func NewMessageBox(state *state.State) *MessageBox {
	messageBox := MessageBox{tview.NewTextView()}
	messageBox.SetBorder(true)
	state.AddOnChangeHook(func() {
		messageBox.SetText(state.Display.UserMessage())
	})
	return &messageBox
}
