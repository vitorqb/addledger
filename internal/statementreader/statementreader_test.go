package statementreader_test

import (
	"fmt"
	"strings"
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
			statementEntry := &finance.StatementEntry{}
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
		statementEntry := &finance.StatementEntry{}
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
		statementEntry := &finance.StatementEntry{}
		err := AmmountImporter{}.Import(statementEntry, tc.ammountStr)
		assert.Equal(t, tc.expectedAmmount, statementEntry.Ammount)
		if tc.expectedError != "" {
			assert.ErrorContains(t, err, tc.expectedError)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestCSVLoader(t *testing.T) {
	type testCase struct {
		name          string
		options       []Option
		csvInput      string
		expected      []finance.StatementEntry
		expectFn      func([]finance.StatementEntry)
		expectedError string
	}
	testCases := []testCase{
		{
			name: "Simple",
			options: []Option{
				WithAccountName("ACC"),
				WithDefaultCommodity("EUR"),
				WithLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: DateImporter{"2006-01-02"}},
					{Column: 1, Importer: DescriptionImporter{}},
					{Column: 2, Importer: AmmountImporter{}},
				}),
			},
			csvInput: `2023-10-31,FOO,12.21`,
			expected: []finance.StatementEntry{
				{
					Account:     "ACC",
					Date:        time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
					Description: "FOO",
					Ammount: finance.Ammount{
						Commodity: "EUR",
						Quantity:  decimal.New(1221, -2),
					},
				},
			},
			expectedError: "",
		},
		{
			name: "Sort by date",
			options: []Option{
				WithLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: DateImporter{"2006-01-02"}},
				}),
				WithSortStrategy(SortByDate{}),
			},
			csvInput: "2023-10-30\n2023-10-29",
			expectFn: func(x []finance.StatementEntry) {
				assert.Equal(t, time.Date(2023, 10, 29, 0, 0, 0, 0, time.UTC), x[0].Date)
				assert.Equal(t, time.Date(2023, 10, 30, 0, 0, 0, 0, time.UTC), x[1].Date)
			},
			expectedError: "",
		},
		{
			name: "Two entries",
			options: []Option{
				WithLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: AccountImporter{}},
					{Column: 1, Importer: DateImporter{"02/01/2006"}},
					{Column: 2, Importer: DescriptionImporter{}},
					{Column: 3, Importer: AmmountImporter{}},
				}),
				WithDefaultCommodity(""),
			},
			csvInput: "ACC,31/10/2023,FOO,12.21\nACC,30/10/2023,BAR,12.00",
			expected: []finance.StatementEntry{
				{
					Account:     "ACC",
					Date:        time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
					Description: "FOO",
					Ammount: finance.Ammount{
						Commodity: "",
						Quantity:  decimal.New(1221, -2),
					},
				},
				{
					Account:     "ACC",
					Date:        time.Date(2023, 10, 30, 0, 0, 0, 0, time.UTC),
					Description: "BAR",
					Ammount: finance.Ammount{
						Commodity: "",
						Quantity:  decimal.New(1200, -2),
					},
				},
			},
		},
		{
			name: "Column out of range",
			options: []Option{
				WithLoaderMapping([]CSVColumnMapping{
					{Column: 10, Importer: DateImporter{}},
				}),
			},
			csvInput:      `10/31/2023`,
			expectedError: "column index out of range",
		},
		{
			name: "Invalid date",
			options: []Option{
				WithLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: DateImporter{"2006-01-02"}},
				}),
			},
			csvInput:      `10/31/2023`,
			expectedError: "invalid date (from format 2006-01-02): 10/31/2023",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loader := NewStatementReader()
			reader := strings.NewReader(tc.csvInput)
			entries, err := loader.Read(reader, tc.options...)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			if tc.expectFn != nil {
				if len(tc.expected) > 0 {
					panic("only one of expected or expectFn should be set")
				}
				tc.expectFn(entries)
			} else {
				assert.Equal(t, tc.expected, entries)
			}
		})
	}
}
