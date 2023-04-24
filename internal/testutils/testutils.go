package testutils

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vitorqb/addledger/internal/input"
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

func JournalEntryInput1(t *testing.T) *input.JournalEntryInput {
	journalEntryInput := input.NewJournalEntryInput()
	journalEntryInput.SetDate(Date1(t))
	journalEntryInput.SetDescription("Description1")
	posting1 := journalEntryInput.AddPosting()
	posting1.SetAccount("ACC1")
	posting1.SetValue("EUR 12.20")
	posting2 := journalEntryInput.AddPosting()
	posting2.SetAccount("ACC2")
	posting2.SetValue("EUR -12.20")
	return journalEntryInput
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

// TestDataPath returns the absolute path to a testdata file.
func TestDataPath(t *testing.T, path string) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(wd, "testdata", path)
}
