package display

import "github.com/rivo/tview"

type (
	Context struct{}
)

func NewContext() *Context {
	return &Context{}
}

func (c Context) GetContent() *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	return textView
}
