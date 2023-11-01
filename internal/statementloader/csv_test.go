package statementloader_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	. "github.com/vitorqb/addledger/internal/statementloader"
)

func TestCSVLoader(t *testing.T) {
	type testCase struct {
		name          string
		options       []CSVLoaderOption
		csvInput      string
		expected      []StatementEntry
		expectedError string
	}
	testCases := []testCase{
		{
			name: "Simple",
			options: []CSVLoaderOption{
				WithCSVLoaderAccountName("ACC"),
				WithCSVLoaderDefaultCommodity("EUR"),
				WithCSVLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: DateImporter},
					{Column: 1, Importer: DescriptionImporter},
					{Column: 2, Importer: AmmountImporter},
				}),
			},
			csvInput: `2023-10-31,FOO,12.21`,
			expected: []StatementEntry{
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
			name: "Two entries",
			options: []CSVLoaderOption{
				WithCSVLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: AccountImporter},
					{Column: 1, Importer: DateImporter},
					{Column: 2, Importer: DescriptionImporter},
					{Column: 3, Importer: AmmountImporter},
				}),
				WithCSVLoaderDefaultCommodity(""),
			},
			csvInput: "ACC,31/10/2023,FOO,12.21\nACC,30/10/2023,BAR,12.00",
			expected: []StatementEntry{
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
			options: []CSVLoaderOption{
				WithCSVLoaderMapping([]CSVColumnMapping{
					{Column: 10, Importer: DateImporter},
				}),
			},
			csvInput:      `10/31/2023`,
			expectedError: "column index out of range",
		},
		{
			name: "Invalid date",
			options: []CSVLoaderOption{
				WithCSVLoaderMapping([]CSVColumnMapping{
					{Column: 0, Importer: DateImporter},
				}),
			},
			csvInput:      `10/31/2023`,
			expectedError: "invalid date format",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			loader := NewCSVLoader(tc.options...)
			reader := strings.NewReader(tc.csvInput)
			entries, err := loader.Load(reader)
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, entries)
		})
	}
}
