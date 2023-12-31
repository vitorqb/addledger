package input

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/pkg/react"
)

type (
	JournalEntryInput struct {
		react.IReact
		inputs map[string]interface{}
	}
)

func NewJournalEntryInput() *JournalEntryInput {
	m := make(map[string]interface{})
	return &JournalEntryInput{react.New(), m}
}

func (i *JournalEntryInput) SetDate(x time.Time) {
	i.inputs["date"] = x
	i.NotifyChange()

}
func (i *JournalEntryInput) GetDate() (time.Time, bool) {
	if rawValue, found := i.inputs["date"]; found {
		if value, ok := rawValue.(time.Time); ok {
			return value, true
		}
	}
	return time.Time{}, false
}

func (i *JournalEntryInput) ClearDate() {
	delete(i.inputs, "date")
	i.NotifyChange()
}

func (i *JournalEntryInput) SetDescription(x string) {
	i.inputs["description"] = x
	i.NotifyChange()
}
func (i *JournalEntryInput) GetDescription() (string, bool) {
	if rawValue, found := i.inputs["description"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
}

func (i *JournalEntryInput) ClearDescription() {
	delete(i.inputs, "description")
	i.NotifyChange()
}

func (i *JournalEntryInput) GetTags() []journal.Tag {
	if rawValue, found := i.inputs["tags"]; found {
		if value, ok := rawValue.([]journal.Tag); ok {
			return value
		}
	}
	return []journal.Tag{}
}

func (i *JournalEntryInput) AppendTag(x journal.Tag) {
	tags := i.GetTags()
	i.inputs["tags"] = append(tags, x)
	i.NotifyChange()
}

func (i *JournalEntryInput) PopTag() {
	tags := i.GetTags()
	if len(tags) > 0 {
		i.inputs["tags"] = tags[:len(tags)-1]
		i.NotifyChange()
	}
}

func (i *JournalEntryInput) ClearTags() {
	delete(i.inputs, "tags")
	i.NotifyChange()
}

func (i *JournalEntryInput) CountPostings() int {
	if postingsInputs, found := i.inputs["postings"]; found {
		if postingsInputs, ok := postingsInputs.([]*PostingInput); ok {
			return len(postingsInputs)
		}
	}
	return 0
}

func (i *JournalEntryInput) AddPosting() (postInput *PostingInput) {
	postInput = NewPostingInput()
	postInput.AddOnChangeHook(i.NotifyChange)
	if rawPostings, found := i.inputs["postings"]; found {
		if postings, ok := rawPostings.([]*PostingInput); ok {
			i.inputs["postings"] = append(postings, postInput)
			i.NotifyChange()
			return
		}
	}
	i.inputs["postings"] = []*PostingInput{postInput}
	i.NotifyChange()
	return
}

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

func TagToText(t journal.Tag) string {
	return fmt.Sprintf("%s:%s", t.Name, t.Value)
}

// DoneSource represents the possible sources of value when an user is done entering
// and input
type DoneSource string

const (
	Context DoneSource = "context"
	Input   DoneSource = "input"
)
