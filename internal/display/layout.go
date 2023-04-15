package display

import (
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/state"
)

type (
	Layout struct {
		state   *state.State
		View    *View
		Input   *Input
		Context *Context
		flex    *tview.Flex
	}
)

func NewLayout(state *state.State) *Layout {
	view := NewView(state)
	input := NewInput(state)
	context := NewContext()
	flex := tview.
		NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(view.GetPrimitive(), 0, 5, false).
		AddItem(input.GetContent(), 0, 1, false).
		AddItem(context.GetTextView(), 0, 10, false)
	return &Layout{
		state:   state,
		View:    view,
		Input:   input,
		Context: context,
		flex:    flex,
	}
}

func (l *Layout) GetContent() tview.Primitive {
	return l.flex
}
