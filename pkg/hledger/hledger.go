package hledger

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/journal"
)

//go:generate $MOCKGEN --source=hledger.go --destination=../../mocks/hledger/hledger_mock.go

// IClient represents an HLedger client.
type IClient interface {
	// Accounts returns a list of all known accounts.
	Accounts() ([]journal.Account, error)
	// Transactions returns a list of all known transactions.
	Transactions() ([]journal.Transaction, error)
}

var _ IClient = &Client{}

func (c *Client) Accounts() (accounts []journal.Account, err error) {
	cmdArgs := []string{"accounts"}
	if c.ledgerFile != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--file=%s", c.ledgerFile))
	}
	cmd := exec.Command(c.executable, cmdArgs...)
	cmdOutputBytes, err := cmd.Output()
	if err != nil {
		return []journal.Account{}, fmt.Errorf("Failed to get accounts: %w", err)
	}
	cmdOutputStr := strings.TrimSpace(string(cmdOutputBytes))
	for _, acc := range strings.Split(cmdOutputStr, "\n") {
		accounts = append(accounts, journal.Account(acc))
	}
	return accounts, nil
}

func (c *Client) Transactions() ([]journal.Transaction, error) {
	transactions := []journal.Transaction{}
	jsontransactions := []JSONTransaction{}
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
	err = json.Unmarshal(cmdOutputBytes, &jsontransactions)
	if err != nil {
		return transactions, fmt.Errorf("failed to unmarshall transactions: %w", err)
	}
	for _, jsontransaction := range jsontransactions {
		date, err := time.Parse("2006-01-02", jsontransaction.Date)
		if err != nil {
			logrus.Warn("Failed to parse date: %w", err)
			continue
		}
		postings, err := ParsePostingsJson(jsontransaction.Postings)
		if err != nil {
			logrus.Warn("Failed to parse postings: %w", err)
		}
		transaction := journal.Transaction{
			Description: jsontransaction.Description,
			Date:        date,
			Posting:     postings,
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func NewClient(executable, ledgerFile string) *Client {
	return &Client{
		executable: executable,
		ledgerFile: ledgerFile,
	}
}
