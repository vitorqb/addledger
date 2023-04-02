package inputbox

import (
	"time"
	"github.com/rivo/tview"
)

func newDateInputField(SetDate func(time.Time)) *tview.InputField {
	return tview.
		NewInputField().
		SetLabel("Date: ").
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			date, err := time.Parse("2006-01-02", text)
			if err != nil {
				return
			}
			SetDate(date)			
		})
}

func newDescriptionInputField(SetDescription func(string)) *tview.InputField {
	return tview.
		NewInputField().
		SetLabel("Description: ").
		SetFieldWidth(100).
		SetChangedFunc(func(x string) {
			SetDescription(x)
		})
}

func NewInputBox(SetDate func(time.Time), SetDescription func(string)) tview.Primitive {
	dateInputField := newDateInputField(SetDate)
	descriptionInputField := newDescriptionInputField(SetDescription)
	form := tview.NewForm()
	form.SetBorder(true)
	form.AddFormItem(dateInputField)
	form.AddFormItem(descriptionInputField)
	return form
}
