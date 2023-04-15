package display

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	Input    struct {
		state *statemod.State
		pages *tview.Pages
	}
)

// !!! TODO Unify with state.Phase
const (
	INPUT_DATE            PageName = "INPUT_DATE"
	INPUT_DESCRIPTION     PageName = "INPUT_DESCRIPTION"
	INPUT_POSTING_ACCOUNT PageName = "INPUT_POSTING_ACCOUNT"
	INPUT_POSTING_VALUE   PageName = "INPUT_POSTING_VALUE"
)

func NewInput(state *statemod.State) *Input {
	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), dateField(state), true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), descriptionField(state), true, false)
	pages.AddPage(string(INPUT_POSTING_ACCOUNT), postingAccountField(state), true, false)
	pages.AddPage(string(INPUT_POSTING_VALUE), postingValueField(state), true, false)

	inputBox := &Input{state, pages}
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
		i.pages.SwitchToPage(string(INPUT_POSTING_ACCOUNT))
	case statemod.InputPostingValue:
		i.pages.SwitchToPage(string(INPUT_POSTING_VALUE))
	default:
	}
}

func descriptionField(state *statemod.State) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		state.JournalEntryInput.SetDescription(inputField.GetText())
		state.NextPhase()
	})
	return inputField
}

func dateField(state *statemod.State) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Date: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		text := inputField.GetText()
		date, err := time.Parse("2006-01-02", text)
		if err != nil {
			return
		}
		state.JournalEntryInput.SetDate(date)
		state.NextPhase()
	})
	return inputField
}

func postingAccountField(state *statemod.State) *tview.InputField {
	accountInputField := tview.NewInputField()
	accountInputField.SetLabel("Account: ")
	accountInputField.SetDoneFunc(func(key tcell.Key) {
		text := accountInputField.GetText()
		posting := state.JournalEntryInput.CurrentPosting()
		posting.SetAccount(text)
		state.NextPhase()
	})
	return accountInputField
}

func postingValueField(state *statemod.State) *tview.InputField {
	valueInputField := tview.NewInputField()
	valueInputField.SetLabel("Value: ")
	valueInputField.SetDoneFunc(func(key tcell.Key) {
		text := valueInputField.GetText()
		posting := state.JournalEntryInput.CurrentPosting()
		posting.SetValue(text)
		state.JournalEntryInput.AdvancePosting()
		state.SetPhase(statemod.InputPostingAccount)
	})
	return valueInputField
}

func (i *Input) GetContent() tview.Primitive {
	return i.pages
}
