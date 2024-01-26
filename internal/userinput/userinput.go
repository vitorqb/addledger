// uerinput is responsible for managing data inputted by the user.
package userinput

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/state"

	"regexp"
	"strings"
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

func PostingFromData(p *state.PostingData) (journal.Posting, error) {
	ammount, found := p.Ammount.Get()
	if !found {
		return journal.Posting{}, ErrMissingAmmount{}
	}
	account, found := p.Account.Get()
	if !found {
		return journal.Posting{}, ErrMissingAccount{}
	}
	return journal.Posting{Account: string(account), Ammount: ammount}, nil
}

func PostingsFromData(postings []*state.PostingData) ([]journal.Posting, error) {
	out := make([]journal.Posting, 0, len(postings))
	for _, posting := range postings {
		p, err := PostingFromData(posting)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func TransactionFromData(t *state.TransactionData) (journal.Transaction, error) {
	ammounts := make([]finance.Ammount, 0, len(t.Postings.Get()))
	postings, err := PostingsFromData(t.Postings.Get())
	if err != nil {
		return journal.Transaction{}, err
	}
	for _, posting := range postings {
		ammounts = append(ammounts, posting.Ammount)
	}

	// If we have a single currency, we can check if the postings are balanced.
	if balance := finance.NewBalance(ammounts); balance.SingleCommodity() && !balance.IsZero() {
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

// CountCommodities returns the number of different commodities in a list of
// postings.
func CountCommodities(postings []state.PostingData) int {
	commodities := make(map[string]bool)
	for _, posting := range postings {
		ammount, found := posting.Ammount.Get()
		if !found {
			continue
		}
		commodities[ammount.Commodity] = true
	}
	return len(commodities)
}

func PostingBalance(postings []*state.PostingData) finance.Balance {
	var ammounts []finance.Ammount
	for _, posting := range postings {
		ammount, found := posting.Ammount.Get()
		if found {
			ammounts = append(ammounts, ammount)
		}
	}
	return finance.NewBalance(ammounts)
}

var TagRegex = regexp.MustCompile(`^(?P<name>[a-zA-Z0-9\-\_]+):(?P<value>[a-zA-Z0-9\-\_]+)$`)

func TextToTag(s string) (journal.Tag, error) {
	match := TagRegex.FindStringSubmatch(strings.TrimSpace(s))
	if len(match) != 3 {
		return journal.Tag{}, fmt.Errorf("invalid tag: %s", s)
	}
	return journal.Tag{
		Name:  match[1],
		Value: match[2],
	}, nil
}

// DoneSource represents the possible sources of value when an user is done entering
// and input
type DoneSource string

const (
	Context DoneSource = "context"
	Input   DoneSource = "input"
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
