package parsing_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	. "github.com/vitorqb/addledger/internal/parsing"
)

func TestTextToAmmount(t *testing.T) {
	type testcase struct {
		text     string
		ammount  finance.Ammount
		errorMsg string
	}
	var testcases = []testcase{
		{
			text: "EUR 12.20",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text: "EUR 99999.99999",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.NewFromFloat(99999.99999),
			},
		},
		{
			text: "12.20",
			ammount: finance.Ammount{
				Commodity: "",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text:     "12,20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12 12",
			errorMsg: "invalid format",
		},
		{
			text:     "12 FOO",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR  12.20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12.20 ",
			errorMsg: "invalid format",
		},
		{
			text:     " EUR 12.20 ",
			errorMsg: "invalid format",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.text, func(t *testing.T) {
			result, err := TextToAmmount(tc.text)
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.ammount, result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}
