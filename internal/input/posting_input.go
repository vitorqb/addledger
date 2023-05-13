package input

import "github.com/vitorqb/addledger/pkg/react"

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

func (i *PostingInput) GetAccount() (string, bool) {
	if rawValue, found := i.inputs["account"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (i *PostingInput) SetValue(value string) {
	i.inputs["value"] = value
	i.NotifyChange()
}

func (i *PostingInput) GetValue() (string, bool) {
	if rawValue, found := i.inputs["value"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
}
