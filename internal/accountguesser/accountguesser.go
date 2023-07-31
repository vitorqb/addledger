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
	Guess(
		transactionHistory TransactionHistory,
		inputPostings []journal.Posting,
		description string,
	) (guess journal.Account, success bool)
}

// AccountGuesser implements IAccountGuesser
type AccountGuesser struct {
	matcher stringmatcher.IStringMatcher
}

var _ IAccountGuesser = &AccountGuesser{}

// Guess implements IAccountGuesser.
func (a *AccountGuesser) Guess(
	transactionHistory TransactionHistory,
	inputPostings []journal.Posting,
	description string,
) (journal.Account, bool) {
	matchedTransaction := journal.Transaction{}
	minDistance := 15

	// Finds matching transaction
	for _, transaction := range transactionHistory {
		distance := a.matcher.Distance(description, transaction.Description)

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
	desiredPostingIndex := len(inputPostings)

	// If the user has already entered more posting than matched transaction,
	// we can't use it.
	if desiredPostingIndex >= len(matchedTransaction.Posting) {
		return "", false
	}

	// Otherwise get the account from the posting with same index.
	matchedPosting := matchedTransaction.Posting[desiredPostingIndex]
	return journal.Account(matchedPosting.Account), true
}

// Options contains the options for an AccountGuesser
type Options struct {
	StringMatcher stringmatcher.IStringMatcher
}

// New returns a new implementation of AccountGuesser
func New(options Options) (*AccountGuesser, error) {
	var err error
	if options.StringMatcher == nil {
		options.StringMatcher, err = stringmatcher.New(&stringmatcher.Options{})
		if err != nil {
			return nil, err
		}
	}
	return &AccountGuesser{options.StringMatcher}, nil
}
