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
		setDate          func(time.Time)
		setDescription   func(string)
		state            *state.State
		pages            *tview.Pages
	}
)

const (
	INPUT_DATE PageName ="INPUT_DATE"
	INPUT_DESCRIPTION PageName = "INPUT_DESCRIPTION"
)

func NewInputBox(SetDate func(time.Time), SetDescription func(string), state *state.State) *InputBox {
	pages := tview.NewPages()
	inputBox := &InputBox{SetDate, SetDescription, state, pages}
	pages.SetBorder(true)
	pages.AddPage(string(INPUT_DATE), inputBox.getDateInputField(), true, false)
	pages.AddPage(string(INPUT_DESCRIPTION), inputBox.getDescriptionInputField(), true, false)
	inputBox.Refresh()
	return inputBox
}

func (i *InputBox) Refresh() {
	switch i.state.CurrentPhase {
	case state.Date:
		i.pages.SwitchToPage(string(INPUT_DATE))
	case state.Description:
		i.pages.SwitchToPage(string(INPUT_DESCRIPTION))
	default:
	}	
}

func (i *InputBox) getDescriptionInputField() *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Description: ")
	inputField.SetChangedFunc(func(x string) {
		i.setDescription(x)
		i.Refresh()
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
		i.setDate(date)
		i.Refresh()
	})
	return inputField
}

func (i *InputBox) GetContent() tview.Primitive {
	return i.pages
}
