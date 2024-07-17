package statement

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/state"
)

//go:generate $MOCKGEN --source=main.go --destination=../../../mocks/display/statement/main_mock.go

const DateFormat = "2006-01-02"

type Controller interface {
	LoadRequest()
	HideModal()
	DiscardStatementEntry(index int)
}

type Modal struct {
	*tview.Flex
	controller Controller
}

// Contextual information about the current state of the statement display.
type Context struct {
	// Tracks the index of the selected statement in the state.GetStatementEntries
	// slice.
	SelectedStatementIndex int
}

// RuneAction represents a modal action from a simple rune being presset. The
// action is ran and the modal is closed afterwards.
type RuneAction struct {
	Rune        rune
	Description string
	Action      func(Controller, *Context)
}

type CommandBar struct{ *tview.TextArea }

func NewCommandBar(actions []RuneAction) tview.Primitive {
	var textBuilder strings.Builder
	textBuilder.WriteString("| ")
	for _, a := range actions {
		textBuilder.WriteString(string(a.Rune) + " - " + a.Description + " | ")
	}
	o := &CommandBar{tview.NewTextArea()}
	o.SetText(textBuilder.String(), true)
	return o
}

type Table struct{ *tview.Table }

func (t *Table) Refresh(entries []finance.StatementEntry) {
	t.Clear()
	for i, e := range entries {
		t.SetCell(i, 0, tview.NewTableCell(e.Account))
		t.SetCell(i, 1, tview.NewTableCell(e.Date.Format(DateFormat)))
		t.SetCell(i, 2, tview.NewTableCell("\""+e.Description+"\""))
		t.SetCell(i, 3, tview.NewTableCell(e.Ammount.Commodity))
		t.SetCell(i, 4, tview.NewTableCell(e.Ammount.Quantity.StringFixed(2)))
	}
	t.SetOffset(0, 0)
}

func NewTable() *Table {
	o := &Table{tview.NewTable()}
	o.SetBorder(true)
	o.SetSelectable(true, false)
	return o
}

func CreateModal(controller Controller, state *state.State, actions []RuneAction) *Modal {
	context := &Context{}

	table := NewTable()
	table.SetSelectionChangedFunc(func(row, col int) {
		context.SelectedStatementIndex = row
	})
	table.Refresh(state.GetStatementEntries())
	state.AddOnChangeHook(func() { table.Refresh(state.GetStatementEntries()) })

	modal := &Modal{tview.NewFlex(), controller}
	modal.SetDirection(tview.FlexRow)
	modal.AddItem(table, 0, 8, true)
	modal.AddItem(NewCommandBar(actions), 0, 2, false)
	modal.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			for _, action := range actions {
				if event.Rune() == action.Rune {
					action.Action(controller, context)
					return nil
				}
			}
		case tcell.KeyEscape, tcell.KeyCtrlQ:
			modal.controller.HideModal()
			return nil
		}
		return event
	})
	return modal
}

// Actions available in the modal.
var defaultActions = []RuneAction{
	{'l', "Load Statement", func(c Controller, ctx *Context) {
		c.LoadRequest()
		c.HideModal()
	}},
	{'d', "Discard Statement Entry", func(c Controller, ctx *Context) { c.DiscardStatementEntry(ctx.SelectedStatementIndex) }},
	{'q', "Quit", func(c Controller, ctx *Context) { c.HideModal() }},
}

func NewModal(controller Controller, state *state.State) *Modal {
	return CreateModal(controller, state, defaultActions)
}
