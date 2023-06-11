package journal

import "github.com/shopspring/decimal"

// Ammount represents an Ammout in hledger. It contains a commodity
// and a quantity.
type Ammount struct {
	Commodity string
	Quantity  decimal.Decimal
}

func (a Ammount) Equal(a2 Ammount) bool {
	return a.Quantity.Equal(a2.Quantity) && a.Commodity == a2.Commodity
}
