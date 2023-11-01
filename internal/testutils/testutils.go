package testutils

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/utils"
)

//
// Test Data

func Date1(t *testing.T) time.Time {
	out, err := time.Parse("2006-01-02", "1993-11-23")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func Date2(t *testing.T) time.Time {
	out, err := time.Parse("2006-01-02", "2001-01-01")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func FillPostingInput_1(t *testing.T, posting *input.PostingInput) {
	posting.SetAccount("ACC1")
	posting.SetAmmount(finance.Ammount{
		Commodity: "EUR",
		Quantity:  decimal.New(1220, -2),
	})
}

func FillPostingInput_2(t *testing.T, posting *input.PostingInput) {
	posting.SetAccount("ACC2")
	posting.SetAmmount(finance.Ammount{
		Commodity: "EUR",
		Quantity:  decimal.New(-1220, -2),
	})
}

func FillPostingInput_3(t *testing.T, posting *input.PostingInput) {
	posting.SetAccount("ACC3")
	posting.SetAmmount(finance.Ammount{
		Commodity: "EUR",
		Quantity:  decimal.New(999, -1),
	})
}

func JournalEntryInput_1(t *testing.T) *input.JournalEntryInput {
	journalEntryInput := input.NewJournalEntryInput()
	journalEntryInput.SetDate(Date1(t))
	journalEntryInput.SetDescription("Description1")
	posting_1 := journalEntryInput.AddPosting()
	FillPostingInput_1(t, posting_1)
	posting_2 := journalEntryInput.AddPosting()
	FillPostingInput_2(t, posting_2)
	return journalEntryInput
}

func Decimal_1(t *testing.T) decimal.Decimal {
	out, err := decimal.NewFromString("2.20")
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func Ammount_1(t *testing.T) *finance.Ammount {
	return &finance.Ammount{Commodity: "EUR", Quantity: Decimal_1(t)}
}

func Transaction_1(t *testing.T) *journal.Transaction {
	return &journal.Transaction{
		Description: "Description1",
		Date:        Date1(t),
		Posting: []journal.Posting{
			{
				Account: "ACC1",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(1220, -2),
				},
			},
			{
				Account: "ACC2",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(-1220, -2),
				},
			},
		},
	}
}

func Transaction_2(t *testing.T) *journal.Transaction {
	return &journal.Transaction{
		Description: "Description2",
		Date:        Date2(t),
		Posting: []journal.Posting{
			{
				Account: "ACC3",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(2000, -2),
				},
			},
			{
				Account: "ACC4",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(-2000, -2),
				},
			},
		},
	}
}

func Transaction_3(t *testing.T) *journal.Transaction {
	return &journal.Transaction{
		Description: "Description3",
		Date:        Date2(t),
		Comment:     "trip:brazil",
		Tags:        []journal.Tag{{Name: "trip", Value: "brazil"}},
		Posting: []journal.Posting{
			{
				Account: "ACC5",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(2001, -2),
				},
			},
			{
				Account: "ACC6",
				Ammount: finance.Ammount{
					Commodity: "EUR",
					Quantity:  decimal.New(-2001, -2),
				},
			},
		},
	}
}

//
// Helpers

func Setenv(t *testing.T, key, newValue string) (cleanup func()) {
	oldValue, existed := os.LookupEnv(key)
	err := os.Setenv(key, newValue)
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		var err error
		if existed {
			err = os.Setenv(key, oldValue)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Unsetenv unsets an environmental variable (if set), while returning
// a function to restore it's previous value (if any).
func Unsetenv(t *testing.T, key string) (cleanup func()) {
	oldValue, existed := os.LookupEnv(key)
	if !existed {
		return func() {}
	}
	err := os.Unsetenv(key)
	if err != nil {
		t.Fatal(err)
	}
	return func() {
		err := os.Setenv(key, oldValue)
		if err != nil {
			t.Fatal(err)
		}
	}
}

// Setenvs sets a batch of environmental variables and returns a
// cleanup function to restore their previous value.
// Usage: Setenvs(t, "VAR1", "VALUE1", "VAR2", "VALUE2")
func Setenvs(t *testing.T, keyValuePairs ...string) (cleanup func()) {
	if math.Mod(float64(len(keyValuePairs)), 2) != 0 {
		t.Fatal("Expected pair number of arguments")
	}
	it, err := utils.SplitArray[string](2, keyValuePairs)
	if err != nil {
		t.Fatal(err)
	}
	var cleanups []func()
	for {
		varSpec, err := it()
		if _, ok := err.(*utils.StopSplitArray); ok {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		varName := varSpec[0]
		varValue := varSpec[1]
		cleanups = append(cleanups, Setenv(t, varName, varValue))
	}
	return func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}

// Setenvs unsets a batch of environmental variables and returns a
// cleanup function to restore their previous value.
// Usage: Setenvs(t, "VAR1", "VALUE1", "VAR2", "VALUE2")
func Unsetenvs(t *testing.T, keys ...string) (cleanup func()) {
	var cleanups []func()
	for _, key := range keys {
		cleanups = append(cleanups, Unsetenv(t, key))
	}
	return func() {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}

// TestDataPath returns the absolute path to a testdata file.
func TestDataPath(t *testing.T, path string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(wd, "testdata", path)
}
