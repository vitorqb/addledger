package parsing

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
)

func TextToAmmount(x string) (finance.Ammount, error) {
	var err error
	var quantity decimal.Decimal
	var commodity string
	switch words := strings.Split(x, " "); len(words) {
	case 1:
		quantity, err = decimal.NewFromString(words[0])
	case 2:
		commodity = words[0]
		quantity, err = decimal.NewFromString(words[1])
	default:
		return finance.Ammount{}, fmt.Errorf("invalid format")
	}
	if err != nil {
		return finance.Ammount{}, fmt.Errorf("invalid format: %w", err)
	}
	return finance.Ammount{Commodity: commodity, Quantity: quantity}, nil
}
