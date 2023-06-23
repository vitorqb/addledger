package journal

import "time"

// Posting represents a Posting inside a transaction
type Posting struct {
	Account string
	Ammount []Ammount
}

// Transaction represents a transaction inside a journal.
type Transaction struct {
	Description string
	Date        time.Time
	Posting     []Posting
}

// An Account represents a hledger account
type Account string
