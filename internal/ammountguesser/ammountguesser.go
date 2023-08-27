package ammountguesser

import (
	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
)

//go:generate $MOCKGEN --source=ammountguesser.go --destination=../../mocks/ammountguesser/ammountguesser_mock.go

type IEngine interface {
	// SetUserInputText must be called to set what's the text the user
	// has inputted.
	SetUserInputText(x string)

	// SetPostingInputs set's all posting inputs that have been inputted so far.
	SetPostingInputs(x []*input.PostingInput)

	// SetMatchingTransactions set's all transactions that match the current
	// user input.
	SetMatchingTransactions(x []journal.Transaction)

	// Guess returns the guess for the current state. If can't
	// guess, success is false.
	Guess() (guess journal.Ammount, success bool)
}

var DefaultCommodity string = "EUR"

// TODO Get rid of this!
var DefaultGuess = journal.Ammount{
	Commodity: "EUR",
	Quantity:  decimal.New(1220, -2),
}

type Engine struct {
	userInput            string
	postingInputs        []*input.PostingInput
	matchingTransactions []journal.Transaction
}

var _ IEngine = &Engine{}

func NewEngine() *Engine { return &Engine{} }

func (e *Engine) Guess() (guess journal.Ammount, success bool) {

	// If user entered an ammount, use it
	if ammountFromUserInput, err := input.TextToAmmount(e.userInput); err == nil {
		// If user didn't enter commodity, use the default
		if ammountFromUserInput.Commodity == "" {
			ammountFromUserInput.Commodity = DefaultCommodity
		}
		return ammountFromUserInput, true
	}

	// If we have pending balance, use it
	for {
		var nonEmptyPostingInputs []*input.PostingInput
		for _, input := range e.postingInputs {
			_, found := input.GetAmmount()
			if found {
				nonEmptyPostingInputs = append(nonEmptyPostingInputs, input)
			}
		}

		// If we have 0 non-empty inputs, can't guess
		if len(nonEmptyPostingInputs) < 1 {
			break
		}

		// Calculate pending balance
		var success bool = true
		firstAmmount, _ := nonEmptyPostingInputs[0].GetAmmount()
		balance := firstAmmount.Quantity.Mul(decimal.NewFromInt(-1))
		for _, posting := range nonEmptyPostingInputs[1:] {
			ammount, _ := posting.GetAmmount()
			// Multiple commodities -> stop
			if ammount.Commodity != firstAmmount.Commodity {
				success = false
				break
			}
			balance = balance.Sub(ammount.Quantity)
		}
		if !success {
			break
		} else {
			return journal.Ammount{
				Commodity: firstAmmount.Commodity,
				Quantity:  balance,
			}, true
		}
	}

	// If we have a matching transaction, use it.
	if len(e.matchingTransactions) > 0 {
		transaction := e.matchingTransactions[0]
		return transaction.Posting[0].Ammount, true
	}

	return DefaultGuess, true
}

func (e *Engine) SetUserInputText(x string)                       { e.userInput = x }
func (e *Engine) SetPostingInputs(x []*input.PostingInput)        { e.postingInputs = x }
func (e *Engine) SetMatchingTransactions(x []journal.Transaction) { e.matchingTransactions = x }
