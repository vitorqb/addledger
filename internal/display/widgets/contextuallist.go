package widgets

import (
	"fmt"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
)

// ContextualListOptions represents all options for a ContextualList
type ContextualListOptions struct {
	// GetItemsFunc is a function that must return the items to be printed.
	GetItemsFunc func() []string
	// SetSelectedFunc is a function called when an item is selected.
	SetSelectedFunc func(string)
	// GetInputFunc is a function that must return the user input.
	GetInputFunc func() string
}

// ContextualList is a List widget for the context that allows the user to
// select an entry from it for an input.
type ContextualList struct {
	*tview.List
	inputCache      string
	getItemsFunc    func() []string
	getInputFunc    func() string
	setSelectedFunc func(string)
}

// NewContextualList creates a new ContextualList. `getItemsFunc` is a
// function that must return all available items. `setSelectedFunc`
// is a function that is called with the selected item every time it
// changes. `getInputFunc` is a function that returns the current
// input used to filter the entries.
func NewContextualList(options ContextualListOptions) (*ContextualList, error) {
	// Validates input
	if options.GetInputFunc == nil {
		return nil, fmt.Errorf("missing GetInputFunc")
	}
	if options.GetItemsFunc == nil {
		return nil, fmt.Errorf("missing GetItemsFunc")
	}

	// Builds list
	list := &ContextualList{
		tview.NewList(),
		"",
		options.GetItemsFunc,
		options.GetInputFunc,
		options.SetSelectedFunc,
	}
	list.ShowSecondaryText(false)
	list.SetChangedFunc(func(_ int, mainText, _ string, _ rune) {
		logrus.WithField("text", mainText).Debug("AccountList changed")
		list.setSelectedFunc(mainText)
	})
	for _, item := range list.getItemsFunc() {
		list.AddItem(item, "", 0, nil)
	}
	return list, nil
}

// HandleAction handles a ListAction (e.g. next, prev, etc).
func (cl *ContextualList) HandleAction(action listaction.ListAction) {
	logrus.WithField("action", action).Debug("Received listAction")
	switch action {
	case listaction.NEXT:
		eventKey := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
		cl.InputHandler()(eventKey, func(p tview.Primitive) {})
	case listaction.PREV:
		eventKey := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
		cl.InputHandler()(eventKey, func(p tview.Primitive) {})
	case listaction.NONE:
	default:
	}
}

// HandleActionFromEvent handles an action inside an eventbus.Event
func (cl *ContextualList) HandleActionFromEvent(e eventbus.Event) {
	listAction, ok := e.Data.(listaction.ListAction)
	if !ok {
		logrus.WithField("event", e).Error("received event w/ unexpected data")
		return
	}
	cl.HandleAction(listAction)
}

// Refresh deletes all items, queries for them again and puts together the list only
// with items that match the current input.
func (cl *ContextualList) Refresh() {
	input := cl.getInputFunc()

	// If cache hit, just return
	if cl.inputCache == input {
		return
	}

	// No cache - clear inputs
	cl.inputCache = input
	cl.Clear()

	// If the input is empty, don't do any matches so we can preserve the order
	if input == "" {
		for _, item := range cl.getItemsFunc() {
			cl.AddItem(item, "", 0, nil)
		}
		return
	}

	// Input is not empty - match and sort by match
	matches := fuzzy.RankFindFold(input, cl.getItemsFunc())
	sort.Sort(matches)
	for _, match := range matches {
		cl.AddItem(match.Target, "", 0, nil)
	}
}
