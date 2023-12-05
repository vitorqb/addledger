package app

import (
	"fmt"
	"os"

	"github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/controller"
	"github.com/vitorqb/addledger/internal/injector"
	statemod "github.com/vitorqb/addledger/internal/state"
)

// CSVStatementLoader can be used to load a CSV statement into the app state.
type CSVStatementLoader struct {
	// The app state.
	state *statemod.State
}

// Load loads a CSV statement into the app state.
func (c *CSVStatementLoader) Load(config config.CSVStatementLoaderConfig) error {
	if config.File == "" {
		return nil
	}
	loader, err := injector.CSVStatementLoader(config)
	if err != nil {
		return fmt.Errorf("failed to load csv statement loader: %w", err)
	}
	csvFile, err := os.Open(config.File)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()
	statmntEntries, err := loader.Load(csvFile)
	if err != nil {
		return fmt.Errorf("failed to load statement: %w", err)
	}
	c.state.SetStatementEntries(statmntEntries)
	return nil
}

// NewCSVStatementLoader creates a new CSVStatementLoader.
func NewCSVStatementLoader(state *statemod.State) *CSVStatementLoader {
	return &CSVStatementLoader{state: state}
}

var _ controller.ICSVStatementLoader = (*CSVStatementLoader)(nil)
