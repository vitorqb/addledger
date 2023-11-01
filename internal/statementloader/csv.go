package statementloader

import (
	"encoding/csv"
	"fmt"
	"io"
)

// CSVColumnMapping maps a csv column to a statement entry field.
type CSVColumnMapping struct {
	// Column is the column index.
	Column int
	// Importer is the field importer.
	Importer FieldImporter
}

// CSVLoader loads a bank statement from a csv file.
type CSVLoader struct {
	config CSVLoaderConfig
}

// Load implements StatementLoader.Load.
func (l *CSVLoader) Load(reader io.Reader) ([]StatementEntry, error) {
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
			if err := columnMapping.Importer(&statementEntry, value); err != nil {
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

// CSVLoaderConfig represents the options for a CSVLoader.
type CSVLoaderConfig struct {
	// AccountName is the default account name for the statement entries.
	AccountName string
	// DefaultCommodity is the default commodity for the statement entries.
	DefaultCommodity string
	// Separator is the csv separator.
	Separator rune
	// ColumnMappings is the csv column mappings.
	ColumnMappings []CSVColumnMapping
}

// DefaultCSVLoaderConfig set the default CSVLoaderConfig.
var DefaultCSVLoaderConfig = CSVLoaderConfig{
	AccountName:      "",
	DefaultCommodity: "EUR",
	Separator:        ',',
	ColumnMappings:   []CSVColumnMapping{},
}

// CSVLoaderOption is a function that configures a CSVLoaderConfig.
type CSVLoaderOption func(*CSVLoaderConfig)

// WithCSVLoaderAccountName returns a CSVLoaderOption that sets the account name.
func WithCSVLoaderAccountName(accountName string) CSVLoaderOption {
	return func(o *CSVLoaderConfig) {
		o.AccountName = accountName
	}
}

// WithCSVLoaderDefaultCommodity returns a CSVLoaderOption that sets the default commodity.
func WithCSVLoaderDefaultCommodity(defaultCommodity string) CSVLoaderOption {
	return func(o *CSVLoaderConfig) {
		o.DefaultCommodity = defaultCommodity
	}
}

// WithCSVSeparator returns a CSVLoaderOption that sets the separator.
func WithCSVSeparator(separator rune) CSVLoaderOption {
	return func(o *CSVLoaderConfig) {
		o.Separator = separator
	}
}

// WithCSVLoaderMapping returns a CSVLoaderOption that sets the column mappings.
func WithCSVLoaderMapping(columnMappings []CSVColumnMapping) CSVLoaderOption {
	return func(o *CSVLoaderConfig) {
		o.ColumnMappings = columnMappings
	}
}

// NewCSVLoader creates a new CSVStatementLoader.
func NewCSVLoader(options ...CSVLoaderOption) *CSVLoader {
	config := DefaultCSVLoaderConfig
	for _, option := range options {
		option(&config)
	}
	return &CSVLoader{config: config}
}
