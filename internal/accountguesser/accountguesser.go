package accountguesser

import (
	"github.com/vitorqb/addledger/internal/journal"
)

//go:generate $MOCKGEN --source=accountguesser.go --destination=../../mocks/accountguesser/accountguesser_mock.go

// MatchedTransactions represents a collection of matched transactions.
type MatchedTransactions []journal.Transaction

// IAccountGuesser is an interface for an AccountGuesser, whose goal is to
// guess which account the user wants to input.
type IAccountGuesser interface {
	Guess() (guess journal.Account, success bool)
}

var _ IAccountGuesser = &MatchedTransactionsGuesser{}

// MatchedTransactionsGuesser uses the matched transactions and
// returns the suggested account based on the best-matched
// transaction. `inputPostings` must be the postings the user has
// already entered for the new transaction.
type MatchedTransactionsGuesser struct {
	matchedTransactions MatchedTransactions
	inputPostings       []journal.Posting
	description         string
}

func (ag *MatchedTransactionsGuesser) SetMatchedTransactions(x MatchedTransactions) {
	ag.matchedTransactions = x
}

func (ag *MatchedTransactionsGuesser) SetInputPostings(x []journal.Posting) {
	ag.inputPostings = x
}

// Guess implements IAccountGuesser.
func (ag *MatchedTransactionsGuesser) Guess() (guess journal.Account, success bool) {
	if len(ag.matchedTransactions) == 0 {
		return "", false
	}
	matchedTransaction := ag.matchedTransactions[0]

	// We had a match, so find posting the user is entering
	desiredPostingIndex := len(ag.inputPostings)

	// If the user has already entered more posting than matched transaction,
	// we can't use it.
	if desiredPostingIndex >= len(matchedTransaction.Posting) {
		return "", false
	}

	// Otherwise get the account from the posting with same index.
	matchedPosting := matchedTransaction.Posting[desiredPostingIndex]
	return journal.Account(matchedPosting.Account), true
}

// NewMatchedTransactionsAccountGuesser returns a new implementation of AccountGuesser
func NewMatchedTransactionsAccountGuesser() (*MatchedTransactionsGuesser, error) {
	return &MatchedTransactionsGuesser{MatchedTransactions{}, []journal.Posting{}, ""}, nil
}

// LastTransactionAccountGuesser uses the last entered transaction to try
// to guess an account
type LastTransactionAccountGuesser struct {
	transactionHistory MatchedTransactions
}

var _ IAccountGuesser = &LastTransactionAccountGuesser{}

func (ag *LastTransactionAccountGuesser) SetTransactionHistory(x MatchedTransactions) {
	ag.transactionHistory = x
}

func (ag *LastTransactionAccountGuesser) Guess() (acc journal.Account, success bool) {
	historyLen := len(ag.transactionHistory)
	if historyLen == 0 {
		return "", false
	}
	lastTransaction := ag.transactionHistory[historyLen-1]
	if len(lastTransaction.Posting) == 0 {
		return "", false
	}
	firstPosting := lastTransaction.Posting[0]
	return journal.Account(firstPosting.Account), true
}

func NewLastTransactionAccountGuesser() (*LastTransactionAccountGuesser, error) {
	return &LastTransactionAccountGuesser{MatchedTransactions{}}, nil
}

// CompositeAccountGuesser composes N different account guessers
type CompositeAccountGuesser struct {
	composedGuessers []IAccountGuesser
}

var _ IAccountGuesser = &CompositeAccountGuesser{}

// Guess implements IAccountGuesser.
func (ag *CompositeAccountGuesser) Guess() (guess journal.Account, success bool) {
	for _, composedGuesser := range ag.composedGuessers {
		if guess, ok := composedGuesser.Guess(); ok {
			return guess, ok
		}
	}
	return journal.Account(""), false
}

func NewCompositeAccountGuesser(accGuessers ...IAccountGuesser) (*CompositeAccountGuesser, error) {
	return &CompositeAccountGuesser{accGuessers}, nil
}
