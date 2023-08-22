package metaloader

import (
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/pkg/hledger"
)

//go:generate $MOCKGEN --source=metaloader.go --destination=../../mocks/metaloader/metaloader_mock.go

// IMetaLoader loads metadata from the journal to a state.
type IMetaLoader interface {
	LoadTransactions() error
	LoadAccounts() error
}

// MetaLoader implements iMetaLoader
type MetaLoader struct {
	state         *state.State
	hledgerClient hledger.IClient
}

var _ IMetaLoader = &MetaLoader{}

// LoadAccounts implements IMetaLoader.
func (ml *MetaLoader) LoadAccounts() error {
	accounts, err := ml.hledgerClient.Accounts()
	if err != nil {
		return err
	}
	ml.state.JournalMetadata.SetAccounts(accounts)
	return nil
}

// LoadTransactions implements IMetaLoader.
func (ml *MetaLoader) LoadTransactions() error {
	postings, err := ml.hledgerClient.Transactions()
	if err != nil {
		return err
	}
	ml.state.JournalMetadata.SetTransactions(postings)
	return nil
}

// New returns a new instance of MetaLoader
func New(state *state.State, hledgerClient hledger.IClient) (*MetaLoader, error) {
	return &MetaLoader{state, hledgerClient}, nil
}
