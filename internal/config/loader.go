package config

import (
	"fmt"
	"os/exec"
	"strings"
)

// ILoader load config variables dinamically.
type ILoader interface {
	// JournalFile is a function that finds which ledger file to use
	// if the user did not specify one. It loads the path used by the ledger
	// executable.
	JournalFile(ledgerExecutable string) (string, error)
}

// Loader is the default implementation of ILoader.
type Loader struct{}

var _ ILoader = (*Loader)(nil)

// JournalFile implements ILoader.JournalFile.
func (l *Loader) JournalFile(ledgerExecutable string) (string, error) {
	output, err := exec.Command(ledgerExecutable, "files").Output()
	if err != nil {
		return "", fmt.Errorf("error finding journal file from executable %s: %w", ledgerExecutable, err)
	}
	outputStr := string(output)
	firstLine := strings.Split(outputStr, "\n")[0]
	return firstLine, nil
}

// NewLoader returns a new instance of Loader.
func NewLoader() *Loader {
	return &Loader{}
}
