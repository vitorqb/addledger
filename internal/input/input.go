package input

import (
	"time"

	"github.com/vitorqb/addledger/internal/utils"
	"github.com/vitorqb/addledger/pkg/react"
)

type (
	JournalEntryInput struct {
		react.IReact
		inputs              map[string]interface{}
		currentPostingIndex int
	}
)

func NewJournalEntryInput() *JournalEntryInput {
	m := make(map[string]interface{})
	return &JournalEntryInput{react.New(), m, 0}
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

// CurrentPosting returns the current posting being edited.
func (i *JournalEntryInput) CurrentPosting() *PostingInput {
	// currentPostingIndex is on range -> get and return.
	if posting, found := i.GetPosting(i.currentPostingIndex); found {
		return posting
	}

	// currentPostingIndex out of range -> add one and return.
	posting := i.AddPosting()
	i.currentPostingIndex = i.CountPostings() - 1
	i.NotifyChange()
	return posting
}

// AdvancePosting is called when a posting has finished to be inputed,
// and we should advance the current posting. If it doesn't exist, adds it.
func (i *JournalEntryInput) AdvancePosting() {
	i.currentPostingIndex++
	if _, found := i.GetPosting(i.currentPostingIndex); !found {
		i.AddPosting()
	}
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

func (i *JournalEntryInput) GetPosting(index int) (*PostingInput, bool) {
	if postingsInputs, found := i.inputs["postings"]; found {
		if postingsInputs, ok := postingsInputs.([]*PostingInput); ok {
			if index <= len(postingsInputs)-1 {
				return postingsInputs[index], true
			}
		}
	}
	return NewPostingInput(), false
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

func (i *JournalEntryInput) DeleteCurrentPosting() {
	if rawPostings, found := i.inputs["postings"]; found {
		if postings, ok := rawPostings.([]*PostingInput); ok {
			if i.currentPostingIndex >= 0 && i.currentPostingIndex < len(postings) {
				postings = utils.RemoveIndex(i.currentPostingIndex, postings)
				i.inputs["postings"] = postings
				if i.currentPostingIndex > 0 && i.currentPostingIndex >= len(postings) {
					i.currentPostingIndex--
				}
				i.NotifyChange()
			}
		}
	}
}

// Repr transforms the JournalEntryInput into a string.
func (jei *JournalEntryInput) Repr() string {
	var text string
	if date, found := jei.GetDate(); found {
		text += date.Format("2006-01-02")
	}
	if description, found := jei.GetDescription(); found {
		text += " " + description
	}
	i := -1
	for {
		i++
		if posting, found := jei.GetPosting(i); found {
			text += "\n" + "    "
			if account, found := posting.GetAccount(); found {
				text += account
			}
			text += "    "
			if value, found := posting.GetValue(); found {
				text += value
			}
		} else {
			break
		}
	}
	return text
}
