package display

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/display/input"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	Input    struct {
		controller          *controller.InputController
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

func NewInput(controller *controller.InputController, state *statemod.State) *Input {
	dateField := dateField(controller)
	descriptionField := descriptionField(controller)
	postingAccountField := input.NewPostingAccount(controller)
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
		i.pages.SwitchToPage(string(INPUT_DATE))
	case statemod.InputDescription:
		i.pages.SwitchToPage(string(INPUT_DESCRIPTION))
	case statemod.InputPostingAccount:
		i.postingAccountField.SetText("")
		i.pages.SwitchToPage(string(INPUT_POSTING_ACCOUNT))
	case statemod.InputPostingValue:
		i.postingValueField.SetText("")
		i.pages.SwitchToPage(string(INPUT_POSTING_VALUE))
	case statemod.Confirmation:
		i.pages.SwitchToPage(string(INPUT_CONFIRMATION))
	default:
	}
}

func descriptionField(controller *controller.InputController) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		text := inputField.GetText()
		controller.OnDescriptionInput(text)
	})
	return inputField
}

func dateField(controller *controller.InputController) *tview.InputField {
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

func postingValueField(controller *controller.InputController) *tview.InputField {
	valueInputField := tview.NewInputField()
	valueInputField.SetLabel("Value: ")
	valueInputField.SetDoneFunc(func(_ tcell.Key) {
		text := valueInputField.GetText()
		controller.OnPostingValueInput(text)
	})
	return valueInputField
}

func inputConfirmationField(controller *controller.InputController) *tview.TextView {
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

func (i *Input) GetContent() tview.Primitive {
	return i.pages
}
