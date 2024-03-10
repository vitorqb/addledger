package widgets

import (
	"fmt"
	"sort"

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
	// GetDefaultFunc is a function that returns the default value.
	GetDefaultFunc func() (defaultValue string, success bool)
	// EmptyInputAction is the action to be taken when the input is empty.
	EmptyInputAction EmptyInputAction
}

// EmptyInputAction represents the action to be taken when the input is empty.
type EmptyInputAction func(c *ContextualList)

// EmptyInputActionHideAll hides all items when input is empty.
var EmptyInputHideItems EmptyInputAction = func(cl *ContextualList) {
	cl.Clear()
}

// EmptyInputActionShowAll shows all items when input is empty.
var EmptyInputActionShowAll EmptyInputAction = func(cl *ContextualList) {
	cl.Clear()
	defaultValue, hasDefault := cl.getDefaultFunc()
	if hasDefault {
		cl.AddItem(defaultValue, "", 0, nil)
	}
	for _, item := range cl.getItemsFunc() {
		cl.AddItem(item, "", 0, nil)
	}
}

// EmptyInputActionShowCustom shows a custom list of items when input is empty.
func EmptyInputActionShowCustom(getItems func() []string) EmptyInputAction {
	return func(cl *ContextualList) {
		cl.Clear()
		// First row stands for "no selection"
		cl.AddItem("", "", 0, nil)
		for _, item := range getItems() {
			cl.AddItem(item, "", 0, nil)
		}
	}
}

// ContextualList is a List widget for the context that allows the user to
// select an entry from it for an input.
type ContextualList struct {
	*tview.List
	getItemsFunc         func() []string
	getInputFunc         func() string
	setSelectedFunc      func(string)
	getDefaultFunc       func() (defaultValue string, success bool)
	emptyInputAction     EmptyInputAction
	isRefreshing         bool
	isHandlingListAction bool
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
	if options.GetDefaultFunc == nil {
		options.GetDefaultFunc = func() (string, bool) { return "", false }
	}
	if options.EmptyInputAction == nil {
		options.EmptyInputAction = EmptyInputActionShowAll
	}

	// Builds list
	list := &ContextualList{
		List:             tview.NewList(),
		getItemsFunc:     options.GetItemsFunc,
		getInputFunc:     options.GetInputFunc,
		setSelectedFunc:  options.SetSelectedFunc,
		getDefaultFunc:   options.GetDefaultFunc,
		emptyInputAction: options.EmptyInputAction,
	}
	list.ShowSecondaryText(false)
	list.SetChangedFunc(func(_ int, mainText, _ string, _ rune) {
		logrus.WithField("text", mainText).Debug("ContextualList changed")
		list.setSelectedFunc(mainText)
	})
	list.Refresh()
	return list, nil
}

// RestoreIndex tries to restore the index of the list to the given index.
func (cl *ContextualList) RestoreIndex(index int) {
	if index < 0 || index >= cl.GetItemCount() {
		return
	}
	cl.SetCurrentItem(index)
}

// HandleAction handles a ListAction (e.g. next, prev, etc).
func (cl *ContextualList) HandleAction(action listaction.ListAction) {
	cl.isHandlingListAction = true
	defer func() { cl.isHandlingListAction = false }()
	logrus.WithField("action", action).Debug("Received listAction")
	switch action {
	case listaction.NEXT:
		currentItem := cl.GetCurrentItem()
		cl.SetCurrentItem(currentItem + 1)
	case listaction.PREV:
		currentItem := cl.GetCurrentItem()
		cl.SetCurrentItem(currentItem - 1)
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
	// If we are already refreshing, do nothing.
	if cl.isRefreshing {
		return
	}
	cl.isRefreshing = true
	defer func() { cl.isRefreshing = false }()

	// If we are handling an action, don't refresh.
	if cl.isHandlingListAction {
		return
	}

	// After we refresh, if we have 0 items, set selected to ""
	defer func() {
		if cl.GetItemCount() == 0 {
			cl.setSelectedFunc("")
		}
	}()

	// Ensure that after we refresh the state is up-to-date with
	// the selected item
	defer func() {
		if cl.GetItemCount() > 0 {
			text, _ := cl.GetItemText(cl.GetCurrentItem())
			cl.setSelectedFunc(text)
		}
	}()

	input := cl.getInputFunc()

	// If the input is empty, dispatch to the empty input action
	if input == "" {
		cl.emptyInputAction(cl)
		return
	}

	defer cl.RestoreIndex(cl.GetCurrentItem())
	cl.Clear()

	// Input is not empty - match and sort by match
	matches := fuzzy.RankFindFold(input, cl.getItemsFunc())
	sort.Sort(matches)
	for _, match := range matches {
		cl.AddItem(match.Target, "", 0, nil)
	}
}
