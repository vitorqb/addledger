// injector package is responsible for injecting dependencies on runtime.
package injector

import (
	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	configmod "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/metaloader"
	"github.com/vitorqb/addledger/internal/printer"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/stringmatcher"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
	"github.com/vitorqb/addledger/pkg/hledger"
)

// HledgerClient injects a new client for HLedger.
func HledgerClient(config *configmod.Config) hledger.IClient {
	return hledger.NewClient(config.HLedgerExecutable, config.LedgerFile)
}

// AmmountGuesserEngine instantiates a new guesser engine for ammounts and
// links it to the state
func AmmountGuesserEngine(state *statemod.State) ammountguesser.IEngine {
	// starts an engine
	ammountGuesserEngine := ammountguesser.NewEngine()

	// sets initial guess
	text := state.InputMetadata.GetPostingAmmountText()
	ammountGuesserEngine.SetUserInputText(text)
	if guess, success := ammountGuesserEngine.Guess(); success {
		state.InputMetadata.SetPostingAmmountGuess(guess)
	}

	busy := false
	// subscribes to changes
	state.AddOnChangeHook(func() {
		if busy {
			return
		}
		busy = true
		defer func() { busy = false }()

		// sync matching transactions
		matchingTransactions := state.InputMetadata.MatchingTransactions()
		ammountGuesserEngine.SetMatchingTransactions(matchingTransactions)

		// sync input text
		newText := state.InputMetadata.GetPostingAmmountText()
		ammountGuesserEngine.SetUserInputText(newText)

		// sync existing postings
		newPostings := state.JournalEntryInput.GetPostings()
		ammountGuesserEngine.SetPostingInputs(newPostings)

		oldGuess, oldGuessFound := state.InputMetadata.GetPostingAmmountGuess()
		guess, success := ammountGuesserEngine.Guess()
		if success {
			if !guess.Equal(oldGuess) {
				state.InputMetadata.SetPostingAmmountGuess(guess)
			}
		} else {
			if oldGuessFound {
				state.InputMetadata.ClearPostingAmmountGuess()
			}
		}

	})

	// returns
	return ammountGuesserEngine
}

func DateGuesser() (dateguesser.IDateGuesser, error) {
	return dateguesser.New()
}

func State(hledgerClient hledger.IClient) (*statemod.State, error) {
	return statemod.InitialState(), nil
}

func MetaLoader(state *statemod.State, hledgerClient hledger.IClient) (*metaloader.MetaLoader, error) {
	return metaloader.New(state, hledgerClient)
}

// DescriptionMatchAccountGuesser instantiates a new DescriptionMatchAccountGuesser and syncs it with
// the state.
func DescriptionMatchAccountGuesser(state *statemod.State) (*accountguesser.MatchedTransactionsGuesser, error) {

	// Creates a new Description Guesser
	accountGuesser, err := accountguesser.NewMatchedTransactionsAccountGuesser()
	if err != nil {
		return nil, err
	}

	// Function that syncs the state with the internal AccountGuesser state.
	syncWithState := func() {
		matchedTransactions := state.InputMetadata.MatchingTransactions()
		accountGuesser.SetMatchedTransactions(matchedTransactions)

		var postings []journal.Posting
		inputPostings := state.JournalEntryInput.GetPostings()
		for _, inputPostings := range inputPostings {
			if inputPostings.IsComplete() {
				postings = append(postings, inputPostings.ToPosting())
			}
		}
		accountGuesser.SetInputPostings(postings)
	}

	// Runs a first sync
	syncWithState()

	// Runs a sync everytime the state changes
	state.AddOnChangeHook(syncWithState)

	// Returns
	return accountGuesser, nil
}

func LastTransactionAccountGuesser(state *statemod.State) (*accountguesser.LastTransactionAccountGuesser, error) {
	// Creates a new LastTransactionAccountGuesser
	accountGuesser, err := accountguesser.NewLastTransactionAccountGuesser()
	if err != nil {
		return nil, err
	}

	// Function that syncs the state with the internal AccountGuesser state.
	sync := func() {
		transactionHistory := state.JournalMetadata.Transactions()
		accountGuesser.SetTransactionHistory(transactionHistory)
	}

	// Runs first sync
	sync()

	// Run sync on state update
	state.AddOnChangeHook(sync)

	return accountGuesser, nil
}

func AccountGuesser(state *statemod.State) (accountguesser.IAccountGuesser, error) {
	// Returns a composite of DescriptionMatch and LastTransaction
	descriptionMatchAccountGuesser, err := DescriptionMatchAccountGuesser(state)
	if err != nil {
		return nil, err
	}
	lastTransactionAccountGuesser, err := LastTransactionAccountGuesser(state)
	if err != nil {
		return nil, err
	}
	return accountguesser.NewCompositeAccountGuesser(descriptionMatchAccountGuesser, lastTransactionAccountGuesser)
}

func Printer(config configmod.PrinterConfig) (printer.IPrinter, error) {
	return printer.New(config.NumLineBreaksBefore, config.NumLineBreaksAfter), nil
}

func TransactionMatcher(state *statemod.State) (transactionmatcher.ITransactionMatcher, error) {
	// We could inject a stringmatcher here if we ever want to make it configurable.
	stringMatcher, err := stringmatcher.New(&stringmatcher.Options{})
	if err != nil {
		return nil, err
	}

	transactionMatcher := transactionmatcher.New(stringMatcher)

	// Ensure the state is updated when the matched transaction changes.
	busy := false
	state.AddOnChangeHook(func() {
		if !busy {
			busy = true
			defer func() { busy = false }()
			descriptionInput, found := state.JournalEntryInput.GetDescription()
			if !found {
				return
			}
			transactionMatcher.SetDescriptionInput(descriptionInput)
			transactionHistory := state.JournalMetadata.Transactions()
			transactionMatcher.SetTransactionHistory(transactionHistory)
			matchingTransactions := transactionMatcher.Match()
			state.InputMetadata.SetMatchingTransactions(matchingTransactions)
		}
	})

	return transactionMatcher, nil
}
