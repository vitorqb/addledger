package finance

import (
	"github.com/shopspring/decimal"
)

// Ammount represents an Ammout in hledger. It contains a commodity
// and a quantity.
type Ammount struct {
	Commodity string
	Quantity  decimal.Decimal
}

func (a Ammount) Equal(a2 Ammount) bool {
	return a.Quantity.Equal(a2.Quantity) && a.Commodity == a2.Commodity
}

func (a Ammount) InvertSign() Ammount {
	return Ammount{a.Commodity, a.Quantity.Neg()}
}

func (a Ammount) Div(d decimal.Decimal) Ammount {
	return Ammount{a.Commodity, a.Quantity.Div(d)}
}

func (a Ammount) Mul(d decimal.Decimal) Ammount {
	return Ammount{a.Commodity, a.Quantity.Mul(d)}
}

func (a Ammount) Round(i int32) Ammount {
	return Ammount{a.Commodity, a.Quantity.Round(i)}
}

// A balance is a list of Ammounts, where each Ammount has a different
// commodity. It represents the balance of a transaction.
type Balance struct {
	ammounts []Ammount // Will always have ONLY 1 ammount per commodity
}

func (b Balance) Ammounts() []Ammount { return b.ammounts }

func (b Balance) SingleCommodity() bool { return len(b.ammounts) == 1 }

func (b Balance) IsZero() bool {
	for _, ammount := range b.ammounts {
		if !ammount.Quantity.IsZero() {
			return false
		}
	}
	return true
}

// Returns the balance for each currency in a list of Ammounts.
func NewBalance(ammounts []Ammount) Balance {
	commoditiesQuantityMap := map[string]decimal.Decimal{}
	for _, ammount := range ammounts {
		commoditiesQuantityMap[ammount.Commodity] = decimal.Zero
	}
	for _, ammount := range ammounts {
		commoditiesQuantityMap[ammount.Commodity] = commoditiesQuantityMap[ammount.Commodity].Add(ammount.Quantity)
	}
	result := []Ammount{}
	for commodity, quantity := range commoditiesQuantityMap {
		if !quantity.Equal(decimal.Zero) {
			result = append(result, Ammount{commodity, quantity})
		}
	}
	return Balance{result}
}
