package finance_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/finance"
)

func TestBalance(t *testing.T) {

	type test struct {
		name     string
		ammounts []Ammount
		expected []Ammount
	}

	tests := []test{
		{
			name:     "Empty",
			ammounts: []Ammount{},
			expected: []Ammount{},
		},
		{
			name: "One long",
			ammounts: []Ammount{
				{
					Commodity: "BRL",
					Quantity:  decimal.New(9999, -3),
				},
			},
			expected: []Ammount{
				{
					Commodity: "BRL",
					Quantity:  decimal.New(9999, -3),
				},
			},
		},
		{
			name: "Balanced",
			ammounts: []Ammount{
				{
					Commodity: "BRL",
					Quantity:  decimal.New(9999, -3),
				},
				{
					Commodity: "BRL",
					Quantity:  decimal.New(-9999, -3),
				},
			},
			expected: []Ammount{},
		},
		{
			name: "Unbalanced",
			ammounts: []Ammount{
				{
					Commodity: "BRL",
					Quantity:  decimal.New(9999, -3),
				},
				{
					Commodity: "BRL",
					Quantity:  decimal.New(-9999, -3),
				},
				{
					Commodity: "EUR",
					Quantity:  decimal.New(2, -3),
				},
			},
			expected: []Ammount{
				{
					Commodity: "EUR",
					Quantity:  decimal.New(2, -3),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Balance(tc.ammounts)
			assert.Equal(t, tc.expected, result)
		})
	}

}
