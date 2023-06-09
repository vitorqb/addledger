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

var aTime, _ = time.Parse(time.RFC3339, "2022-01-01")

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
				c.controller.OnDateInput(aTime)
				assert.Equal(t, statemod.InputDescription, c.state.CurrentPhase())
				foundDate, _ := c.state.JournalEntryInput.GetDate()
				assert.Equal(t, aTime, foundDate)
			},
		},
		{
			name: "Description input changes and done",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputDescription)
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
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
		{
			name: "OnPostingAccountChanged",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAccountChanged("FOO")
				assert.Equal(t, "FOO", c.state.InputMetadata.PostingAccountText())
			},
		},
		{
			name: "OnPostingAccountSelectedFromContext",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetSelectedPostingAccount("FOO")
				c.controller.OnPostingAccountSelectedFromContext()
				posting := c.state.JournalEntryInput.CurrentPosting()
				acc, ok := posting.GetAccount()
				assert.True(t, ok)
				assert.Equal(t, "FOO", acc)
			},
		},
		{
			name: "OnPostingAccountInsertFromContext",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetSelectedPostingAccount("FOO")
				expEvent := eventbus.Event{
					Topic: "input.postingaccount.settext",
					Data:  "FOO",
				}
				c.eventBus.EXPECT().Send(expEvent)
				c.controller.OnPostingAccountInsertFromContext()
			},
		},
		{
			name: "OnUndo moves the state page back",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputDescription)
				c.controller.OnUndo()
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputDate)
			},
		},
		{
			name: "OnUndo cleans up the last user input",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.JournalEntryInput.SetDate(aTime)
				c.state.NextPhase()
				c.controller.OnUndo()
				_, ok := c.state.JournalEntryInput.GetDate()
				assert.False(t, ok)
			},
		},
		{
			name: "OnUndo with posting account ",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.state.JournalEntryInput.SetDate(aTime)
				c.state.NextPhase()
				c.controller.OnUndo()
				_, ok := c.state.JournalEntryInput.GetDate()
				assert.False(t, ok)
			},
		},
		{
			name: "Undo after first posting is entered ",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnDateInput(aTime)
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingValueInput("EUR 12.20")

				// Must have 2 postings - the filled one and an empty one.
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputPostingAccount)

				// Undo
				c.controller.OnUndo()

				// Must not have a single posting without value
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 1)
				_, valueFound := c.state.JournalEntryInput.CurrentPosting().GetValue()
				assert.False(t, valueFound)
			},
		},
		{
			name: "Must remove empty posting after finished entering postings",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnDateInput(aTime)
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingValueInput("EUR 12.20")
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone("BAR2")
				c.controller.OnPostingValueInput("EUR -12.20")

				// Should have 3 postings - 2 filled and 1 empty
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 3)
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				_, accFound := lastPosting.GetAccount()
				assert.False(t, accFound)
				_, valueFound := lastPosting.GetValue()
				assert.False(t, valueFound)

				// Submit
				c.controller.OnPostingAccountDone("")

				// Must have deleted the last posting
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				lastPosting = c.state.JournalEntryInput.CurrentPosting()
				accValue, accFound := lastPosting.GetAccount()
				assert.True(t, accFound)
				assert.Equal(t, "BAR2", accValue)
				valueValue, valueFound := lastPosting.GetValue()
				assert.True(t, valueFound)
				assert.Equal(t, "EUR -12.20", valueValue)
			},
		},
		{
			name: "Must have empty posting after confirmation rejection",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnDateInput(aTime)
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingValueInput("EUR 12.20")
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone("BAR2")
				c.controller.OnPostingValueInput("EUR -12.20")

				// Goes to confirm page
				c.controller.OnPostingAccountDone("")

				// Should have 2 filled postings (empty one deleted)
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)

				// User decides to go back
				c.controller.OnInputRejection()

				// Should have 3 postings - 2 filled and 1 empty
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 3)
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				_, accFound := lastPosting.GetAccount()
				assert.False(t, accFound)
				_, valueFound := lastPosting.GetValue()
				assert.False(t, valueFound)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
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
