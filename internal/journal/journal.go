package journal

import (
	"time"
)

// Posting represents a Posting inside a transaction
type Posting struct {
	Account string
	Ammount Ammount
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
