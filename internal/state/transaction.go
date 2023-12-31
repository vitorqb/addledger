package state

import (
	"time"

	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/pkg/react"
)

// PostingData is a struct that holds the data of a posting inputted by the user.
type PostingData struct {
	react.React
	Account MaybeValue[journal.Account]
	Ammount MaybeValue[finance.Ammount]
}

func NewPostingData() *PostingData {
	out := &PostingData{}
	out.Account.AddOnChangeHook(out.NotifyChange)
	out.Ammount.AddOnChangeHook(out.NotifyChange)
	return out
}

// TransactionData is a struct that holds the data of a transaction inputted by the user.
type TransactionData struct {
	react.React
	Date        MaybeValue[time.Time]
	Description MaybeValue[string]
	Tags        ArrayValue[journal.Tag]
	Postings    ArrayValue[*PostingData]
}

func NewTransactionData() *TransactionData {
	out := &TransactionData{}
	out.Date.AddOnChangeHook(out.NotifyChange)
	out.Description.AddOnChangeHook(out.NotifyChange)
	out.Tags.AddOnChangeHook(out.NotifyChange)
	out.Postings.AddOnChangeHook(out.NotifyChange)
	return out
}
