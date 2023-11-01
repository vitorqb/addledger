package hledger

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
)

func ParsePostingsJson(jsonpostings []JSONPosting) ([]journal.Posting, error) {
	postings := []journal.Posting{}
	for _, jsonposting := range jsonpostings {
		ammounts := []finance.Ammount{}
		for _, jsonammount := range jsonposting.Ammount {
			quantity := decimal.New(jsonammount.Quantity.DecimalMantissa, -1*jsonammount.Quantity.DecimalPlaces)
			ammount := finance.Ammount{
				Commodity: jsonammount.Commodity,
				Quantity:  quantity,
			}
			ammounts = append(ammounts, ammount)
		}
		// IMPORTANT NOTE: In the JSON each posting has an array of Ammounts.
		// I don't why that's the case. In my personal journal all postings
		// have exactly 1 ammount.
		// For now, WE DO NOT SUPPORT postings with more than 1 ammount
		if len(ammounts) > 1 {
			logrus.
				WithField("posting", jsonposting).
				Warn("Found posting with more than 1 ammount. This is not supported.")
			continue
		}
		if len(ammounts) < 1 {
			logrus.
				WithField("posting", jsonposting).
				Warn("Found posting with no ammount. This is not supported.")
			continue
		}
		posting := journal.Posting{
			Account: jsonposting.Account,
			Ammount: ammounts[0],
		}
		postings = append(postings, posting)
	}
	return postings, nil
}
