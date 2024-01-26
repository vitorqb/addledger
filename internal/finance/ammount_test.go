package finance_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/finance"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

func TestNewBalance(t *testing.T) {

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
			result := NewBalance(tc.ammounts)
			assert.Equal(t, NewBalance(tc.expected), result)
		})
	}

}

func TestBalance(t *testing.T) {
	t.Run("SingleCommodity", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			balance := Balance{}
			assert.Equal(t, false, balance.SingleCommodity())
		})
		t.Run("Single", func(t *testing.T) {
			balance := NewBalance([]Ammount{*tu.Ammount_1(t)})
			assert.Equal(t, true, balance.SingleCommodity())
		})
		t.Run("Multiple ammounts single currency", func(t *testing.T) {
			balance := NewBalance([]Ammount{*tu.Ammount_1(t), *tu.Ammount_2(t)})
			assert.Equal(t, true, balance.SingleCommodity())
		})
		t.Run("Multiple currencies", func(t *testing.T) {
			ammounts := make([]Ammount, 2)
			ammounts[0] = *tu.Ammount_1(t)
			ammounts[1] = *tu.Ammount_1(t)
			ammounts[1].Commodity = "BRL"
			balance := NewBalance(ammounts)
			assert.Equal(t, false, balance.SingleCommodity())
		})
	})
	t.Run("IsZero", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			balance := NewBalance([]Ammount{})
			assert.Equal(t, true, balance.IsZero())
		})
		t.Run("Zero", func(t *testing.T) {
			ammounts := make([]Ammount, 2)
			ammounts[0] = *tu.Ammount_1(t)
			ammounts[1] = tu.Ammount_1(t).InvertSign()
			balance := NewBalance(ammounts)
			assert.Equal(t, true, balance.IsZero())
		})
		t.Run("Zero multiple currencies", func(t *testing.T) {
			ammounts := make([]Ammount, 4)
			ammounts[0] = *tu.Ammount_1(t)
			ammounts[1] = tu.Ammount_1(t).InvertSign()
			ammounts[2] = *tu.Ammount_1(t)
			ammounts[2].Commodity = "BRL"
			ammounts[3] = tu.Ammount_1(t).InvertSign()
			ammounts[3].Commodity = "BRL"
			balance := NewBalance(ammounts)
			assert.Equal(t, true, balance.IsZero())
		})
		t.Run("NonZero", func(t *testing.T) {
			ammounts := make([]Ammount, 2)
			ammounts[0] = *tu.Ammount_1(t)
			ammounts[1] = *tu.Ammount_2(t)
			balance := NewBalance(ammounts)
			assert.Equal(t, false, balance.IsZero())
		})
	})
}
