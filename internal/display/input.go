package display

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display/input"
	"github.com/vitorqb/addledger/internal/eventbus"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	Input    struct {
		controller          controller.IInputController
		state               *statemod.State
		pages               *tview.Pages
		dateField           *tview.InputField
		descriptionField    *tview.InputField
		postingAccountField *input.PostingAccountField
		postingValueField   *tview.InputField
	}
)

// !!! TODO Unify with state.Phase
const (
	INPUT_DATE            PageName = "INPUT_DATE"
	INPUT_DESCRIPTION     PageName = "INPUT_DESCRIPTION"
	INPUT_POSTING_ACCOUNT PageName = "INPUT_POSTING_ACCOUNT"
	INPUT_POSTING_VALUE   PageName = "INPUT_POSTING_VALUE"
	INPUT_CONFIRMATION    PageName = "INPUT_CONFIRMATION"
)

func NewInput(
	controller controller.IInputController,
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) *Input {
	dateField := dateField(controller)
	descriptionField := DescriptionField(controller, eventbus)
	postingAccountField := input.NewPostingAccount(controller, eventbus)
	postingValueField := postingValueField(controller)
	inputConfirmationField := inputConfirmationField(controller)

	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), dateField, true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), descriptionField, true, false)
	pages.AddPage(string(INPUT_POSTING_ACCOUNT), postingAccountField, true, false)
	pages.AddPage(string(INPUT_POSTING_VALUE), postingValueField, true, false)
	pages.AddPage(string(INPUT_CONFIRMATION), inputConfirmationField, true, false)

	inputBox := &Input{
		controller:          controller,
		state:               state,
		pages:               pages,
		dateField:           dateField,
		postingValueField:   postingValueField,
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
			i.pages.SwitchToPage(string(INPUT_DATE))
		}
	case statemod.InputDescription:
		if i.CurrentPageName() != string(INPUT_DESCRIPTION) {
			i.pages.SwitchToPage(string(INPUT_DESCRIPTION))
		}
	case statemod.InputPostingAccount:
		if i.CurrentPageName() != string(INPUT_POSTING_ACCOUNT) {
			i.postingAccountField.SetText("")
			i.pages.SwitchToPage(string(INPUT_POSTING_ACCOUNT))
		}
	case statemod.InputPostingValue:
		if i.CurrentPageName() != string(INPUT_POSTING_VALUE) {
			i.postingValueField.SetText("")
			i.pages.SwitchToPage(string(INPUT_POSTING_VALUE))
		}
	case statemod.Confirmation:
		if i.CurrentPageName() != string(INPUT_CONFIRMATION) {
			i.pages.SwitchToPage(string(INPUT_CONFIRMATION))
		}
	default:
	}
}

func DescriptionField(
	controller controller.IInputController,
	eventbus eventbus.IEventBus,
) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetChangedFunc(controller.OnDescriptionChanged)
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyDown, tcell.KeyCtrlN:
			controller.OnDescriptionListAction(listaction.NEXT)
			return nil
		// Arrow up -> move contextual list up
		case tcell.KeyUp, tcell.KeyCtrlP:
			controller.OnDescriptionListAction(listaction.PREV)
			return nil
		// If user hit enter...
		case tcell.KeyEnter:
			controller.OnDescriptionSelectedFromContext()
			return nil
		// if Ctrl+J, use input as it is
		case tcell.KeyCtrlJ:
			controller.OnDescriptionDone()
			return nil
		// if Tab then autocompletes
		case tcell.KeyTab:
			controller.OnDescriptionInsertFromContext()
			return nil
		}
		return event
	})
	err := eventbus.Subscribe(eventbusmod.Subscription{
		Topic: "input.description.settext",
		Handler: func(e eventbusmod.Event) {
			text, ok := e.Data.(string)
			if !ok {
				logrus.WithField("event", e).Error("Received invalid event")
				return
			}
			inputField.SetText(text)
		},
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to subscribe to Topic")
	}
	return inputField
}

func dateField(controller controller.IInputController) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Date: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		text := inputField.GetText()
		date, err := time.Parse("2006-01-02", text)
		if err != nil {
			return
		}
		controller.OnDateInput(date)
	})
	return inputField
}

func postingValueField(controller controller.IInputController) *tview.InputField {
	valueInputField := tview.NewInputField()
	valueInputField.SetLabel("Value: ")
	valueInputField.SetDoneFunc(func(_ tcell.Key) {
		text := valueInputField.GetText()
		controller.OnPostingValueInput(text)
	})
	return valueInputField
}

func inputConfirmationField(controller controller.IInputController) *tview.TextView {
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
	out, _ := i.pages.GetFrontPage()
	return out
}

func (i *Input) GetContent() tview.Primitive {
	return i.pages
}
