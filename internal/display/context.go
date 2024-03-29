package display

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/userinput"
)

type Refreshable interface {
	Refresh()
}

type Context struct {
	*tview.Pages
	state *statemod.State
}

// ContextEntry represents a widget inside the context
type ContextWidget struct {
	// The name of the page where the widget is located
	PageName string
	// The widget itself
	Widget tview.Primitive
}

func NewContext(state *statemod.State, widgets []ContextWidget) (*Context, error) {
	// Creates Context
	context := new(Context)
	context.state = state

	// Add all pages to the context
	context.Pages = tview.NewPages()
	context.SetBorder(true)
	for _, widget := range widgets {
		context.AddPage(widget.PageName, widget.Widget, true, false)
	}

	// Add a hook to refresh the widgets when the current page changes.
	context.SetChangedFunc(func() {
		_, page := context.GetFrontPage()
		if refreshablePage, ok := page.(Refreshable); ok {
			refreshablePage.Refresh()
		}
	})

	// Switch to the initial page
	context.SwitchToPage("dateGuesser")

	context.Refresh()
	state.AddOnChangeHook(context.Refresh)

	return context, nil
}

func (c Context) Refresh() {
	switch c.state.CurrentPhase() {
	case statemod.InputDate:
		c.MaybeSwitchToPage("dateGuesser")
	case statemod.InputPostingAccount:
		c.MaybeSwitchToPage("accountList")
	case statemod.InputDescription:
		c.MaybeSwitchToPage("descriptionPicker")
	case statemod.InputPostingAmmount:
		c.MaybeSwitchToPage("ammountGuesser")
	case statemod.InputTags:
		c.MaybeSwitchToPage("tagsPicker")
	default:
		c.MaybeSwitchToPage("empty")
	}
}

func (c Context) MaybeSwitchToPage(pageName string) {
	currentPage, _ := c.GetFrontPage()
	if currentPage != pageName {
		c.SwitchToPage(pageName)
	}
}

func NewAccountList(state *statemod.State, eventbus eventbusmod.IEventBus) (*widgets.ContextualList, error) {
	list, err := widgets.NewContextualList(widgets.ContextualListOptions{
		GetItemsFunc: func() (out []string) {
			// List all accounts
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
			guess, success := state.InputMetadata.GetPostingAccountGuess()
			return string(guess), success
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

type TagsPicker struct{ *widgets.ContextualList }

// TagsPicker presents a list of tags to the user, and allows they to select one.
func NewTagsPicker(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) (*TagsPicker, error) {
	contextListOpts := widgets.ContextualListOptions{
		GetItemsFunc: func() []string {
			tagsStr := []string{}
			for _, tag := range state.JournalMetadata.Tags() {
				tagsStr = append(tagsStr, userinput.TagToText(tag))
			}
			return tagsStr
		},
		SetSelectedFunc: func(x string) {
			tag, _ := userinput.TextToTag(x)
			state.InputMetadata.SetSelectedTag(tag)
		},
		GetInputFunc: func() string {
			return state.InputMetadata.TagText()
		},
		EmptyInputAction: widgets.EmptyInputActionShowCustom(func() []string {
			matchingTransaction := state.InputMetadata.MatchingTransactions()
			if len(matchingTransaction) == 0 {
				return []string{}
			}
			tags := userinput.TagsToText(matchingTransaction[0].Tags)
			return tags
		}),
	}
	contextualList, err := widgets.NewContextualList(contextListOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create contextual list: %w", err)
	}
	state.AddOnChangeHook(func() { contextualList.Refresh() })
	err = eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.tag.listaction",
		Handler: contextualList.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return &TagsPicker{contextualList}, nil
}
