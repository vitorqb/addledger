package statementloader

import (
	"fmt"
	"time"

	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/input"
)

// StatementEntry represents a single entry in a bank statement.
type StatementEntry struct {
	// Account is the account of the entry.
	Account string
	// Date is the date of the entry.
	Date time.Time
	// Description is a description of the entry.
	Description string
	// Amount is the amount of the entry.
	Ammount finance.Ammount
}

// StatementLoader is an interface representing a bank statement loader.
type StatementLoader interface {
	// Load loads a bank statement from a file, and returns a list of
	// statement entries. Those entires contain infromation that will help
	// the user to create journal entries.
	Load(file string) ([]StatementEntry, error)
}

// A FieldImporter knows how to import a field from a string.
type FieldImporter func(statementEntry *StatementEntry, value string) error

// AccountImporter imports the account field.
func AccountImporter(statementEntry *StatementEntry, value string) error {
	statementEntry.Account = value
	return nil
}

var _ FieldImporter = AccountImporter

// DateImporter imports the date field.
func DateImporter(statementEntry *StatementEntry, value string) error {
	// Note: we are hardcoding the date formats here, which is not ideal.
	// We should probably allow the user to configure the date formats.

	// ISO format (yyyy-mm-dd)
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		statementEntry.Date = parsed
		return nil
	}

	// EU format (dd/mm/yyyy)
	if parsed, err := time.Parse("02/01/2006", value); err == nil {
		statementEntry.Date = parsed
		return nil
	}

	return fmt.Errorf("invalid date format: %s", value)
}

var _ FieldImporter = DateImporter

// DescriptionImporter imports the description field.
func DescriptionImporter(statementEntry *StatementEntry, value string) error {
	statementEntry.Description = value
	return nil
}

var _ FieldImporter = DescriptionImporter

// AmmountImporter imports the amount field.
func AmmountImporter(statementEntry *StatementEntry, value string) error {
	if parsed, err := input.TextToAmmount(value); err == nil {
		statementEntry.Ammount = parsed
		return nil
	}
	return fmt.Errorf("invalid amount format: %s", value)
}

var _ FieldImporter = AmmountImporter
