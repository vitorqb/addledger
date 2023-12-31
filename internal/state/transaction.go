package state

import (
	"time"

	"github.com/vitorqb/addledger/pkg/react"
)

type TransactionData struct {
	react.React
	Date MaybeValue[time.Time]
}

func NewTransactionData() *TransactionData {
	out := &TransactionData{}
	out.Date = MaybeValue[time.Time]{}
	out.Date.AddOnChangeHook(out.NotifyChange)
	return out
}
