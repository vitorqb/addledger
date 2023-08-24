package transactionmatcher

import (
	"sort"

	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/stringmatcher"
)

// match is a transaction and its distance from the inputs.
type match struct {
	transaction journal.Transaction
	distance    int
}

// matchList is a list of matches.
type matchList []match

// Implement sort.Interface for matchList
func (ml matchList) Len() int           { return len(ml) }
func (ml matchList) Less(i, j int) bool { return ml[i].distance < ml[j].distance }
func (ml matchList) Swap(i, j int)      { ml[i], ml[j] = ml[j], ml[i] }

// ITransactionMatcher implements the logic to match transactions
// based on user input and a transaction history.
type ITransactionMatcher interface {
	// SetDescriptionInput sets the description inputted by the user.
	SetDescriptionInput(x string)

	// SetTransactionHistory sets the transaction history.
	SetTransactionHistory(x []journal.Transaction)

	// Match return the best matches from the transaction history
	// for the current description input.
	Match() []journal.Transaction
}

// TransactionMatcher implements ITransactionMatcher.
type TransactionMatcher struct {
	descriptionInput   string
	transactionHistory []journal.Transaction
	stringMatcher      stringmatcher.IStringMatcher
}

// SetDescriptionInput sets the description inputted by the user.
func (tm *TransactionMatcher) SetDescriptionInput(x string) {
	tm.descriptionInput = x
}

// SetTransactionHistory sets the transaction history.
func (tm *TransactionMatcher) SetTransactionHistory(x []journal.Transaction) {
	tm.transactionHistory = x
}

// Match return the best matches from the transaction history
// for the current description input.
func (tm *TransactionMatcher) Match() []journal.Transaction {
	var matches []match
	for _, transaction := range tm.transactionHistory {
		descriptionDistance := tm.stringMatcher.Distance(tm.descriptionInput, transaction.Description)
		// !!!! TODO Make 6 a configurable value
		if descriptionDistance <= 6 {
			matches = append(matches, match{transaction, descriptionDistance})
		}
	}
	sort.Sort(matchList(matches))
	// !!!! TODO Make 20 a configurable value
	if len(matches) > 20 {
		matches = matches[:20]
	}
	transactions := []journal.Transaction{}
	for _, match := range matches {
		transactions = append(transactions, match.transaction)
	}
	return transactions
}

var _ ITransactionMatcher = &TransactionMatcher{}

// New returns a new TransactionMatcher.
func New(stringMatcher stringmatcher.IStringMatcher) *TransactionMatcher {
	descriptionInput := ""
	transactionHistory := []journal.Transaction{}
	return &TransactionMatcher{descriptionInput, transactionHistory, stringMatcher}
}
