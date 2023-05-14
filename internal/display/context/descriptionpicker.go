package context

import (
	"fmt"

	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	statemod "github.com/vitorqb/addledger/internal/state"
)

// DescriptionPicker presents a list of known descriptions to the user,
// and allows it to pick one.
func NewDescriptionPicker(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) (*widgets.ContextualList, error) {
	list := widgets.NewContextualList(
		func() []string {
			descriptionsMap := make(map[string]interface{})
			for _, p := range state.JournalMetadata.Transactions() {
				descriptionsMap[p.Description] = 1
			}
			var descriptions []string
			for k := range descriptionsMap {
				descriptions = append(descriptions, k)
			}
			return descriptions
		},
		func(s string) {
			state.InputMetadata.SetSelectedDescription(s)
		},
		func() string {
			return state.InputMetadata.DescriptionText()
		},
	)
	state.AddOnChangeHook(list.Refresh)
	err := eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.description.listaction",
		Handler: list.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return list, nil
}
