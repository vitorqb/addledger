package journal

// Transaction represents a transaction inside a journal.
type Transaction struct {
	Description string `json:"tdescription"`
}

// An Account represents a hledger account
type Account string
