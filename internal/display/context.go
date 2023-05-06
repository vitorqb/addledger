package display

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type Context struct {
	state *statemod.State
	pages *tview.Pages
}

func NewContext(
	state *statemod.State,
	eventBus eventbus.IEventBus,
) (*Context, error) {
	context := new(Context)
	accountList := accountList(state)
	context.state = state
	context.pages = tview.NewPages()
	context.pages.SetBorder(true)
	context.pages.AddPage("accountList", accountList, true, false)
	context.pages.AddPage("empty", tview.NewBox(), true, false)
	context.pages.SwitchToPage("accountList")
	context.Refresh()
	state.AddOnChangeHook(context.Refresh)
	err := eventBus.Subscribe(eventbus.Subscription{
		Topic: "input.postingaccount.listaction",
		Handler: func(e eventbus.Event) {
			listAction, ok := e.Data.(listaction.ListAction)
			if !ok {
				logrus.Errorf("received event w/ unexpected data %+v", e)
			}
			accountList.handleAction(listAction)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return context, nil
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

func accountList(state *statemod.State) *AccountList {
	list := &AccountList{tview.NewList()}
	for _, acc := range state.GetAccounts() {
		list.AddItem(acc, "", 0, nil)
	}
	list.ShowSecondaryText(false)
	return list
}

func (al *AccountList) handleAction(action listaction.ListAction) {
	switch action {
	case listaction.NEXT:
		eventKey := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
		al.InputHandler()(eventKey, func(p tview.Primitive) {})
	case listaction.PREV:
		eventKey := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
		al.InputHandler()(eventKey, func(p tview.Primitive) {})
	case listaction.NONE:
	default:
	}
}
