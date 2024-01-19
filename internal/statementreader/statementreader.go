package statementreader

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"github.com/vitorqb/addledger/internal/finance"
)

//go:generate $MOCKGEN --source=statementreader.go --destination=../../mocks/statementreader/statementreader_mock.go

// IStatementReader is an interface for reading a bank statement from a file and
// converting it to a list of statement entries.
type IStatementReader interface {
	Read(file io.Reader, options ...Option) ([]finance.StatementEntry, error)
}

// CSVColumnMapping maps a csv column to a statement entry field.
type CSVColumnMapping struct {
	// Column is the column index.
	Column int
	// Importer is the field importer.
	Importer FieldImporter
}

type StatementReader struct{}

func (s *StatementReader) Read(reader io.Reader, options ...Option) ([]finance.StatementEntry, error) {
	config := parseOptions(options)
	csvReader := csv.NewReader(reader)
	csvReader.Comma = config.Separator

	// Parse statement entries
	var statementEntries []finance.StatementEntry
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading csv file: %w", err)
		}
		var statementEntry finance.StatementEntry
		for _, columnMapping := range config.ColumnMappings {
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
			statementEntry.Account = config.AccountName
		}
		if statementEntry.Ammount.Commodity == "" {
			statementEntry.Ammount.Commodity = config.DefaultCommodity
		}
		statementEntries[i] = statementEntry
	}

	// Sort
	if config.SortStrategy != nil {
		sort.Sort(config.SortStrategy.Clone(statementEntries))
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
	// Sort strategy to use (if any)
	SortStrategy SortStrategy
}

var DefaultConfig = Config{
	AccountName:      "",
	DefaultCommodity: "EUR",
	Separator:        ',',
	ColumnMappings:   []CSVColumnMapping{},
}

// A SortStrategy represents a strategy for sorting an array of StatementEntry.
// The `Clone` method must return a new SortStrategy of the same implementation
// but wrapping the given array.
type SortStrategy interface {
	sort.Interface
	Clone([]finance.StatementEntry) SortStrategy
}

type SortByDate []finance.StatementEntry

func (x SortByDate) Len() int                                    { return len(x) }
func (x SortByDate) Swap(i, j int)                               { x[i], x[j] = x[j], x[i] }
func (x SortByDate) Less(i, j int) bool                          { return x[i].Date.Before(x[j].Date) }
func (SortByDate) Clone(x []finance.StatementEntry) SortStrategy { return SortByDate(x) }

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

func WithSortStrategy(strategy SortStrategy) Option {
	return func(o *Config) {
		o.SortStrategy = strategy
	}
}

func NewStatementReader() *StatementReader { return &StatementReader{} }

func parseOptions(options []Option) Config {
	config := DefaultConfig
	for _, option := range options {
		option(&config)
	}
	return config
}
