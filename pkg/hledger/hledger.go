package hledger

import (
	"fmt"
	"os/exec"
	"strings"
)

//go:generate mockgen --source=hledger.go --destination=../../mocks/hledger/hledger_mock.go

// IClient represents an HLedger client.
type IClient interface {
	// Accounts returns a list of all known accounts.
	Accounts() ([]string, error)
}

// Client is the default implementation for IClient.
type Client struct {
	executable string
	ledgerFile string
}

var _ IClient = &Client{}

func (c *Client) Accounts() (accounts []string, err error) {
	cmdArgs := []string{"accounts"}
	if c.ledgerFile != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--file=%s", c.ledgerFile))
	}
	cmd := exec.Command(c.executable, cmdArgs...)
	cmdOutputBytes, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("Failed to get accounts: %w", err)
	}
	cmdOutputStr := strings.TrimSpace(string(cmdOutputBytes))
	return strings.Split(cmdOutputStr, "\n"), nil
}

func NewClient(executable, ledgerFile string) *Client {
	return &Client{
		executable: executable,
		ledgerFile: ledgerFile,
	}
}
