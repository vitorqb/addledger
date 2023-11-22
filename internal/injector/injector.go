// injector package is responsible for injecting dependencies on runtime.
package injector

import (
	"fmt"

	"github.com/vitorqb/addledger/internal/accountguesser"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	configmod "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/metaloader"
	"github.com/vitorqb/addledger/internal/printer"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
	"github.com/vitorqb/addledger/internal/stringmatcher"
	"github.com/vitorqb/addledger/internal/transactionmatcher"
	"github.com/vitorqb/addledger/pkg/hledger"
)

// HledgerClient injects a new client for HLedger.
func HledgerClient(config *configmod.Config) hledger.IClient {
	return hledger.NewClient(config.HLedgerExecutable, config.LedgerFile)
}

// AmmountGuesserEngine instantiates a new guesser engine for ammount.
func AmmountGuesserEngine() ammountguesser.IEngine {
	// starts an engine
	ammountGuesserEngine := ammountguesser.NewEngine()

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
		completePostings := state.JournalEntryInput.GetCompletePostings()
		accountGuesser.SetInputPostings(completePostings)
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

func StatementAccountGuesser(state *statemod.State) (accountguesser.IAccountGuesser, error) {
	// Creates a new StatementAccountGuesser
	accountGuesser, err := accountguesser.NewStatementAccountGuesser()
	if err != nil {
		return nil, err
	}

	// Function that syncs the state with the internal AccountGuesser state.
	sync := func() {
		statementEntry, _ := state.CurrentStatementEntry()
		accountGuesser.SetStatementEntry(statementEntry)
		completedPostings := state.JournalEntryInput.GetCompletePostings()
		accountGuesser.SetInputPostings(completedPostings)
	}

	// Runs first sync
	sync()

	// Run sync on state update
	state.AddOnChangeHook(sync)

	return accountGuesser, nil
}

func AccountGuesser(state *statemod.State) (accountguesser.IAccountGuesser, error) {
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

func CSVStatementLoaderOptions(config configmod.CSVStatementLoaderConfig) ([]statementloader.CSVLoaderOption, error) {
	options := []statementloader.CSVLoaderOption{}
	if acc := config.Account; acc != "" {
		options = append(options, statementloader.WithCSVLoaderAccountName(acc))
	}
	if comm := config.Commodity; comm != "" {
		options = append(options, statementloader.WithCSVLoaderDefaultCommodity(comm))
	}
	if sep := config.Separator; sep != "" {
		if len(sep) != 1 {
			return nil, fmt.Errorf("invalid csv separator: %s", sep)
		}
		options = append(options, statementloader.WithCSVSeparator([]rune(sep)[0]))
	}
	mapping := []statementloader.CSVColumnMapping{}
	if idate := config.DateFieldIndex; idate != -1 {
		importer := statementloader.DateImporter{Format: config.DateFormat}
		mapping = append(mapping, statementloader.CSVColumnMapping{Column: idate, Importer: importer})
	}
	if idescription := config.DescriptionFieldIndex; idescription != -1 {
		mapping = append(mapping, statementloader.CSVColumnMapping{
			Column: idescription, Importer: statementloader.DescriptionImporter{},
		})
	}
	if iaccount := config.AccountFieldIndex; iaccount != -1 {
		mapping = append(mapping, statementloader.CSVColumnMapping{
			Column: iaccount, Importer: statementloader.AccountImporter{},
		})
	}
	if iammount := config.AmmountFieldIndex; iammount != -1 {
		mapping = append(mapping, statementloader.CSVColumnMapping{
			Column: iammount, Importer: statementloader.AmmountImporter{},
		})
	}
	options = append(options, statementloader.WithCSVLoaderMapping(mapping))
	return options, nil
}

func CSVStatementLoader(config configmod.CSVStatementLoaderConfig) (*statementloader.CSVLoader, error) {
	options, err := CSVStatementLoaderOptions(config)
	if err != nil {
		return nil, err
	}
	return statementloader.NewCSVLoader(options...), nil
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
