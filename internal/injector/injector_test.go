package injector_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/injector"
	. "github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/statementloader"
	"github.com/vitorqb/addledger/internal/testutils"
	hledger_mock "github.com/vitorqb/addledger/mocks/hledger"
)

func TestStateAndMetaLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	transactions := []journal.Transaction{{Description: "FOO"}, {Description: "Bar"}}
	hledgerClient := hledger_mock.NewMockIClient(ctrl)
	hledgerClient.EXPECT().Accounts().Return([]journal.Account{"FOO"}, nil)
	hledgerClient.EXPECT().Transactions().Return(transactions, nil)

	state, err := State(hledgerClient)
	assert.Nil(t, err)

	metaLoader, err := MetaLoader(state, hledgerClient)
	assert.Nil(t, err)
	err = metaLoader.LoadAccounts()
	assert.Nil(t, err)
	err = metaLoader.LoadTransactions()
	assert.Nil(t, err)
	assert.Equal(t, []journal.Account{"FOO"}, state.JournalMetadata.Accounts())
	assert.Equal(t, transactions, state.JournalMetadata.Transactions())
}

func TestPrinter(t *testing.T) {
	config := config.PrinterConfig{NumLineBreaksBefore: 2, NumLineBreaksAfter: 3}
	printer, err := injector.Printer(config)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	printer.Print(&buf, *testutils.Transaction_1(t))
	expectedPrint := "\n\n1993-11-23 Description1\n    ACC1    EUR 12.2\n    ACC2    EUR -12.2\n\n\n"
	assert.Equal(t, expectedPrint, buf.String())
}

func TestTransactionMatcher(t *testing.T) {
	matcher, err := injector.TransactionMatcher()
	assert.Nil(t, err)
	transactions := []journal.Transaction{{Description: "test"}, {Description: "INVALID"}}
	expectedMatchedTransactions := []journal.Transaction{{Description: "test"}}
	matcher.SetTransactionHistory(transactions)
	matcher.SetDescriptionInput("test")
	matchedTransactions := matcher.Match()
	assert.Equal(t, expectedMatchedTransactions, matchedTransactions)
}

func TestCSVStatementLoaderOptions(t *testing.T) {
	type testcase struct {
		name            string
		config          config.CSVStatementLoaderConfig
		expectedOptions []statementloader.CSVLoaderOption
	}
	testcases := []testcase{
		{
			name: "empty",
			config: config.CSVStatementLoaderConfig{
				DateFieldIndex:        -1,
				DescriptionFieldIndex: -1,
				AccountFieldIndex:     -1,
				AmmountFieldIndex:     -1,
			},
			expectedOptions: []statementloader.CSVLoaderOption{
				statementloader.WithCSVLoaderMapping([]statementloader.CSVColumnMapping{}),
			},
		},
		{
			name: "full",
			config: config.CSVStatementLoaderConfig{
				Separator:             ";",
				Account:               "acc",
				Commodity:             "com",
				DateFieldIndex:        0,
				DateFormat:            "01/02/2006",
				DescriptionFieldIndex: 1,
				AccountFieldIndex:     2,
				AmmountFieldIndex:     3,
			},
			expectedOptions: []statementloader.CSVLoaderOption{
				statementloader.WithCSVSeparator(';'),
				statementloader.WithCSVLoaderAccountName("acc"),
				statementloader.WithCSVLoaderDefaultCommodity("com"),
				statementloader.WithCSVLoaderMapping([]statementloader.CSVColumnMapping{
					{Column: 0, Importer: statementloader.DateImporter{Format: "01/02/2006"}},
					{Column: 1, Importer: statementloader.DescriptionImporter{}},
					{Column: 2, Importer: statementloader.AccountImporter{}},
					{Column: 3, Importer: statementloader.AmmountImporter{}},
				}),
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actualConfig := statementloader.CSVLoaderConfig{}
			expectedConfig := statementloader.CSVLoaderConfig{}
			options, err := CSVStatementLoaderOptions(testcase.config)
			assert.Nil(t, err)
			for _, option := range options {
				option(&actualConfig)
			}
			for _, option := range testcase.expectedOptions {
				option(&expectedConfig)
			}
			assert.Equal(t, expectedConfig, actualConfig)
		})
	}
}
