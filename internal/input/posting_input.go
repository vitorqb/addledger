package input

import (
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/pkg/react"
)

type (
	PostingInput struct {
		react.IReact
		inputs map[string]interface{}
	}
)

func NewPostingInput() *PostingInput {
	inputs := map[string]interface{}{}
	return &PostingInput{IReact: react.New(), inputs: inputs}
}

func (i *PostingInput) SetAccount(account string) {
	i.inputs["account"] = account
	i.NotifyChange()
}

func (i *PostingInput) ClearAccount() {
	delete(i.inputs, "account")
	i.NotifyChange()
}

func (i *PostingInput) GetAccount() (string, bool) {
	if rawValue, found := i.inputs["account"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (i *PostingInput) SetAmmount(value finance.Ammount) {
	i.inputs["ammount"] = value
	i.NotifyChange()
}

func (i *PostingInput) GetAmmount() (finance.Ammount, bool) {
	if rawValue, found := i.inputs["ammount"]; found {
		if value, ok := rawValue.(finance.Ammount); ok {
			return value, true
		}
	}
	return finance.Ammount{}, false
}
