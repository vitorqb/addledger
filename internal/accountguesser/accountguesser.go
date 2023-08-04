package accountguesser

import (
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/stringmatcher"
)

//go:generate $MOCKGEN --source=accountguesser.go --destination=../../mocks/accountguesser/accountguesser_mock.go

// TransactionHistory represents a history of transactions.
type TransactionHistory []journal.Transaction

// IAccountGuesser is an interface for an AccountGuesser, whose goal is to
// guess which account the user wants to input.
type IAccountGuesser interface {
	Guess() (guess journal.Account, success bool)
}

var _ IAccountGuesser = &DescriptionMatchAccountGuesser{}

// DescriptionMatchAccountGuesser matches a transaction from the transaction history that is
// similar to a description, and returns the suggested account based on the matched transaction.
// `inputPostings` must be the postings the user has already entered for the new transaction.
type DescriptionMatchAccountGuesser struct {
	matcher            stringmatcher.IStringMatcher
	transactionHistory TransactionHistory
	inputPostings      []journal.Posting
	description        string
}

func (ag *DescriptionMatchAccountGuesser) SetTransactionHistory(x TransactionHistory) {
	ag.transactionHistory = x
}

func (ag *DescriptionMatchAccountGuesser) SetInputPostings(x []journal.Posting) {
	ag.inputPostings = x
}

func (ag *DescriptionMatchAccountGuesser) SetDescription(x string) {
	ag.description = x
}

// Guess implements IAccountGuesser.
func (ag *DescriptionMatchAccountGuesser) Guess() (guess journal.Account, success bool) {
	matchedTransaction := journal.Transaction{}
	minDistance := 15

	// Finds matching transaction
	for _, transaction := range ag.transactionHistory {
		distance := ag.matcher.Distance(ag.description, transaction.Description)

		// We found a transaction with better score
		if distance < minDistance {
			minDistance = distance
			matchedTransaction = transaction
			continue
		}

		// We found a transaction with same score that is more recent.
		if distance < 20 && distance == minDistance {
			if transaction.Date.After(matchedTransaction.Date) {
				matchedTransaction = transaction
				continue
			}
		}
	}

	// If we had no match, return
	if minDistance >= 20 {
		return "", false
	}

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

// DescriptionMatchOption contains the options for a DescriptionMatchAccountGuesser
type DescriptionMatchOption struct {
	StringMatcher stringmatcher.IStringMatcher
}

// NewDescriptionMatchAccountGuesser returns a new implementation of AccountGuesser
func NewDescriptionMatchAccountGuesser(options DescriptionMatchOption) (*DescriptionMatchAccountGuesser, error) {
	var err error
	if options.StringMatcher == nil {
		options.StringMatcher, err = stringmatcher.New(&stringmatcher.Options{})
		if err != nil {
			return nil, err
		}
	}
	return &DescriptionMatchAccountGuesser{
		options.StringMatcher,
		TransactionHistory{},
		[]journal.Posting{},
		"",
	}, nil
}
