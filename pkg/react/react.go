// react implements a reactifull state, which calls hooks on change.
package react

// OnChangeHook is a hook called when the state changes
type OnChangeHook func()

// IReact allows having a reactiful state that notifies subscribers on change.
type IReact interface {
	AddOnChangeHook(h OnChangeHook)
	NotifyChange()
}

// React implements IReact
type React struct {
	onChangeHooks []OnChangeHook
}

var _ IReact = &React{}

// AddOnChangeHook allows adding a new hook to be called on changes.
func (r *React) AddOnChangeHook(h OnChangeHook) {
	r.onChangeHooks = append(r.onChangeHooks, h)
}

// NotifyChange must be called everytime the state changes to notify subscribers.
func (r *React) NotifyChange() {
	for _, h := range r.onChangeHooks {
		h()
	}
}

// New creates a new React
func New() *React {
	return &React{[]OnChangeHook{}}
}
