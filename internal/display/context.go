package display

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type Context struct {
	state *statemod.State
	pages *tview.Pages
}

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

// AccountList represents a list of accounts
type AccountList struct {
	*tview.List
}

// HandleRequest handles an AccountListRequest
func (al *AccountList) HandleRequest(req string) {
	switch strings.ToLower(req) {
	case "next":
		event := tcell.NewEventKey(tcell.KeyDown, tcell.RuneDArrow, tcell.ModNone)
		al.InputHandler()(event, func(p tview.Primitive) {})
	}
}

func accountList(state *statemod.State) *AccountList {
	list := &AccountList{tview.NewList()}
	for _, acc := range state.GetAccounts() {
		list.AddItem(acc, "", 0, nil)
	}
	list.ShowSecondaryText(false)
	return list
}
