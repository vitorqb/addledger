package input

import (
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
