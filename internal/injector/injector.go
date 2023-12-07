// injector package is responsible for injecting dependencies on runtime.
package injector

import (
	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	configmod "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/metaloader"
	"github.com/vitorqb/addledger/internal/printer"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
	"github.com/vitorqb/addledger/internal/stringmatcher"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
	"github.com/vitorqb/addledger/pkg/hledger"
)

// HledgerClient injects a new client for HLedger.
func HledgerClient(config *configmod.Config) hledger.IClient {
	return hledger.NewClient(config.HLedgerExecutable, config.LedgerFile)
}

// AmmountGuesser instantiates a new guesser for ammount.
func AmmountGuesser() ammountguesser.IAmmountGuesser {
	return ammountguesser.New()
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
	return accountguesser.NewMatchedTransactionsAccountGuesser()
}

func LastTransactionAccountGuesser(state *statemod.State) (*accountguesser.LastTransactionAccountGuesser, error) {
	return accountguesser.NewLastTransactionAccountGuesser()
}

func StatementAccountGuesser(state *statemod.State) (accountguesser.AccountGuesser, error) {
	return accountguesser.NewStatementAccountGuesser()
}

func AccountGuesser(state *statemod.State) (accountguesser.AccountGuesser, error) {
	// Returns a composite of all account guessers
	statementAccountGuesser, err := StatementAccountGuesser(state)
	if err != nil {
		return nil, err
	}
	descriptionMatchAccountGuesser, err := DescriptionMatchAccountGuesser(state)
	if err != nil {
		return nil, err
	}
	lastTransactionAccountGuesser, err := LastTransactionAccountGuesser(state)
	if err != nil {
		return nil, err
	}
	return accountguesser.NewCompositeAccountGuesser(
		statementAccountGuesser,
		descriptionMatchAccountGuesser,
		lastTransactionAccountGuesser,
	)
}

func Printer(config configmod.PrinterConfig) (printer.IPrinter, error) {
	return printer.New(config.NumLineBreaksBefore, config.NumLineBreaksAfter), nil
}

func StatementReader() statementreader.IStatementReader {
	return statementreader.NewStatementReader()
}

func TransactionMatcher() (transactionmatcher.ITransactionMatcher, error) {
	// We could inject a stringmatcher here if we ever want to make it configurable.
	stringMatcher, err := stringmatcher.New(&stringmatcher.Options{})
	if err != nil {
		return nil, err
	}

	transactionMatcher := transactionmatcher.New(stringMatcher)

	return transactionMatcher, nil
}
