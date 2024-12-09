package statementloader

import (
	"fmt"
	"os"
	"strings"

	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
)

// Service can be used to load a statement into the app state.
type Service struct {
	state  *statemod.State
	reader statementreader.IStatementReader
}

// Load loads a statement into the app state.
func (c *Service) Load(config Config) error {
	if config.File == "" {
		return nil
	}
	options, err := ParseConfig(config)
	if err != nil {
		return fmt.Errorf("failed to load csv statement loader: %w", err)
	}
	csvFile, err := os.Open(config.File)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()
	statmntEntries, err := c.reader.Read(csvFile, options...)
	if err != nil {
		return fmt.Errorf("failed to load statement: %w", err)
	}
	c.state.SetStatementEntries(statmntEntries)
	return nil
}

// LoadFromFiles do the same as `Load` but reads the config from a json file.
func (c *Service) LoadFromFiles(statementFile, presetFile string) error {
	config, err := LoadConfig(statementFile, presetFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	return c.Load(config)
}

// New creates a new StatementLoaderSvc.
func New(state *statemod.State, reader statementreader.IStatementReader) *Service {
	return &Service{state: state, reader: reader}
}

// ParseConfig parses a statement loader config into statemtn reader options.
func ParseConfig(config Config) ([]statementreader.Option, error) {
	options := []statementreader.Option{}
	if acc := config.Account; acc != "" {
		options = append(options, statementreader.WithAccountName(acc))
	}
	if comm := config.Commodity; comm != "" {
		options = append(options, statementreader.WithDefaultCommodity(comm))
	}
	if sep := config.Separator; sep != "" {
		if len(sep) != 1 {
			return nil, fmt.Errorf("invalid csv separator: %s", sep)
		}
		options = append(options, statementreader.WithSeparator([]rune(sep)[0]))
	}
	if sortByStr := config.SortBy; sortByStr != "" {
		switch strings.ToLower(sortByStr) {
		case "date":
			options = append(options, statementreader.WithSortStrategy(statementreader.SortByDate{}))
		default:
			return nil, fmt.Errorf("invalid SortBy: %s", sortByStr)
		}
	}
	mapping := []statementreader.CSVColumnMapping{}
	if idate := config.DateFieldIndex; idate != -1 {
		importer := statementreader.DateImporter{Format: config.DateFormat}
		mapping = append(mapping, statementreader.CSVColumnMapping{Column: idate, Importer: importer})
	}
	if idescription := config.DescriptionFieldIndex; idescription != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: idescription, Importer: statementreader.DescriptionImporter{},
		})
	}
	if iaccount := config.AccountFieldIndex; iaccount != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: iaccount, Importer: statementreader.AccountImporter{},
		})
	}
	if iammount := config.AmmountFieldIndex; iammount != -1 {
		mapping = append(mapping, statementreader.CSVColumnMapping{
			Column: iammount, Importer: statementreader.AmmountImporter{},
		})
	}
	options = append(options, statementreader.WithLoaderMapping(mapping))
	return options, nil
}
