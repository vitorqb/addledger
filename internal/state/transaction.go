package state

import (
	"time"

	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/pkg/react"
)

type TransactionData struct {
	react.React
	Date        MaybeValue[time.Time]
	Description MaybeValue[string]
	Tags        ArrayValue[journal.Tag]
}

func NewTransactionData() *TransactionData {
	out := &TransactionData{}
	out.Date = MaybeValue[time.Time]{}
	out.Date.AddOnChangeHook(out.NotifyChange)
	out.Description = MaybeValue[string]{}
	out.Description.AddOnChangeHook(out.NotifyChange)
	out.Tags = ArrayValue[journal.Tag]{}
	out.Tags.AddOnChangeHook(out.NotifyChange)
	return out
}
