package inputbox

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	PageName string
	InputBox struct {
		state *state.State
		pages *tview.Pages
	}
)

// !!! TODO Unify with state.Phase
const (
	INPUT_DATE        PageName = "INPUT_DATE"
	INPUT_DESCRIPTION PageName = "INPUT_DESCRIPTION"
	INPUT_POSTINGS    PageName = "INPUT_POSTINGS"
)

func NewInputBox(state *state.State) *InputBox {
	pages := tview.NewPages()
	inputBox := &InputBox{state, pages}
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), inputBox.getDateInputField(), true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), inputBox.getDescriptionInputField(), true, false)
	pages.AddPage(string(INPUT_POSTINGS), inputBox.getPostingInputPages(), true, false)
	inputBox.refresh()
	state.AddOnChangeHook(inputBox.refresh)
	return inputBox
}

func (i *InputBox) refresh() {
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

func (i *InputBox) getDescriptionInputField() *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		i.state.JournalEntryInput.SetDescription(inputField.GetText())
		i.state.NextPhase()
	})
	return inputField
}

func (i *InputBox) getDateInputField() *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Date: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		text := inputField.GetText()
		date, err := time.Parse("2006-01-02", text)
		if err != nil {
			return
		}
		i.state.JournalEntryInput.SetDate(date)
		i.state.NextPhase()
	})
	return inputField
}

func (i *InputBox) getPostingInputPages() *tview.Pages {
	pages := tview.NewPages()

	accountInputField := tview.NewInputField()
	accountInputField.SetLabel("Account: ")
	accountInputField.SetDoneFunc(func(key tcell.Key) {
		text := accountInputField.GetText()
		if _, found := i.state.JournalEntryInput.GetPosting(0); ! found {
			i.state.JournalEntryInput.AddPosting()
		}
		posting, _ := i.state.JournalEntryInput.GetPosting(0)
		posting.SetAccount(text)
	})
	pages.AddAndSwitchToPage("1", accountInputField, true)

	return pages
}

func (i *InputBox) GetContent() tview.Primitive {
	return i.pages
}
