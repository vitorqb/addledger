// uerinput is responsible for managing data inputted by the user.
package userinput

import (
	"fmt"

	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"
)

type ErrMissingAmmount struct{}

func (e ErrMissingAmmount) Error() string {
	return "one of the postings is missing the ammount"
}

type ErrMissingAccount struct{}

func (e ErrMissingAccount) Error() string {
	return "one of the postings is missing the account"
}

type ErrUnbalancedPosting struct{}

func (e ErrUnbalancedPosting) Error() string {
	return "postings are not balanced"
}

func TransactionRepr(t *state.TransactionData) string {
	var out string
	if date, found := t.Date.Get(); found {
		out += date.Format("2006-01-02")
	}
	if description, found := t.Description.Get(); found {
		out += " " + description
	}
	for i, tag := range t.Tags.Get() {
		if i == 0 {
			out += " ;"
		}
		out += " " + tag.Name + ":" + tag.Value
	}
	for _, posting := range t.Postings.Get() {
		out += "\n" + "    " + PostingRepr(posting)
	}
	return out
}

func PostingRepr(p *state.PostingData) string {
	out := ""
	if account, found := p.Account.Get(); found {
		out += string(account)
	}
	out += "    "
	if ammount, found := p.Ammount.Get(); found {
		if ammount.Commodity != "" {
			out += ammount.Commodity + " "
		}
		out += ammount.Quantity.String()
	}
	return out
}

func TransactionFromData(t *state.TransactionData) (journal.Transaction, error) {
	ammounts := make([]finance.Ammount, 0, len(t.Postings.Get()))
	postings := make([]journal.Posting, 0, len(t.Postings.Get()))
	for _, postingData := range t.Postings.Get() {
		ammount, found := postingData.Ammount.Get()
		if !found {
			return journal.Transaction{}, ErrMissingAmmount{}
		}
		account, found := postingData.Account.Get()
		if !found {
			return journal.Transaction{}, ErrMissingAccount{}
		}
		ammounts = append(ammounts, ammount)
		posting := journal.Posting{Account: string(account), Ammount: ammount}
		postings = append(postings, posting)
	}

	// If we have a single currency, we can check if the postings are balanced.
	if balance := finance.Balance(ammounts); len(balance) == 1 && !balance[0].Quantity.IsZero() {
		return journal.Transaction{}, ErrUnbalancedPosting{}
	}

	description, found := t.Description.Get()
	if !found {
		return journal.Transaction{}, fmt.Errorf("missing description")
	}

	date, found := t.Date.Get()
	if !found {
		return journal.Transaction{}, fmt.Errorf("missing date")
	}

	comment := ""
	for i, tag := range t.Tags.Get() {
		if i != 0 {
			comment += " "
		}
		comment += TagToText(tag)
	}

	return journal.Transaction{
		Description: description,
		Date:        date,
		Comment:     comment,
		Posting:     postings,
	}, nil
}

func TagToText(tag journal.Tag) string {
	return tag.Name + ":" + tag.Value
}

// ExtractPostings returns a list of journal.Posting from a list of
// state.PostingData.
func ExtractPostings(postings []*state.PostingData) []journal.Posting {
	var out []journal.Posting
	for _, posting := range postings {
		account, found := posting.Account.Get()
		if !found {
			continue
		}
		ammount, found := posting.Ammount.Get()
		if !found {
			continue
		}
		posting := journal.Posting{Account: string(account), Ammount: ammount}
		out = append(out, posting)
	}
	return out
}
