// ammountguesser package is responsible for guessing which is the ammount
// an user wants to input for a given posting.
package ammountguesser

import (
	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/statementloader"
)

//go:generate $MOCKGEN --source=ammountguesser.go --destination=../../mocks/ammountguesser/ammountguesser_mock.go

// Inputs are the inputs that are used for guessing.
type Inputs struct {
	// UserInput is the text the user has inputted in the ammount field. E.g. "EUR 12.20"
	UserInput string

	// PostingInputs are the postings that have been inputted so far by the user,
	// in the current journal entry.
	PostingInputs []*input.PostingInput

	// StatementEntry is the statement entry that has been loaded and is being used
	// for the current journal entry.
	StatementEntry statementloader.StatementEntry

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
	if ammountFromUserInput, err := input.TextToAmmount(inputs.UserInput); err == nil {
		if ammountFromUserInput.Commodity == "" {
			ammountFromUserInput.Commodity = DefaultCommodity
		}
		return ammountFromUserInput, true
	}

	// If we have pending balance, use it
	for {
		nonEmptyPostingInputs := selectNonEmptyPostingInputs(inputs.PostingInputs)
		if len(nonEmptyPostingInputs) < 1 {
			break
		}

		// Calculate pending balance
		var success bool = true
		balance := decimal.Zero
		commodity := ""
		for _, posting := range nonEmptyPostingInputs[1:] {
			ammount, _ := posting.GetAmmount()
			// Multiple commodities -> stop
			if commodity != "" && ammount.Commodity != commodity {
				success = false
				break
			}
			commodity = ammount.Commodity
			balance = balance.Sub(ammount.Quantity)
		}
		if !success {
			break
		} else {
			return finance.Ammount{Commodity: commodity, Quantity: balance}, true
		}
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

func selectNonEmptyPostingInputs(postingInputs []*input.PostingInput) []*input.PostingInput {
	var nonEmptyPostingInputs []*input.PostingInput
	for _, input := range postingInputs {
		_, found := input.GetAmmount()
		if found {
			nonEmptyPostingInputs = append(nonEmptyPostingInputs, input)
		}
	}
	return nonEmptyPostingInputs
}
