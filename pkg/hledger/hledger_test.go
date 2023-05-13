package hledger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
}
