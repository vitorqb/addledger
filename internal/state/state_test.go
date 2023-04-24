package state_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/state"
	hledger_mock "github.com/vitorqb/addledger/mocks/hledger"
)

func TestState(t *testing.T) {

	type testcontext struct {
		hookCallCounter int
		state           *State
	}

	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "Notify on change of JournalEntryInput",
			run: func(t *testing.T, c *testcontext) {
				c.state.JournalEntryInput.SetDescription("FOO")
				assert.Equal(t, 1, c.hookCallCounter)

			},
		},
		{
			name: "NextPhase",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, c.state.CurrentPhase(), InputDate)
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputDescription)
				assert.Equal(t, 1, c.hookCallCounter)
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), InputPostingAccount)
				assert.Equal(t, 2, c.hookCallCounter)

			},
		},
		{
			name: "SetPhase",
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(InputPostingAccount)
				assert.Equal(t, InputPostingAccount, c.state.CurrentPhase())
				assert.Equal(t, 1, c.hookCallCounter)

				c.state.SetPhase(InputDescription)
				assert.Equal(t, InputDescription, c.state.CurrentPhase())
				assert.Equal(t, 2, c.hookCallCounter)
			},
		},
		{
			name: "SetAccounts",
			run: func(t *testing.T, c *testcontext) {
				c.state.SetAccounts([]string{"FOO"})
				assert.Equal(t, 1, c.hookCallCounter)
				accounts := c.state.GetAccounts()
				assert.Equal(t, []string{"FOO"}, accounts)
			},
		},
		{
			name: "LoadMetadata",
			run: func(t *testing.T, c *testcontext) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				hledgerClient := hledger_mock.NewMockIClient(ctrl)
				hledgerClient.EXPECT().Accounts().Return([]string{"FOO"}, nil)

				err := c.state.LoadMetadata(hledgerClient)
				assert.Nil(t, err)
				assert.Equal(t, 1, c.hookCallCounter)
				assert.Equal(t, []string{"FOO"}, c.state.GetAccounts())
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.hookCallCounter = 0
			c.state = InitialState()
			c.state.AddOnChangeHook(func() { c.hookCallCounter++ })
			tc.run(t, c)
		})
	}
}
