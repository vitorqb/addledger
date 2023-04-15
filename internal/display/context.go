package display

import "github.com/rivo/tview"

type (
	Context struct{}
)

func NewContext() *Context {
	return &Context{}
}

// !!!! TODO Rename to GetContent
func (c Context) GetTextView() *tview.TextView {
	textView := tview.NewTextView()
	textView.SetBorder(true)
	return textView
}
