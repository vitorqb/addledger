package display

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/accountguesser"
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
	accountGuesser accountguesser.IAccountGuesser,
) (*Context, error) {
	// Creates an AccountList widget
	accountList, err := NewAccountList(state, eventbus, accountGuesser)
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

	// Creates a date guesser
	dateGuesser, err := NewDateGuesser(state)
	if err != nil {
		return nil, fmt.Errorf("failed to create date guesser: %w", err)
	}

	// Creates Context
	context := new(Context)
	context.state = state
	context.pages = tview.NewPages()
	context.pages.SetBorder(true)
	context.pages.AddPage("accountList", accountList, true, false)
	context.pages.AddPage("descriptionPicker", descriptionPicker, true, false)
	context.pages.AddPage("ammountGuesser", ammountGuesser, true, false)
	context.pages.AddPage("dateGuesser", dateGuesser, true, false)
	context.pages.AddPage("empty", tview.NewBox(), true, false)
	context.pages.SwitchToPage("dateGuesser")
	context.Refresh()
	state.AddOnChangeHook(context.Refresh)

	return context, nil
}

func (c Context) GetContent() *tview.Pages { return c.pages }

func (c Context) Refresh() {
	switch c.state.CurrentPhase() {
	case statemod.InputDate:
		c.pages.SwitchToPage("dateGuesser")
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

func NewAccountList(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
	accountGuesser accountguesser.IAccountGuesser,
) (*widgets.ContextualList, error) {
	list, err := widgets.NewContextualList(widgets.ContextualListOptions{
		GetItemsFunc: func() (out []string) {
			for _, acc := range state.JournalMetadata.Accounts() {
				out = append(out, string(acc))
			}
			return out
		},
		SetSelectedFunc: func(s string) {
			state.InputMetadata.SetSelectedPostingAccount(s)
		},
		GetInputFunc: func() string {
			return state.InputMetadata.PostingAccountText()
		},
		GetDefaultFunc: func() (defaultValue string, success bool) {
			acc, success := accountGuesser.Guess()
			return string(acc), success
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build contextual list: %w", err)
	}

	state.AddOnChangeHook(func() { list.Refresh() })
	err = eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.postingaccount.listaction",
		Handler: list.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return list, nil
}

func NewDateGuesser(state *statemod.State) (*tview.TextView, error) {
	guesser := tview.NewTextView()
	refresh := func() {
		if guess, found := state.InputMetadata.GetDateGuess(); found {
			guesser.SetText(guess.Format("2006-01-02") + "\n" + guess.Format("Mon, 02 Jan 2006"))
		} else {
			guesser.SetText("")
		}
	}
	refresh()
	state.AddOnChangeHook(refresh)
	return guesser, nil
}
