package hledger

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/vitorqb/addledger/internal/journal"
)

//go:generate $MOCKGEN --source=hledger.go --destination=../../mocks/hledger/hledger_mock.go

// IClient represents an HLedger client.
type IClient interface {
	// Accounts returns a list of all known accounts.
	Accounts() ([]string, error)
	// Transactions returns a list of all known transactions.
	Transactions() ([]journal.Transaction, error)
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

func (c *Client) Transactions() ([]journal.Transaction, error) {
	var transactions []journal.Transaction
	cmdArgs := []string{}
	if c.ledgerFile != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--file=%s", c.ledgerFile))
	}
	cmdArgs = append(cmdArgs, "print", "--output-format=json")
	cmd := exec.Command(c.executable, cmdArgs...)
	cmdOutputBytes, err := cmd.Output()
	if err != nil {
		return transactions, fmt.Errorf("failed to get transactions: %w", err)
	}
	err = json.Unmarshal(cmdOutputBytes, &transactions)
	if err != nil {
		return transactions, fmt.Errorf("failed to unmarshall transactions: %w", err)
	}
	return transactions, nil
}

func NewClient(executable, ledgerFile string) *Client {
	return &Client{
		executable: executable,
		ledgerFile: ledgerFile,
	}
}
