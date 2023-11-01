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

// Returns the balance for each currency in a list of Ammounts.
func Balance(ammounts []Ammount) []Ammount {
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
	return result
}
