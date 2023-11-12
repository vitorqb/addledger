// The app package contains the logic to manage and run the entire application. It
// coordinates most the functionality of the app and the state management.
package app

import (
	"fmt"
	"os"

	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
)

// LoadStatement loads a statement from a file and saves it to the state.
func LoadStatement(
	loader statementloader.StatementLoader,
	file string,
	state *state.State,
) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	entries, err := loader.Load(f)
	if err != nil {
		return fmt.Errorf("failed to load statement: %w", err)
	}

	state.SetStatementEntries(entries)
	return nil
}
