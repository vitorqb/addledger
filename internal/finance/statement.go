package finance

import "time"

// StatementEntry represents a single entry in a bank/credit card statement.
type StatementEntry struct {
	// Account is the account of the entry.
	Account string
	// Date is the date of the entry.
	Date time.Time
	// Description is a description of the entry.
	Description string
	// Amount is the amount of the entry.
	Ammount Ammount
}
