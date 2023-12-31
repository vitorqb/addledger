// uerinput is responsible for managing data inputted by the user.
package userinput

import "github.com/vitorqb/addledger/internal/state"

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
