package inputbox

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	PostingInputBox struct {
		state        *state.State
		postingIndex int
		pages        *tview.Pages
	}

	PostingPageName string
)

const (
	INPUT_POSTING_ACCOUNT PostingPageName = "INPUT_POSTING_ACCOUNT"
	INPUT_POSTING_VALUE   PostingPageName = "INPUT_POSTING_VALUE"
)

func NewPostingInputBox(s *state.State) *PostingInputBox {
	pages := tview.NewPages()
	postingInputBox := &PostingInputBox{s, 0, pages}
	pages.AddPage(string(INPUT_POSTING_ACCOUNT), postingInputBox.getAccountField(), true, false)
	pages.SwitchToPage(string(INPUT_POSTING_ACCOUNT))
	return postingInputBox
}

func (i *PostingInputBox) getAccountField() *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Account: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		account := inputField.GetText()
		posting, exists := i.state.JournalEntryInput.GetPosting(0)
		if !exists {
			i.state.JournalEntryInput.AddPosting()
			posting, _ = i.state.JournalEntryInput.GetPosting(0)
		}
		posting.SetAccount(account)
		i.state.NextPhase()
	})
	return inputField
}

func (p *PostingInputBox) GetContent() tview.Primitive {
	return p.pages
}
