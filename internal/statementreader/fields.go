package statementreader

import (
	"fmt"
	"time"

	"github.com/vitorqb/addledger/internal/parsing"
)

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
	if parsed, err := parsing.TextToAmmount(value); err == nil {
		statementEntry.Ammount = parsed
		return nil
	}
	return fmt.Errorf("invalid amount format: %s", value)
}

var _ FieldImporter = AmmountImporter{}
