package controller_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/eventbus"
)

func TestInputController(t *testing.T) {

	type testcontext struct {
		state       *statemod.State
		controller  *InputController
		initError   error
		bytesBuffer *bytes.Buffer
		eventBus    *MockIEventBus
	}

	type testcase struct {
		name string
		opts func(t *testing.T, c *testcontext) []Opt
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "NewController missing output causes error",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{}
			},
			run: func(t *testing.T, c *testcontext) {
				assert.ErrorContains(t, c.initError, "missing output")
			},
		},
		{
			name: "NewController missing eventBus causes error",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{WithOutput(c.bytesBuffer)}
			},
			run: func(t *testing.T, c *testcontext) {
				assert.ErrorContains(t, c.initError, "missing Event Bus")
			},
		},
		{
			name: "OnDateInput",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				assert.Nil(t, c.initError)
				date, _ := time.Parse(time.RFC3339, "2022-01-01")
				c.controller.OnDateInput(date)
				assert.Equal(t, statemod.InputDescription, c.state.CurrentPhase())
				foundDate, _ := c.state.JournalEntryInput.GetDate()
				assert.Equal(t, date, foundDate)
			},
		},
		{
			name: "OnDescriptionInput",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputDescription)
				c.controller.OnDescriptionInput("FOO")
				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())
				foundDescription, _ := c.state.JournalEntryInput.GetDescription()
				assert.Equal(t, "FOO", foundDescription)
			},
		},
		{
			name: "OnAccountInput and empty account",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAccountDone("")
				assert.Equal(t, statemod.Confirmation, c.state.CurrentPhase())
			},
		},
		{
			name: "OnAccountInput not empty",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.controller.OnPostingAccountDone("FOO")
				assert.Equal(t, statemod.InputPostingValue, c.state.CurrentPhase())
				account, _ := c.state.JournalEntryInput.CurrentPosting().GetAccount()
				assert.Equal(t, "FOO", account)
			},
		},
		{
			name: "OnInputRejection",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingValue)
				c.controller.OnInputRejection()
				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())
			},
		},
		{
			name: "OnInputConfirmation",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.JournalEntryInput = testutils.JournalEntryInput1(t)
				c.controller.OnInputConfirmation()
				expected := "\n\n" + testutils.JournalEntryInput1(t).Repr()
				assert.Equal(t, expected, c.bytesBuffer.String())
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputDate)
				_, dateFound := c.state.JournalEntryInput.GetDate()
				assert.False(t, dateFound)
			},
		},
		{
			name: "OnPostingAccountListAcction",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				expectedEvent := eventbus.Event{
					Topic: "input.postingaccount.listaction",
					Data:  listaction.NEXT,
				}
				c.eventBus.EXPECT().Send(expectedEvent)
				c.controller.OnPostingAccountListAcction(listaction.NEXT)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := new(testcontext)
			var bytesBuffer bytes.Buffer
			c.bytesBuffer = &bytesBuffer
			c.state = statemod.InitialState()
			c.eventBus = NewMockIEventBus(ctrl)
			opts := tc.opts(t, c)
			c.controller, c.initError = NewController(c.state, opts...)
			tc.run(t, c)
		})
	}
}
