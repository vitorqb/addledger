package controller_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/dateguesser"
	. "github.com/vitorqb/addledger/mocks/eventbus"
)

var aTime, _ = time.Parse(time.RFC3339, "2022-01-01")
var anAmmount = journal.Ammount{
	Commodity: "BRL",
	Quantity:  decimal.New(9999, -3),
}
var anAmmountStr = "BRL 9.999"
var anotherAmmount = journal.Ammount{
	Commodity: "EUR",
	Quantity:  decimal.New(1220, -2),
}

var anotherAmmountStr = "EUR 12.20"

func TestInputController(t *testing.T) {

	type testcontext struct {
		state       *statemod.State
		controller  *InputController
		initError   error
		bytesBuffer *bytes.Buffer
		eventBus    *MockIEventBus
		dateGuesser *MockIDateGuesser
	}

	type testcase struct {
		name string
		opts func(t *testing.T, c *testcontext) []Opt
		run  func(t *testing.T, c *testcontext)
	}

	defaultOpts := func(t *testing.T, c *testcontext) []Opt {
		return []Opt{
			WithOutput(c.bytesBuffer),
			WithEventBus(c.eventBus),
			WithDateGuesser(c.dateGuesser),
		}
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
			name: "NewController missing eventGuesser causes error",
			opts: func(t *testing.T, c *testcontext) []Opt {
				return []Opt{
					WithOutput(c.bytesBuffer),
					WithEventBus(c.eventBus),
				}
			},
			run: func(t *testing.T, c *testcontext) {
				assert.ErrorContains(t, c.initError, "missing DateGuesser")
			},
		},
		{
			name: "On date change and done",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				assert.Nil(t, c.initError)
				c.dateGuesser.EXPECT().Guess("2022-01-01").Return(aTime, true)
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				assert.Equal(t, statemod.InputDescription, c.state.CurrentPhase())
				foundDate, _ := c.state.JournalEntryInput.GetDate()
				assert.Equal(t, aTime, foundDate)
			},
		},
		{
			name: "On date change but no guess",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				assert.Nil(t, c.initError)
				c.dateGuesser.EXPECT().Guess("aaa").Return(time.Time{}, false)
				c.controller.OnDateChanged("aaa")
				c.controller.OnDateDone()
				assert.Equal(t, statemod.InputDate, c.state.CurrentPhase())
				_, dateFound := c.state.JournalEntryInput.GetDate()
				assert.False(t, dateFound)
			},
		},
		{
			name: "On date cleans up date on second entry",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				assert.Nil(t, c.initError)
				c.dateGuesser.EXPECT().Guess("2023-01-01").Return(aTime, true)
				c.dateGuesser.EXPECT().Guess("aaa").Return(time.Time{}, false)
				c.controller.OnDateChanged("2023-01-01")
				foundDate, _ := c.state.JournalEntryInput.GetDate()
				assert.Equal(t, aTime, foundDate)
				c.controller.OnDateChanged("aaa")
				_, dateFound := c.state.JournalEntryInput.GetDate()
				assert.False(t, dateFound)
			},
		},
		{
			name: "Description input changes and done",
			opts: defaultOpts,
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
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAccountDone("")
				assert.Equal(t, statemod.Confirmation, c.state.CurrentPhase())
			},
		},
		{
			name: "OnAccountInput not empty",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.controller.OnPostingAccountDone("FOO")
				assert.Equal(t, statemod.InputPostingAmmount, c.state.CurrentPhase())
				account, _ := c.state.JournalEntryInput.CurrentPosting().GetAccount()
				assert.Equal(t, "FOO", account)
			},
		},
		{
			name: "OnInputRejection",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAmmount)
				c.controller.OnInputRejection()
				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())
			},
		},
		{
			name: "OnInputConfirmation",
			opts: defaultOpts,
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
			opts: defaultOpts,
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
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAccountChanged("FOO")
				assert.Equal(t, "FOO", c.state.InputMetadata.PostingAccountText())
			},
		},
		{
			name: "OnPostingAmmountDone from input",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetPostingAmmountInput(anAmmount)
				c.controller.OnPostingAmmountDone(input.Input)
				posting, _ := c.state.JournalEntryInput.GetPosting(0)
				ammount, _ := posting.GetAmmount()
				assert.Equal(t, anAmmount, ammount)
			},
		},
		{
			name: "OnPostingAmmountDone from guess",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetPostingAmmountGuess(anAmmount)
				c.controller.OnPostingAmmountDone(input.Context)
				posting, _ := c.state.JournalEntryInput.GetPosting(0)
				ammount, _ := posting.GetAmmount()
				assert.Equal(t, anAmmount, ammount)
			},
		},
		{
			name: "OnPostingAmmountChanged saves to state success",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAmmountChanged("EUR 12.20")
				text := c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "EUR 12.20", text)
				ammount, found := c.state.InputMetadata.GetPostingAmmountInput()
				assert.True(t, found)
				assert.Equal(t, anotherAmmount, ammount)
			},
		},
		{
			name: "OnPostingAmmountChanged saves to state parse fails",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.controller.OnPostingAmmountChanged("aaa")
				text := c.state.InputMetadata.GetPostingAmmountText()
				assert.Equal(t, "aaa", text)
				_, found := c.state.InputMetadata.GetPostingAmmountInput()
				assert.False(t, found)
			},
		},
		{
			name: "OnPostingAccountSelectedFromContext",
			opts: defaultOpts,
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
			opts: defaultOpts,
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
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.NextPhase()
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputDescription)
				c.controller.OnUndo()
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputDate)
			},
		},
		{
			name: "OnUndo cleans up the last user input",
			opts: defaultOpts,
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
			opts: defaultOpts,
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
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Must have 2 postings - the filled one and an empty one.
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				assert.Equal(t, c.state.CurrentPhase(), statemod.InputPostingAccount)

				// Undo
				c.controller.OnUndo()

				// Must not have a single posting without value
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 1)
				_, ammountFound := c.state.JournalEntryInput.CurrentPosting().GetAmmount()
				assert.False(t, ammountFound)
			},
		},
		{
			name: "Must remove empty posting after finished entering postings",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone("BAR2")
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Should have 3 postings - 2 filled and 1 empty
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 3)
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				_, accFound := lastPosting.GetAccount()
				assert.False(t, accFound)
				_, ammountFound := lastPosting.GetAmmount()
				assert.False(t, ammountFound)

				// Submit
				c.controller.OnPostingAccountDone("")

				// Must have deleted the last posting
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				lastPosting = c.state.JournalEntryInput.CurrentPosting()
				accValue, accFound := lastPosting.GetAccount()
				assert.True(t, accFound)
				assert.Equal(t, "BAR2", accValue)
				ammountValue, ammountFound := lastPosting.GetAmmount()
				assert.True(t, ammountFound)
				assert.Equal(t, journal.Ammount{Commodity: "EUR", Quantity: decimal.New(1220, -2)}, ammountValue)
			},
		},
		{
			name: "Must have empty posting after confirmation rejection",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone("BAR2")
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)

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
				_, ammountFound := lastPosting.GetAmmount()
				assert.False(t, ammountFound)
			},
		},
		{
			name: "OnPostingAmmountDone",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone("BAR")
				c.controller.OnPostingAmmountChanged(anAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// The ammount has been saved to the posting
				posting, found := c.state.JournalEntryInput.GetPosting(0)
				assert.True(t, found)
				ammount, found := posting.GetAmmount()
				assert.True(t, found)
				assert.Equal(t, anAmmount, ammount)

				// Phase is set to posting account
				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())

				// A new empty posting is there
				assert.Equal(t, 2, c.state.JournalEntryInput.CountPostings())
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
			c.dateGuesser = NewMockIDateGuesser(ctrl)
			opts := tc.opts(t, c)
			c.controller, c.initError = NewController(c.state, opts...)
			tc.run(t, c)
		})
	}
}
