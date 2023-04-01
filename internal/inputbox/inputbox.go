package inputbox

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/utils"
)

func NewInputBox(SetDate func(time.Time)) tview.Primitive {
	dateInputValue := ""
	dateInputField := tview.
		NewInputField().
		SetLabel("Date: ").
		SetFieldWidth(20).
		SetChangedFunc(func(text string) {
			dateInputValue = text
		}).
		SetDoneFunc(func(key tcell.Key) {
			date, err := time.Parse("2006-01-02", dateInputValue)
			if err != nil {
				panic(err)
			}

			SetDate(date)
		})
	return utils.Center(20, 1, dateInputField)
}
