package display

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	controllermod "github.com/vitorqb/addledger/internal/controller"
	display_input "github.com/vitorqb/addledger/internal/display/input"
	"github.com/vitorqb/addledger/internal/display/widgets"
	"github.com/vitorqb/addledger/internal/eventbus"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	Input    struct {
		*tview.Pages
		controller          controllermod.IInputController
		state               *statemod.State
		dateField           *tview.InputField
		descriptionField    *widgets.InputField
		postingAccountField *widgets.InputField
		postingAmmountField *tview.InputField
	}
)

// !!! TODO Unify with state.Phase
const (
	INPUT_DATE            PageName = "INPUT_DATE"
	INPUT_DESCRIPTION     PageName = "INPUT_DESCRIPTION"
	INPUT_TAGS            PageName = "INPUT_TAGS"
	INPUT_POSTING_ACCOUNT PageName = "INPUT_POSTING_ACCOUNT"
	INPUT_POSTING_AMMOUNT PageName = "INPUT_POSTING_AMMOUNT"
	INPUT_CONFIRMATION    PageName = "INPUT_CONFIRMATION"
)

func NewInput(
	controller controllermod.IInputController,
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) *Input {
	dateField := DateField(controller)
	descriptionField := DescriptionField(controller, eventbus)
	tagsField := NewTagsField(controller, eventbus)
	postingAccountField := display_input.NewPostingAccount(controller, eventbus)
	postingAmmountField := PostingAmmountField(controller)
	inputConfirmationField := inputConfirmationField(controller)

	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), dateField, true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), descriptionField, true, false)
	pages.AddPage(string(INPUT_POSTING_ACCOUNT), postingAccountField, true, false)
	pages.AddPage(string(INPUT_POSTING_AMMOUNT), postingAmmountField, true, false)
	pages.AddPage(string(INPUT_CONFIRMATION), inputConfirmationField, true, false)
	pages.AddPage(string(INPUT_TAGS), tagsField, true, false)

	inputBox := &Input{
		Pages:               pages,
		controller:          controller,
		state:               state,
		dateField:           dateField,
		postingAmmountField: postingAmmountField,
		descriptionField:    descriptionField,
		postingAccountField: postingAccountField,
	}
	inputBox.refresh()

	state.AddOnChangeHook(inputBox.refresh)

	return inputBox
}

func (i *Input) refresh() {
	switch i.state.CurrentPhase() {
	case statemod.InputDate:
		if i.CurrentPageName() != string(INPUT_DATE) {
			if i.dateField.GetText() != "" {
				i.dateField.SetText("")
			}
			i.SwitchToPage(string(INPUT_DATE))
			i.controller.OnDateChanged("")
		}
	case statemod.InputDescription:
		if i.CurrentPageName() != string(INPUT_DESCRIPTION) {
			if i.descriptionField.GetText() != "" {
				i.descriptionField.SetText("")
			}
			i.SwitchToPage(string(INPUT_DESCRIPTION))
		}
	case statemod.InputTags:
		if i.CurrentPageName() != string(INPUT_TAGS) {
			i.SwitchToPage(string(INPUT_TAGS))
		}
	case statemod.InputPostingAccount:
		if i.CurrentPageName() != string(INPUT_POSTING_ACCOUNT) {
			i.postingAccountField.SetText("")
			i.SwitchToPage(string(INPUT_POSTING_ACCOUNT))
		}
	case statemod.InputPostingAmmount:
		if i.CurrentPageName() != string(INPUT_POSTING_AMMOUNT) {
			i.postingAmmountField.SetText("")
			i.SwitchToPage(string(INPUT_POSTING_AMMOUNT))
		}
	case statemod.Confirmation:
		if i.CurrentPageName() != string(INPUT_CONFIRMATION) {
			i.SwitchToPage(string(INPUT_CONFIRMATION))
		}
	default:
	}
}

func DescriptionField(
	controller controllermod.IInputController,
	eventbus eventbus.IEventBus,
) *widgets.InputField {
	inputField := widgets.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetChangedFunc(controller.OnDescriptionChanged)
	inputField.LinkContextualList(eventbus, widgets.ContextualListLinkOpts{
		InputName:           "description",
		OnListAction:        controller.OnDescriptionListAction,
		OnDone:              controller.OnDescriptionDone,
		OnInsertFromContext: controller.OnDescriptionInsertFromContext,
	})
	return inputField
}

func DateField(controller controllermod.IInputController) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Date: ")
	inputField.SetChangedFunc(controller.OnDateChanged)
	inputField.SetText("")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		controller.OnDateDone()
	})
	return inputField
}

func PostingAmmountField(controller controllermod.IInputController) *tview.InputField {
	postingAmmountField := tview.NewInputField()
	postingAmmountField.SetLabel("Ammount: ")
	postingAmmountField.SetChangedFunc(func(text string) {
		controller.OnPostingAmmountChanged(text)
	})
	postingAmmountField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			controller.OnPostingAmmountDone(input.Context)
			return nil
		case tcell.KeyCtrlJ:
			controller.OnPostingAmmountDone(input.Input)
			return nil
		}
		return event
	})
	return postingAmmountField
}

func inputConfirmationField(controller controllermod.IInputController) *tview.TextView {
	field := tview.NewTextView()
	field.SetText("Do you want to commit the transaction? [Y/n]")
	field.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			controller.OnInputConfirmation()
		}
		switch event.Rune() {
		case 'y':
			controller.OnInputConfirmation()
		case 'Y':
			controller.OnInputConfirmation()
		case 'n':
			controller.OnInputRejection()
		case 'N':
			controller.OnInputRejection()
		}
		return event
	})
	return field
}

func (i *Input) CurrentPageName() string {
	out, _ := i.GetFrontPage()
	return out
}

func NewTagsField(
	controller controllermod.IInputController,
	eventbus eventbus.IEventBus,
) *widgets.InputField {
	field := widgets.NewInputField()
	field.SetLabel("Tags: ")
	field.SetChangedFunc(controller.OnTagChanged)
	field.LinkContextualList(eventbus, widgets.ContextualListLinkOpts{
		InputName:           "tag",
		OnListAction:        controller.OnTagListAction,
		OnDone:              controller.OnTagDone,
		OnInsertFromContext: controller.OnTagInsertFromContext,
	})
	return field
}
