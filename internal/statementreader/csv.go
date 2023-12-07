package statementreader

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

// CSVStatementReader loads a bank statement from a csv file.
type CSVStatementReader struct {
	config CSVReaderConfig
}

// Read implements StatementLoader.Read.
func (l *CSVStatementReader) Read(reader io.Reader) ([]StatementEntry, error) {
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

// CSVReaderConfig represents the options for a CSVReader.
type CSVReaderConfig struct {
	// AccountName is the default account name for the statement entries.
	AccountName string
	// DefaultCommodity is the default commodity for the statement entries.
	DefaultCommodity string
	// Separator is the csv separator.
	Separator rune
	// ColumnMappings is the csv column mappings.
	ColumnMappings []CSVColumnMapping
}

// DefaultCSVReaderConfig set the default CSVLoaderConfig.
var DefaultCSVReaderConfig = CSVReaderConfig{
	AccountName:      "",
	DefaultCommodity: "EUR",
	Separator:        ',',
	ColumnMappings:   []CSVColumnMapping{},
}

// CSVReaderOption is a function that configures a CSVLoaderConfig.
type CSVReaderOption func(*CSVReaderConfig)

// WithCSVReaderAccountName returns a CSVLoaderOption that sets the account name.
func WithCSVReaderAccountName(accountName string) CSVReaderOption {
	return func(o *CSVReaderConfig) {
		o.AccountName = accountName
	}
}

// WithCSVReaderDefaultCommodity returns a CSVLoaderOption that sets the default commodity.
func WithCSVReaderDefaultCommodity(defaultCommodity string) CSVReaderOption {
	return func(o *CSVReaderConfig) {
		o.DefaultCommodity = defaultCommodity
	}
}

// WithCSVSeparator returns a CSVLoaderOption that sets the separator.
func WithCSVSeparator(separator rune) CSVReaderOption {
	return func(o *CSVReaderConfig) {
		o.Separator = separator
	}
}

// WithCSVLoaderMapping returns a CSVLoaderOption that sets the column mappings.
func WithCSVLoaderMapping(columnMappings []CSVColumnMapping) CSVReaderOption {
	return func(o *CSVReaderConfig) {
		o.ColumnMappings = columnMappings
	}
}

// NewCSVLoader creates a new CSVStatementLoader.
func NewCSVLoader(options ...CSVReaderOption) *CSVStatementReader {
	config := DefaultCSVReaderConfig
	for _, option := range options {
		option(&config)
	}
	return &CSVStatementReader{config: config}
}
