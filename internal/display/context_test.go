package display_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/display/widgets"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	mock_accountguesser "github.com/vitorqb/addledger/mocks/accountguesser"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

var expectedDate1String = "1993-11-23\nTue, 23 Nov 1993"

func TestNewDateGuesser(t *testing.T) {
	type testcontext struct {
		state   *state.State
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
			c.state = state.InitialState()
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
		state          *state.State
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
				c.accountGuesser.
					EXPECT().
					Guess(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(journal.Account("GUESS"), true)
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
			c.state = state.InitialState()
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
