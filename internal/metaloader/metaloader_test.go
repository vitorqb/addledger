package metaloader_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	. "github.com/vitorqb/addledger/internal/metaloader"
	statemod "github.com/vitorqb/addledger/internal/state"
	hledger_mocks "github.com/vitorqb/addledger/mocks/hledger"
)

var accounts = []journal.Account{"assets:bank:current:bnext", "assets:bank:savings:itau"}
var transactions = []journal.Transaction{
	{
		Description: "Supermarket",
		Date:        time.Date(2018, 12, 1, 0, 0, 0, 0, time.UTC),
		Posting: []journal.Posting{
			{
				Account: "liabilities:other",
				Ammount: journal.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(-4000000, -5),
				},
			},
			{
				Account: "expenses:sports",
				Ammount: journal.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(4000000, -5),
				},
			},
		},
	},
}

func TestMetaLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := statemod.InitialState()
	hledgerClient := hledger_mocks.NewMockIClient(ctrl)
	hledgerClient.EXPECT().Accounts().Return(accounts, nil)
	hledgerClient.EXPECT().Transactions().Return(transactions, nil)
	metaLoader, err := New(state, hledgerClient)
	assert.Nil(t, err)
	err = metaLoader.LoadTransactions()
	assert.Nil(t, err)
	err = metaLoader.LoadAccounts()
	assert.Nil(t, err)
	assert.Equal(t, state.JournalMetadata.Accounts(), accounts)
	assert.Equal(t, state.JournalMetadata.Transactions(), transactions)
}
