// ammountguesser package is responsible for guessing which is the ammount
// an user wants to input for a given posting.
package ammountguesser

import (
	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/userinput"
)

//go:generate $MOCKGEN --source=ammountguesser.go --destination=../../mocks/ammountguesser/ammountguesser_mock.go

// Inputs are the inputs that are used for guessing.
type Inputs struct {
	// UserInput is the text the user has inputted in the ammount field. E.g. "EUR 12.20"
	UserInput string

	// PostingsData are the postings that have been inputted so far by the user,
	// in the current transaction.
	PostingsData []*state.PostingData

	// StatementEntry is the statement entry that has been loaded and is being used
	// for the current journal entry.
	StatementEntry finance.StatementEntry

	// MatchingTransactions are the transactions that match the current user input.
	MatchingTransactions []journal.Transaction
}

// IAmmountGuesser is a strategy for guessing the ammount an user may want for an journal entry.
type IAmmountGuesser interface {
	Guess(inputs Inputs) (guess finance.Ammount, success bool)
}

// AmmountGuesser is the default implementation of IAmmountGuesser.
type AmmountGuesser struct{}

// Guess implements IAmmountGuesser.
func (*AmmountGuesser) Guess(inputs Inputs) (guess finance.Ammount, success bool) {
	// If user entered an ammount, use it
	if ammountFromUserInput, err := userinput.TextToAmmount(inputs.UserInput); err == nil {
		if ammountFromUserInput.Commodity == "" {
			ammountFromUserInput.Commodity = DefaultCommodity
		}
		return ammountFromUserInput, true
	}

	// If we have pending balance, use it
	nonEmptyPostingData := selectNonEmptyPostingData(inputs.PostingsData)
	postings, _ := userinput.PostingsFromData(nonEmptyPostingData)
	balance := journal.PostingsBalance(postings)
	// The if below checks (a) we have a single currency with (b) non-zero balance
	if len(balance) == 1 && !balance[0].Quantity.IsZero() {
		return balance[0].InvertSign(), true
	}

	// If we have a statement entry, use it
	if inputs.StatementEntry.Ammount.Quantity.Abs().GreaterThan(decimal.Zero) {
		return inputs.StatementEntry.Ammount.InvertSign(), true
	}

	// If we have a matching transaction, use it.
	if len(inputs.MatchingTransactions) > 0 {
		return inputs.MatchingTransactions[0].Posting[0].Ammount, true
	}

	return DefaultGuess, true
}

var _ IAmmountGuesser = &AmmountGuesser{}

func New() *AmmountGuesser { return &AmmountGuesser{} }

var DefaultCommodity string = "EUR"

// TODO Get rid of this!
var DefaultGuess = finance.Ammount{
	Commodity: "EUR",
	Quantity:  decimal.New(1220, -2),
}

func selectNonEmptyPostingData(postingsData []*state.PostingData) []*state.PostingData {
	var nonEmpty []*state.PostingData
	for _, data := range postingsData {
		_, found := data.Ammount.Get()
		if found {
			nonEmpty = append(nonEmpty, data)
		}
	}
	return nonEmpty
}
