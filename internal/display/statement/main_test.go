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
			{'a', "Foo", func(c Controller) {}},
			{'b', "Bar", func(c Controller) {}},
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
}

func TestNewModal(t *testing.T) {
	fakeSetFocus := func(p tview.Primitive) {}
	lEventKey := tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone)
	qEventKey := tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone)
	escapeEventKey := tcell.NewEventKey(tcell.KeyEscape, '\x00', tcell.ModNone)

	setup := func() (*MockController, *Modal, *statemod.State, func()) {
		ctrl := gomock.NewController(t)
		controller := NewMockController(ctrl)
		teardown := func() {
			ctrl.Finish()
		}
		state := statemod.InitialState()
		modal := NewModal(controller, state)
		return controller, modal, state, teardown
	}

	t.Run("displays command bar", func(t *testing.T) {
		_, modal, _, teardown := setup()
		defer teardown()
		commandBar := modal.GetItem(1)
		assert.IsType(t, &CommandBar{}, commandBar)
	})

	t.Run("displays statement table", func(t *testing.T) {
		_, modal, _, teardown := setup()
		defer teardown()
		table := modal.GetItem(0)
		assert.IsType(t, &Table{}, table)
	})

	t.Run("dispatches to load statement", func(t *testing.T) {
		controller, modal, _, teardown := setup()
		defer teardown()
		controller.EXPECT().LoadRequest()
		controller.EXPECT().HideModal()
		modal.InputHandler()(lEventKey, fakeSetFocus)
	})

	t.Run("hides modal on q", func(t *testing.T) {
		controller, modal, _, teardown := setup()
		defer teardown()
		controller.EXPECT().HideModal()
		modal.InputHandler()(qEventKey, fakeSetFocus)
	})

	t.Run("hides modal on escape", func(t *testing.T) {
		controller, modal, _, teardown := setup()
		defer teardown()
		controller.EXPECT().HideModal()
		modal.InputHandler()(escapeEventKey, fakeSetFocus)
	})

	t.Run("refreshes table with entries", func(t *testing.T) {
		_, modal, state, teardown := setup()
		defer teardown()
		entry := finance.StatementEntry{Account: "foooo"}
		state.SetStatementEntries([]finance.StatementEntry{entry})
		table := modal.GetItem(0)
		assert.Contains(t, table.(*Table).GetCell(0, 0).Text, "foooo")
	})
}
