package journal

import (
	"time"

	"github.com/vitorqb/addledger/internal/finance"
)

// Posting represents a Posting inside a transaction
type Posting struct {
	Account string
	Ammount finance.Ammount
}

// Transaction represents a transaction inside a journal.
type Transaction struct {
	Description string
	Date        time.Time
	Posting     []Posting
	Comment     string
	Tags        []Tag
}

// An Account represents a hledger account
type Account string

// A Tag represents a hledger tag
type Tag struct {
	Name  string
	Value string
}

// PostingsBalance returns the balance of the postings.
func PostingsBalance(postings []Posting) finance.Balance {
	ammounts := make([]finance.Ammount, len(postings))
	for i, posting := range postings {
		ammounts[i] = posting.Ammount
	}
	return finance.NewBalance(ammounts)
}
