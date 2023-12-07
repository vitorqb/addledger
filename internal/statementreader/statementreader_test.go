package statementreader_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/vitorqb/addledger/internal/finance"
	. "github.com/vitorqb/addledger/internal/statementreader"
)

func TestDateImporter(t *testing.T) {
	type testCase struct {
		dateStr       string
		format        string
		expectedDate  time.Time
		expectedError string
	}
	testCases := []testCase{
		{
			dateStr:       "2020-01-01",
			format:        "2006-01-02",
			expectedDate:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedError: "",
		},
		{
			dateStr:       "31/10/2023",
			format:        "02/01/2006",
			expectedDate:  time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
			expectedError: "",
		},
		{
			dateStr:      "10/31/2023",
			format:       "01/02/2006",
			expectedDate: time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			dateStr:       "10/31/2023",
			format:        "02/01/2006",
			expectedError: "invalid date (from format 02/01/2006): 10/31/2023",
		},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("%s-%s", tc.format, tc.dateStr)
		t.Run(testName, func(t *testing.T) {
			statementEntry := &StatementEntry{}
			err := DateImporter{tc.format}.Import(statementEntry, tc.dateStr)
			assert.Equal(t, tc.expectedDate, statementEntry.Date)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAccountImporter(t *testing.T) {
	type testCase struct {
		accountStr      string
		expectedAccount string
		expectedError   error
	}
	testCases := []testCase{
		{
			accountStr:      "Assets:Checking",
			expectedAccount: "Assets:Checking",
			expectedError:   nil,
		},
	}
	for _, tc := range testCases {
		statementEntry := &StatementEntry{}
		err := AccountImporter{}.Import(statementEntry, tc.accountStr)
		assert.Equal(t, tc.expectedAccount, statementEntry.Account)
		assert.ErrorIs(t, err, tc.expectedError)
	}
}

func TestAmmountImporter(t *testing.T) {
	type testCase struct {
		ammountStr      string
		expectedAmmount finance.Ammount
		expectedError   string
	}
	testCases := []testCase{
		{
			ammountStr:      "EUR 12.2",
			expectedAmmount: finance.Ammount{Commodity: "EUR", Quantity: decimal.New(122, -1)},
			expectedError:   "",
		},
		{
			ammountStr:      "12.2",
			expectedAmmount: finance.Ammount{Commodity: "", Quantity: decimal.New(122, -1)},
			expectedError:   "",
		},
		{
			ammountStr:    "FOO",
			expectedError: "invalid amount format: FOO",
		},
	}
	for _, tc := range testCases {
		statementEntry := &StatementEntry{}
		err := AmmountImporter{}.Import(statementEntry, tc.ammountStr)
		assert.Equal(t, tc.expectedAmmount, statementEntry.Ammount)
		if tc.expectedError != "" {
			assert.ErrorContains(t, err, tc.expectedError)
		} else {
			assert.NoError(t, err)
		}
	}
}
