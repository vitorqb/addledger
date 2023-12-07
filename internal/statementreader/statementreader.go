package statementreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/vitorqb/addledger/internal/finance"
)

//go:generate $MOCKGEN --source=statementreader.go --destination=../../mocks/statementreader/statementreader_mock.go

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

// IStatementReader is an interface for reading a bank statement from a file and
// converting it to a list of statement entries.
type IStatementReader interface {
	Read(file io.Reader) ([]StatementEntry, error)
}

// CSVColumnMapping maps a csv column to a statement entry field.
type CSVColumnMapping struct {
	// Column is the column index.
	Column int
	// Importer is the field importer.
	Importer FieldImporter
}

// StatementReader implements IStatementReader.
type StatementReader struct {
	config Config
}

// Read implements StatementLoader.Read.
func (l *StatementReader) Read(reader io.Reader) ([]StatementEntry, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = l.config.Separator

	// Parse statement entries
	var statementEntries []StatementEntry
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading csv file: %w", err)
		}
		var statementEntry StatementEntry
		for _, columnMapping := range l.config.ColumnMappings {
			if columnMapping.Column >= len(record) {
				return nil, fmt.Errorf("column index out of range for field %T", columnMapping.Importer)
			}
			value := record[columnMapping.Column]
			if err := columnMapping.Importer.Import(&statementEntry, value); err != nil {
				return nil, fmt.Errorf("error importing field %T: %w", columnMapping.Importer, err)
			}
		}
		statementEntries = append(statementEntries, statementEntry)
	}

	// Set default values
	for i, statementEntry := range statementEntries {
		if statementEntry.Account == "" {
			statementEntry.Account = l.config.AccountName
		}
		if statementEntry.Ammount.Commodity == "" {
			statementEntry.Ammount.Commodity = l.config.DefaultCommodity
		}
		statementEntries[i] = statementEntry
	}

	return statementEntries, nil
}

// Config represents the options for a CSVReader.
type Config struct {
	// AccountName is the default account name for the statement entries.
	AccountName string
	// DefaultCommodity is the default commodity for the statement entries.
	DefaultCommodity string
	// Separator is the csv separator.
	Separator rune
	// ColumnMappings is the csv column mappings.
	ColumnMappings []CSVColumnMapping
}

var DefaultConfig = Config{
	AccountName:      "",
	DefaultCommodity: "EUR",
	Separator:        ',',
	ColumnMappings:   []CSVColumnMapping{},
}

// Option is a function that configures a CSVLoaderConfig.
type Option func(*Config)

func WithAccountName(accountName string) Option {
	return func(o *Config) {
		o.AccountName = accountName
	}
}

func WithDefaultCommodity(defaultCommodity string) Option {
	return func(o *Config) {
		o.DefaultCommodity = defaultCommodity
	}
}

func WithSeparator(separator rune) Option {
	return func(o *Config) {
		o.Separator = separator
	}
}

func WithLoaderMapping(columnMappings []CSVColumnMapping) Option {
	return func(o *Config) {
		o.ColumnMappings = columnMappings
	}
}

func NewStatementReader(options ...Option) *StatementReader {
	config := DefaultConfig
	for _, option := range options {
		option(&config)
	}
	return &StatementReader{config: config}
}
