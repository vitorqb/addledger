// injector package is responsible for injecting dependencies on runtime.
package injector

import (
	"github.com/vitorqb/addledger/internal/ammountguesser"
	configmod "github.com/vitorqb/addledger/internal/config"
	statemod "github.com/vitorqb/addledger/internal/state"
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

	// subscribes to changes
	state.AddOnChangeHook(func() {

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

func State(hledgerClient hledger.IClient) (*statemod.State, error) {
	// Initializes a new state
	state := statemod.InitialState()

	// load accounts
	accounts, err := hledgerClient.Accounts()
	if err != nil {
		return &statemod.State{}, err
	}
	state.JournalMetadata.SetAccounts(accounts)

	// load transactions
	postings, err := hledgerClient.Transactions()
	if err != nil {
		return &statemod.State{}, err
	}
	state.JournalMetadata.SetTransactions(postings)

	return state, nil
}
