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
var anAmmountNeg = journal.Ammount{
	Commodity: anAmmount.Commodity,
	Quantity:  anAmmount.Quantity.Neg(),
}
var anAmmountNegStr = "BRL -9.999"
var anotherAmmount = journal.Ammount{
	Commodity: "EUR",
	Quantity:  decimal.New(1220, -2),
}
var anotherAmmountStr = "EUR 12.20"
var anotherAmmountNeg = journal.Ammount{
	Commodity: anotherAmmount.Commodity,
	Quantity:  anotherAmmount.Quantity.Neg(),
}
var anotherAmmountNegStr = "EUR -12.20"

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
			name: "OnPostingAccountDone from context",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.state.InputMetadata.SetPostingAccountText("BAR")
				c.state.InputMetadata.SetSelectedPostingAccount("FOO")

				c.controller.OnPostingAccountDone(input.Context)

				assert.Equal(t, statemod.InputPostingAmmount, c.state.CurrentPhase())
				posting := c.state.JournalEntryInput.CurrentPosting()
				account, _ := posting.GetAccount()
				assert.Equal(t, "FOO", account)
			},
		},
		{
			name: "OnPostingAccountDone from context no account",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.state.InputMetadata.SetSelectedPostingAccount("")
				c.state.InputMetadata.SetPostingAccountText("BAR")

				c.controller.OnPostingAccountDone(input.Context)

				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())
			},
		},
		{
			name: "OnPostingAccountDone from input",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.state.InputMetadata.SetSelectedPostingAccount("BAR")
				c.state.InputMetadata.SetPostingAccountText("FOO")

				c.controller.OnPostingAccountDone(input.Input)

				assert.Equal(t, statemod.InputPostingAmmount, c.state.CurrentPhase())
				posting := c.state.JournalEntryInput.CurrentPosting()
				account, _ := posting.GetAccount()
				assert.Equal(t, "FOO", account)
			},
		},
		{
			name: "OnPostingAccountDone from input no account",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.state.SetPhase(statemod.InputPostingAccount)
				c.state.InputMetadata.SetSelectedPostingAccount("BAR")
				c.state.InputMetadata.SetPostingAccountText("")

				c.controller.OnPostingAccountDone(input.Input)

				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())
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
				c.state.InputMetadata.SetPostingAccountText("BAR")
				c.state.InputMetadata.SetSelectedPostingAccount("FOO")
				c.controller.OnPostingAccountDone(input.Context)
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
			name: "After two unbalanced postings dont advance to confirmation",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("FOO")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anAmmountNegStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountNegStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Should still be on entering postings
				assert.Equal(t, statemod.InputPostingAccount, c.state.CurrentPhase())

				// First posting
				firstPosting, _ := c.state.JournalEntryInput.GetPosting(0)
				firstPostingAccount, _ := firstPosting.GetAccount()
				assert.Equal(t, "FOO", firstPostingAccount)
				firstPostingAmmount, _ := firstPosting.GetAmmount()
				assert.Equal(t, anAmmountNeg, firstPostingAmmount)

				// Second posting
				secondPosting, _ := c.state.JournalEntryInput.GetPosting(1)
				secondPostingAccount, _ := secondPosting.GetAccount()
				assert.Equal(t, "BAR", secondPostingAccount)
				secondPostingAmmount, _ := secondPosting.GetAmmount()
				assert.Equal(t, anotherAmmountNeg, secondPostingAmmount)
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
				c.controller.OnPostingAccountDone(input.Input)
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
			name: "Must write to ",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountNegStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Should have 2 filled postings and be on confirmation page
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				assert.Equal(t, statemod.Confirmation, c.state.CurrentPhase())
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				acc, _ := lastPosting.GetAccount()
				assert.Equal(t, "BAR2", acc)
				ammount, _ := lastPosting.GetAmmount()
				assert.Equal(t, anotherAmmountNeg, ammount)

				// Confirms submission
				c.controller.OnInputConfirmation()

				// Must have written to output
				wrote := c.bytesBuffer.String()
				assert.Contains(t, wrote, "BAR2")
				assert.Contains(t, wrote, "EUR -12.2")
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
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountNegStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Should have gone to confirmation page
				assert.Equal(t, statemod.Confirmation, c.state.CurrentPhase())

				// Should have 2 filled postings (empty one deleted)
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)

				// User decides to go back
				c.controller.OnInputRejection()

				// Should have 3 postings, 2 filled + 1 empty
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 3)
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				_, accFound := lastPosting.GetAccount()
				assert.False(t, accFound)
				_, ammountFound := lastPosting.GetAmmount()
				assert.False(t, ammountFound)
			},
		},
		{
			name: "User undo on confirmation page",
			opts: defaultOpts,
			run: func(t *testing.T, c *testcontext) {
				c.dateGuesser.EXPECT().Guess(gomock.Any())
				c.controller.OnDateChanged("2022-01-01")
				c.controller.OnDateDone()
				c.controller.OnDescriptionChanged("FOO")
				c.controller.OnDescriptionDone()
				c.controller.OnPostingAccountChanged("BAR")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountStr)
				c.controller.OnPostingAmmountDone(input.Input)
				c.controller.OnPostingAccountChanged("BAR2")
				c.controller.OnPostingAccountDone(input.Input)
				c.controller.OnPostingAmmountChanged(anotherAmmountNegStr)
				c.controller.OnPostingAmmountDone(input.Input)

				// Should have gone to confirmation page
				assert.Equal(t, statemod.Confirmation, c.state.CurrentPhase())

				// Should have 2 filled postings (empty one deleted)
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)

				// User decides to undo
				c.controller.OnUndo()

				// Should have 2 filled postings
				assert.Equal(t, c.state.JournalEntryInput.CountPostings(), 2)
				lastPosting := c.state.JournalEntryInput.CurrentPosting()
				acc, _ := lastPosting.GetAccount()
				assert.Equal(t, "BAR2", acc)
				ammount, _ := lastPosting.GetAmmount()
				assert.Equal(t, anotherAmmountNeg, ammount)
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
				c.controller.OnPostingAccountDone(input.Input)
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
