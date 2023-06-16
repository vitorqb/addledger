package display

import (
	"fmt"

	"github.com/rivo/tview"
	contextmod "github.com/vitorqb/addledger/internal/display/context"
	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	statemod "github.com/vitorqb/addledger/internal/state"
)

type Context struct {
	state *statemod.State
	pages *tview.Pages
}

func NewContext(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) (*Context, error) {

	// Creates an AccountList widget
	accountList, err := newAccountList(state, eventbus)
	if err != nil {
		return nil, fmt.Errorf("failed to create account list: %w", err)
	}

	// Creates an DescriptionPicker widget
	descriptionPicker, err := contextmod.NewDescriptionPicker(state, eventbus)
	if err != nil {
		return nil, fmt.Errorf("failed to create description picker: %w", err)
	}

	// Creates a Ammount Guesser
	ammountGuesser, err := contextmod.NewAmmountGuesser(state)
	if err != nil {
		return nil, fmt.Errorf("failed to create ammount guesser: %w", err)
	}

	// Creates Context
	context := new(Context)
	context.state = state
	context.pages = tview.NewPages()
	context.pages.SetBorder(true)
	context.pages.AddPage("accountList", accountList, true, false)
	context.pages.AddPage("descriptionPicker", descriptionPicker, true, false)
	context.pages.AddPage("ammountGuesser", ammountGuesser, true, false)
	context.pages.AddPage("empty", tview.NewBox(), true, false)
	context.pages.SwitchToPage("accountList")
	context.Refresh()
	state.AddOnChangeHook(context.Refresh)

	return context, nil
}

func (c Context) GetContent() *tview.Pages { return c.pages }

func (c Context) Refresh() {
	switch c.state.CurrentPhase() {
	case statemod.InputPostingAccount:
		c.pages.SwitchToPage("accountList")
	case statemod.InputDescription:
		c.pages.SwitchToPage("descriptionPicker")
	case statemod.InputPostingAmmount:
		c.pages.SwitchToPage("ammountGuesser")
	default:
		c.pages.SwitchToPage("empty")
	}
}

func newAccountList(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) (*widgets.ContextualList, error) {
	list := widgets.NewContextualList(
		func() []string {
			return state.GetAccounts()
		},
		func(s string) {
			state.InputMetadata.SetSelectedPostingAccount(s)
		},
		func() string {
			return state.InputMetadata.PostingAccountText()
		},
	)
	state.AddOnChangeHook(func() { list.Refresh() })
	err := eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.postingaccount.listaction",
		Handler: list.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return list, nil
}
