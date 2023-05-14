package widgets

import (
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
)

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
func NewContextualList(
	getItemsFunc func() []string,
	setSelectedFunc func(string),
	getInputFunc func() string,
) *ContextualList {
	list := &ContextualList{tview.NewList(), "", getItemsFunc, getInputFunc, setSelectedFunc}
	list.ShowSecondaryText(false)
	list.SetChangedFunc(func(_ int, mainText, _ string, _ rune) {
		logrus.WithField("text", mainText).Debug("AccountList changed")
		setSelectedFunc(mainText)
	})
	for _, item := range getItemsFunc() {
		list.AddItem(item, "", 0, nil)
	}
	return list
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
	logrus.WithField("input", input).Debug("Refreshing ContextualList")
	if cl.inputCache == input {
		return
	}
	cl.inputCache = input
	cl.Clear()
	matches := fuzzy.RankFindFold(input, cl.getItemsFunc())
	sort.Sort(matches)
	for _, match := range matches {
		cl.AddItem(match.Target, "", 0, nil)
	}
}
