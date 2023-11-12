package statementloader

import (
	"fmt"
	"io"
	"time"

	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/input"
)

//go:generate $MOCKGEN --source=statementloader.go --destination=../../mocks/statementloader/statementloader_mock.go

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
	Load(file io.Reader) ([]StatementEntry, error)
}

// A FieldImporter knows how to import a field from a string.
type FieldImporter interface {
	// Import imports the field from the string.
	Import(statementEntry *StatementEntry, value string) error
}

// AccountImporter imports the account field.
type AccountImporter struct{}

func (a AccountImporter) Import(statementEntry *StatementEntry, value string) error {
	statementEntry.Account = value
	return nil
}

var _ FieldImporter = AccountImporter{}

// DateImporter imports the date field.
type DateImporter struct {
	Format string
}

func (d DateImporter) Import(statementEntry *StatementEntry, value string) error {
	// Note: we are hardcoding the date formats here, which is not ideal.
	// We should probably allow the user to configure the date formats.
	if d.Format != "" {
		if parsed, err := time.Parse(d.Format, value); err == nil {
			statementEntry.Date = parsed
			return nil
		}
	}
	return fmt.Errorf("invalid date (from format %s): %s", d.Format, value)
}

var _ FieldImporter = DateImporter{}

// DescriptionImporter imports the description field.
type DescriptionImporter struct{}

func (d DescriptionImporter) Import(statementEntry *StatementEntry, value string) error {
	statementEntry.Description = value
	return nil
}

var _ FieldImporter = DescriptionImporter{}

// AmmountImporter imports the amount field.
type AmmountImporter struct{}

func (a AmmountImporter) Import(statementEntry *StatementEntry, value string) error {
	if parsed, err := input.TextToAmmount(value); err == nil {
		statementEntry.Ammount = parsed
		return nil
	}
	return fmt.Errorf("invalid amount format: %s", value)
}

var _ FieldImporter = AmmountImporter{}
