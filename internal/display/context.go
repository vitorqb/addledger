package display

import (
	"github.com/rivo/tview"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type (
	Context struct {
		state *statemod.State
		pages *tview.Pages
	}
)

func NewContext(state *statemod.State) *Context {
	context := new(Context)
	context.state = state
	context.pages = tview.NewPages()
	context.pages.SetBorder(true)
	context.pages.AddPage("accountList", accountList(state), true, false)
	context.pages.AddPage("empty", tview.NewBox(), true, false)
	context.pages.SwitchToPage("accountList")
	context.Refresh()
	state.AddOnChangeHook(context.Refresh)
	return context
}

func (c Context) GetContent() *tview.Pages { return c.pages }

func (c Context) Refresh() {
	switch c.state.CurrentPhase() {
	case statemod.InputPostingAccount:
		c.pages.SwitchToPage("accountList")
	default:
		c.pages.SwitchToPage("empty")
	}
}

func accountList(state *statemod.State) *tview.List {
	list := tview.NewList()
	for _, acc := range state.GetAccounts() {
		list.AddItem(acc, "", 0, nil)
	}
	list.ShowSecondaryText(false)
	return list
}
