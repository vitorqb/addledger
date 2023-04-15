package input

import "time"

type (
	OnChangeHook      func()
	JournalEntryInput struct {
		onChangeHooks       []OnChangeHook
		inputs              map[string]interface{}
		currentPostingIndex int
	}
)

func NewJournalEntryInput() *JournalEntryInput {
	m := make(map[string]interface{})
	ws := []OnChangeHook{}
	return &JournalEntryInput{ws, m, 0}
}

func (i *JournalEntryInput) AddOnChangeHook(hook OnChangeHook) {
	i.onChangeHooks = append(i.onChangeHooks, hook)
}
func (i *JournalEntryInput) notifyChange() {
	for _, h := range i.onChangeHooks {
		h()
	}
}

func (i *JournalEntryInput) SetDate(x time.Time) {
	i.inputs["date"] = x
	i.notifyChange()

}
func (i *JournalEntryInput) GetDate() (time.Time, bool) {
	if rawValue, found := i.inputs["date"]; found {
		if value, ok := rawValue.(time.Time); ok {
			return value, true
		}
	}
	return time.Time{}, false
}

func (i *JournalEntryInput) SetDescription(x string) {
	i.inputs["description"] = x
	i.notifyChange()
}
func (i *JournalEntryInput) GetDescription() (string, bool) {
	if rawValue, found := i.inputs["description"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
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
	i.notifyChange()
	return posting
}

// AdvancePosting is called when a posting has finished to be inputed,
// and we should advance the current posting.
func (i *JournalEntryInput) AdvancePosting() {
	i.currentPostingIndex++
	i.notifyChange()
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
	postInput.AddOnChangeHook(i.notifyChange)
	if rawPostings, found := i.inputs["postings"]; found {
		if postings, ok := rawPostings.([]*PostingInput); ok {
			i.inputs["postings"] = append(postings, postInput)
			i.notifyChange()
			return
		}
	}
	i.inputs["postings"] = []*PostingInput{postInput}
	i.notifyChange()
	return
}
