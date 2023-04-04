package inputcontextbox

import "github.com/rivo/tview"

type (
	InputContextBox struct{}
)

func NewInputContextBox() *InputContextBox {
	return &InputContextBox{}
}

func (i InputContextBox) GetTextView() *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	return textView
}
