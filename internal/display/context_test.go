package display_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/display/widgets"
	"github.com/vitorqb/addledger/internal/journal"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	mock_accountguesser "github.com/vitorqb/addledger/mocks/accountguesser"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

var expectedDate1String = "1993-11-23\nTue, 23 Nov 1993"

// FakeRefreshablePrimitive is a fake tview.Primitive that implements
// the Refreshable interface
type FakeRefreshablePrimitive struct {
	tview.Primitive
	RefreshCallCount int
}

// Refresh implements the Refreshable interface
func (f *FakeRefreshablePrimitive) Refresh() {
	f.RefreshCallCount++
}

// NewFakeRefreshablePrimitive creates a new FakeRefreshablePrimitive
func NewFakeRefreshablePrimitive() *FakeRefreshablePrimitive {
	return &FakeRefreshablePrimitive{tview.NewBox(), 0}
}

func TestNewContext(t *testing.T) {
	type testcontext struct {
		state *statemod.State
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Calls Refresh on widgets on page change",
			run: func(c *testcontext, t *testing.T) {
				primitive1 := NewFakeRefreshablePrimitive()
				primitive2 := NewFakeRefreshablePrimitive()
				widget1 := ContextWidget{"dateGuesser", primitive1}
				widget2 := ContextWidget{"accountList", primitive2}
				widgets := []ContextWidget{widget1, widget2}
				_, err := NewContext(c.state, widgets)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, 1, primitive1.RefreshCallCount)
				c.state.SetPhase(statemod.InputDate)
				assert.Equal(t, 2, primitive1.RefreshCallCount)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.state = statemod.InitialState()
			tc.run(c, t)
		})
	}

}

func TestNewDateGuesser(t *testing.T) {
	type testcontext struct {
		state   *statemod.State
		guesser *tview.TextView
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "sets text from state",
			run: func(c *testcontext, t *testing.T) {
				c.state.InputMetadata.SetDateGuess(testutils.Date1(t))
				assert.Equal(t, expectedDate1String, c.guesser.GetText(true))
			},
		},
		{
			name: "clears when state clears",
			run: func(c *testcontext, t *testing.T) {
				c.state.InputMetadata.SetDateGuess(testutils.Date1(t))
				c.state.InputMetadata.ClearDateGuess()
				assert.Equal(t, "", c.guesser.GetText(true))
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.guesser, err = NewDateGuesser(c.state)
			if err != nil {
				t.Fatal(err)
			}
			tc.run(c, t)
		})
	}
}

func TestAccountList(t *testing.T) {
	type testcontext struct {
		state          *statemod.State
		eventBus       *mock_eventbus.MockIEventBus
		accountGuesser *mock_accountguesser.MockIAccountGuesser
		accountList    *widgets.ContextualList
	}
	type testcase struct {
		name  string
		setup func(c *testcontext)
		run   func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Set's default account as first item",
			setup: func(c *testcontext) {
				c.accountGuesser.EXPECT().Guess().Return(journal.Account("GUESS"), true)
			},
			run: func(c *testcontext, t *testing.T) {
				assert.Equal(t, "GUESS", c.state.InputMetadata.SelectedPostingAccount())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.eventBus = mock_eventbus.NewMockIEventBus(ctrl)
			c.eventBus.EXPECT().Subscribe(gomock.Any())
			c.accountGuesser = mock_accountguesser.NewMockIAccountGuesser(ctrl)
			if tc.setup != nil {
				tc.setup(c)
			}
			c.accountList, err = NewAccountList(c.state, c.eventBus, c.accountGuesser)
			if err != nil {
				t.Fatal(err)
			}
			tc.run(c, t)
		})
	}
}

func TestTagsPicker(t *testing.T) {
	type testcontext struct {
		state    *statemod.State
		eventBus *mock_eventbus.MockIEventBus
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Displays all tags",
			run: func(c *testcontext, t *testing.T) {
				transaction := testutils.Transaction_3(t)
				c.state.JournalMetadata.AppendTransaction(*transaction)
				tagsPicker, err := NewTagsPicker(c.state, c.eventBus)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, 1, tagsPicker.GetItemCount())
				displayedTag, _ := tagsPicker.GetItemText(0)
				assert.Equal(t, "trip:brazil", displayedTag)
			},
		},
		{
			name: "Sets selected tag after user input",
			run: func(c *testcontext, t *testing.T) {
				transaction := testutils.Transaction_3(t)
				transaction.Tags = []journal.Tag{
					{Name: "aaaa", Value: "bbbb"},
					{Name: "cccc", Value: "dddd"},
				}
				c.state.JournalMetadata.AppendTransaction(*transaction)
				tagsPicker, err := NewTagsPicker(c.state, c.eventBus)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, 2, tagsPicker.GetItemCount())

				c.state.InputMetadata.SetTagText("cccc:d")
				assert.Equal(t, 1, tagsPicker.GetItemCount())
				displayedTag, _ := tagsPicker.GetItemText(0)
				assert.Equal(t, "cccc:dddd", displayedTag)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.eventBus = mock_eventbus.NewMockIEventBus(ctrl)
			c.eventBus.EXPECT().Subscribe(gomock.Any())
			tc.run(c, t)
		})
	}
}
