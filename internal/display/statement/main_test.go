package statement_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/statement"
	"github.com/vitorqb/addledger/internal/finance"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/display/statement"
)

func TestNewCommandBar(t *testing.T) {
	t.Run("displays actions", func(t *testing.T) {
		actions := []RuneAction{
			{'a', "Foo", func(c Controller, ctx *Context) {}},
			{'b', "Bar", func(c Controller, ctx *Context) {}},
		}
		commandBar := NewCommandBar(actions)
		assert.Contains(t, commandBar.(*CommandBar).GetText(), "a - Foo | ")
		assert.Contains(t, commandBar.(*CommandBar).GetText(), "b - Bar | ")
	})
}

func TestTable(t *testing.T) {
	var setup = func() (tcell.SimulationScreen, *Table) {
		ss := tcell.NewSimulationScreen("UTF-8")
		_ = ss.Init()
		ss.SetSize(60, 10)
		table := NewTable()
		table.SetRect(0, 0, 60, 10)
		table.Draw(ss)
		ss.Sync()
		return ss, table
	}

	t.Run("displays cells", func(t *testing.T) {
		ss, table := setup()
		entries := []finance.StatementEntry{
			{
				Account:     "Foo",
				Date:        testutils.Date1(t),
				Description: "Bar",
				Ammount:     *testutils.Ammount_1(t),
			},
		}
		table.Refresh(entries)
		table.Draw(ss)
		ss.Sync()
		text := testutils.ExtractText(ss)
		exp := []string{
			"┌──────────────────────────────────────────────────────────┐",
			"│Foo 1993-11-23 \"Bar\" EUR 2.20                             │",
			"│                                                          │",
			"│                                                          │",
			"│                                                          │",
			"│                                                          │",
			"│                                                          │",
			"│                                                          │",
			"│                                                          │",
			"└──────────────────────────────────────────────────────────┘\n",
		}
		assert.Equal(t, strings.Join(exp, "\n"), text)
	})

	t.Run("displays first row by default", func(t *testing.T) {
		ss, table := setup()
		entries := []finance.StatementEntry{}
		for i := 0; i < 20; i++ {
			entry := finance.StatementEntry{
				Account:     "Foo" + fmt.Sprint(i),
				Date:        testutils.Date1(t),
				Description: "Bar",
				Ammount:     *testutils.Ammount_1(t),
			}
			entries = append(entries, entry)
		}
		table.Refresh(entries)
		table.Draw(ss)
		ss.Sync()
		text := testutils.ExtractText(ss)
		exp := []string{
			"┌──────────────────────────────────────────────────────────┐",
			"│Foo0 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo1 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo2 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo3 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo4 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo5 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo6 1993-11-23 \"Bar\" EUR 2.20                            │",
			"│Foo7 1993-11-23 \"Bar\" EUR 2.20                            │",
			"└──────────────────────────────────────────────────────────┘\n",
		}
		assert.Equal(t, strings.Join(exp, "\n"), text)
	})

	t.Run("clear table on refresh", func(t *testing.T) {
		_, table := setup()
		entry1 := finance.StatementEntry{Account: "Foo"}
		entry2 := finance.StatementEntry{Account: "Bar"}
		entries := []finance.StatementEntry{entry1, entry2}
		table.Refresh(entries)
		assert.Contains(t, table.GetCell(0, 0).Text, "Foo")
		assert.Contains(t, table.GetCell(1, 0).Text, "Bar")
		assert.Equal(t, 2, table.GetRowCount())
		table.Refresh([]finance.StatementEntry{entry1})
		assert.Contains(t, table.GetCell(0, 0).Text, "Foo")
		assert.Equal(t, 1, table.GetRowCount())
	})
}

func TestNewModal(t *testing.T) {
	fakeSetFocus := func(p tview.Primitive) {}
	lEventKey := tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone)
	escapeEventKey := tcell.NewEventKey(tcell.KeyEscape, '\x00', tcell.ModNone)

	type setupOptions struct {
		Actions []RuneAction
	}

	setup := func(opts setupOptions) (*MockController, *Modal, *statemod.State, func()) {
		ctrl := gomock.NewController(t)
		controller := NewMockController(ctrl)
		teardown := func() {
			ctrl.Finish()
		}
		state := statemod.InitialState()
		modal := CreateModal(controller, state, opts.Actions)
		return controller, modal, state, teardown
	}

	t.Run("displays command bar", func(t *testing.T) {
		_, modal, _, teardown := setup(setupOptions{})
		defer teardown()
		commandBar := modal.GetItem(1)
		assert.IsType(t, &CommandBar{}, commandBar)
	})

	t.Run("displays statement table", func(t *testing.T) {
		_, modal, _, teardown := setup(setupOptions{})
		defer teardown()
		table := modal.GetItem(0)
		assert.IsType(t, &Table{}, table)
	})

	t.Run("dispatches to action on input", func(t *testing.T) {
		actionCalled := false
		fakeAction := func(c Controller, ctx *Context) { actionCalled = true }
		_, modal, _, teardown := setup(setupOptions{
			Actions: []RuneAction{{'l', "Foo", fakeAction}},
		})
		defer teardown()
		modal.InputHandler()(lEventKey, fakeSetFocus)
		assert.True(t, actionCalled)
	})

	t.Run("updates context selected statement index", func(t *testing.T) {
		var selectedStatementIndex int
		fakeAction := func(c Controller, ctx *Context) {
			selectedStatementIndex = ctx.SelectedStatementIndex
		}
		_, modal, state, teardown := setup(setupOptions{
			Actions: []RuneAction{{'l', "Foo", fakeAction}},
		})
		defer teardown()

		state.SetStatementEntries([]finance.StatementEntry{
			{Account: "foo"},
			{Account: "bar"},
		})

		// Selects the second row & dispatches fakeAction
		modal.GetItem(0).(*Table).Select(1, 0)
		modal.InputHandler()(lEventKey, fakeSetFocus)

		assert.Equal(t, 1, selectedStatementIndex)
	})

	t.Run("hides modal on escape", func(t *testing.T) {
		controller, modal, _, teardown := setup(setupOptions{})
		defer teardown()
		controller.EXPECT().HideModal()
		modal.InputHandler()(escapeEventKey, fakeSetFocus)
	})

	t.Run("refreshes table with entries", func(t *testing.T) {
		_, modal, state, teardown := setup(setupOptions{})
		defer teardown()
		entry := finance.StatementEntry{Account: "foooo"}
		state.SetStatementEntries([]finance.StatementEntry{entry})
		table := modal.GetItem(0)
		assert.Contains(t, table.(*Table).GetCell(0, 0).Text, "foooo")
	})
}
