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

func CSVStatementLoaderOptions(config configmod.CSVStatementLoaderConfig) ([]statementreader.CSVReaderOption, error) {
	options := []statementreader.CSVReaderOption{}
	if acc := config.Account; acc != "" {
		options = append(options, statementreader.WithCSVReaderAccountName(acc))
	}
	if comm := config.Commodity; comm != "" {
		options = append(options, statementreader.WithCSVReaderDefaultCommodity(comm))
	}
	if sep := config.Separator; sep != "" {
		if len(sep) != 1 {
			return nil, fmt.Errorf("invalid csv separator: %s", sep)
		}
		options = append(options, statementreader.WithCSVSeparator([]rune(sep)[0]))
	}
	mapping := []statementreader.CSVColumnMapping{}
	if idate := config.DateFieldIndex; idate != -1 {
		importer := statementreader.DateImporter{Format: config.DateFormat}
		mapping = append(mapping, statementreader.CSVColumnMapping{Column: idate, Importer: importer})
	}
	if idescription := config.DescriptionFieldIndex; idescription != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: idescription, Importer: statementreader.DescriptionImporter{},
		})
	}
	if iaccount := config.AccountFieldIndex; iaccount != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: iaccount, Importer: statementreader.AccountImporter{},
		})
	}
	if iammount := config.AmmountFieldIndex; iammount != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: iammount, Importer: statementreader.AmmountImporter{},
		})
	}
	options = append(options, statementreader.WithCSVLoaderMapping(mapping))
	return options, nil
}

func CSVStatementLoader(config configmod.CSVStatementLoaderConfig) (*statementreader.CSVStatementReader, error) {
	options, err := CSVStatementLoaderOptions(config)
	if err != nil {
		return nil, err
	}
	return statementreader.NewCSVLoader(options...), nil
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
