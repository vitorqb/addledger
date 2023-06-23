package hledger

import (
	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/journal"
)

func ParsePostingsJson(jsonpostings []JSONPosting) ([]journal.Posting, error) {
	postings := []journal.Posting{}
	for _, jsonposting := range jsonpostings {
		ammounts := []journal.Ammount{}
		for _, jsonammount := range jsonposting.Ammount {
			quantity := decimal.New(jsonammount.Quantity.DecimalMantissa, -1*jsonammount.Quantity.DecimalPlaces)
			ammount := journal.Ammount{
				Commodity: jsonammount.Commodity,
				Quantity:  quantity,
			}
			ammounts = append(ammounts, ammount)
		}
		posting := journal.Posting{
			Account: jsonposting.Account,
			Ammount: ammounts,
		}
		postings = append(postings, posting)
	}
	return postings, nil
}
