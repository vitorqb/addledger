package inputbox

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	InputBox struct {
		dateInput        *tview.InputField
		descriptionInput *tview.InputField
	}
)

func newDateInputField(SetDate func(time.Time)) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Date: ")
	inputField.SetDoneFunc(func(_ tcell.Key) {
		text := inputField.GetText()
		date, err := time.Parse("2006-01-02", text)
		if err != nil {
			return
		}
		SetDate(date)
	})
	return inputField
}

func newDescriptionInputField(SetDescription func(string)) *tview.InputField {
	return tview.
		NewInputField().
		SetLabel("Description: ").
		SetChangedFunc(func(x string) {
			SetDescription(x)
		})
}

func NewInputBox(SetDate func(time.Time), SetDescription func(string)) InputBox {
	dateInputField := newDateInputField(SetDate)
	descriptionInputField := newDescriptionInputField(SetDescription)
	return InputBox{dateInputField, descriptionInputField}
}

func (i InputBox) GetInputField() *tview.InputField {
	i.dateInput.SetBorder(true)
	return i.dateInput
}
