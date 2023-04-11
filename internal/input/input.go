package input

import "time"

type (
	OnChangeHook      func()
	JournalEntryInput struct {
		onChangeHooks []OnChangeHook
		inputs        map[string]interface{}
	}
)

func NewJournalEntryInput() *JournalEntryInput {
	m := make(map[string]interface{})
	ws := []OnChangeHook{}
	return &JournalEntryInput{ws, m}
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

func (i *JournalEntryInput) AddPosting() {
	postInput := NewPostingInput()
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
}
