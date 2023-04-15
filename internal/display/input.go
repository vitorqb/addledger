package display

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	Input    struct {
		state *state.State
		pages *tview.Pages
	}
)

// !!! TODO Unify with state.Phase
const (
	INPUT_DATE        PageName = "INPUT_DATE"
	INPUT_DESCRIPTION PageName = "INPUT_DESCRIPTION"
	// !!!! TODO Change to `INPUT_POSTING_ACCOUNT` AND `INPUT_POSTING_VALUE`
	// !!!! we don't need to have a separated pages for posting
	INPUT_POSTINGS PageName = "INPUT_POSTINGS"
)

func NewInput(state *state.State) *Input {
	pages := tview.NewPages()
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), dateField(state), true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), descriptionField(state), true, false)
	pages.AddPage(string(INPUT_POSTINGS), postingPages(state), true, false)

	inputBox := &Input{state, pages}
	inputBox.refresh()

	state.AddOnChangeHook(inputBox.refresh)

	return inputBox
}

func (i *Input) refresh() {
	switch i.state.CurrentPhase {
	case state.InputDate:
		i.pages.SwitchToPage(string(INPUT_DATE))
	case state.InputDescription:
		i.pages.SwitchToPage(string(INPUT_DESCRIPTION))
	case state.InputPostings:
		i.pages.SwitchToPage(string(INPUT_POSTINGS))
	default:
	}
}

func descriptionField(state *state.State) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		state.JournalEntryInput.SetDescription(inputField.GetText())
		state.NextPhase()
	})
	return inputField
}

func dateField(state *state.State) *tview.InputField {
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

func postingPages(state *state.State) *tview.Pages {
	pages := tview.NewPages()
	currentIndex := 0

	accountInputField := tview.NewInputField()
	accountInputField.SetLabel("Account: ")
	accountInputField.SetDoneFunc(func(key tcell.Key) {
		text := accountInputField.GetText()
		posting, found := state.JournalEntryInput.GetPosting(currentIndex)
		if !found {
			posting = state.JournalEntryInput.AddPosting()
		}
		posting.SetAccount(text)
		pages.SwitchToPage("value")
	})

	valueInputField := tview.NewInputField()
	valueInputField.SetLabel("Value: ")
	valueInputField.SetDoneFunc(func(key tcell.Key) {
		text := valueInputField.GetText()
		posting, found := state.JournalEntryInput.GetPosting(currentIndex)
		if !found {
			posting = state.JournalEntryInput.AddPosting()
		}
		posting.SetValue(text)
		currentIndex++
		pages.SwitchToPage("account")
	})

	pages.AddAndSwitchToPage("account", accountInputField, true)
	pages.AddPage("value", valueInputField, true, false)

	return pages
}

func (i *Input) GetContent() tview.Primitive {
	return i.pages
}
