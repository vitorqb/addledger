package hledger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/journal"
	tu "github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/pkg/hledger"
)

// from testdata/fake_hledger.sh
var expectedAccounts = []string{
	"assets:bank:current:bnext",
	"assets:bank:savings:itau",
	"assets:cash",
	"assets:other",
	"expenses:bank-fees",
	"expenses:trips-and-travels",
	"expenses:unknown",
	"expenses:urban-transportation:public",
	"expenses:urban-transportation:taxi-uber-others",
	"initial-balance",
	"liabilities:credit-cards:amex",
	"liabilities:other",
	"revenues:earned-interests",
	"revenues:salary",
}

// from testdata/fake_hledger.sh
var expectedTransactions = []journal.Transaction{
	{Description: "Supermarket"},
	{Description: "Bar"},
}

func TestClient(t *testing.T) {
	t.Run("Accounts (no ledger file)", func(t *testing.T) {
		client := NewClient(tu.TestDataPath(t, "fake_hledger.sh"), "")
		accounts, err := client.Accounts()
		assert.NoError(t, err)
		assert.Equal(t, expectedAccounts, accounts)
	})
	t.Run("Accounts (ledger file)", func(t *testing.T) {
		client := NewClient(tu.TestDataPath(t, "fake_hledger.sh"), "foo")
		accounts, err := client.Accounts()
		assert.NoError(t, err)
		assert.Equal(t, expectedAccounts, accounts)
	})
	t.Run("Transactions (ledger file)", func(t *testing.T) {
		client := NewClient(tu.TestDataPath(t, "fake_hledger.sh"), "foo")
		transactions, err := client.Transactions()
		assert.NoError(t, err)
		assert.Equal(t, expectedTransactions, transactions)
	})
}
