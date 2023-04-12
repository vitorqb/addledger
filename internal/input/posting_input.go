package input

type (
	PostingInput struct {
		onChangeHooks []OnChangeHook
		inputs        map[string]interface{}
	}
)

func NewPostingInput() *PostingInput {
	onChangeHooks := []OnChangeHook{}
	inputs := map[string]interface{}{}
	return &PostingInput{onChangeHooks: onChangeHooks, inputs: inputs}
}

func (i *PostingInput) AddOnChangeHook(hook OnChangeHook) {
	i.onChangeHooks = append(i.onChangeHooks, hook)
}

func (i *PostingInput) notifyChange() {
	for _, h := range i.onChangeHooks {
		h()
	}
}

func (i *PostingInput) SetAccount(account string) {
	i.inputs["account"] = account
	i.notifyChange()
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
	i.notifyChange()
}

func (i *PostingInput) GetValue() (string, bool) {
	if rawValue, found := i.inputs["value"]; found {
		if value, ok := rawValue.(string); ok {
			return value, true
		}
	}
	return "", false
}
