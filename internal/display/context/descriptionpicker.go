package context

import (
	"fmt"

	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/utils"
)

// DescriptionPicker presents a list of known descriptions to the user,
// and allows it to pick one.
func NewDescriptionPicker(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
) (*widgets.ContextualList, error) {
	list, err := widgets.NewContextualList(widgets.ContextualListOptions{
		GetItemsFunc: func() []string {
			descriptions := []string{}
			for _, transaction := range state.JournalMetadata.Transactions() {
				descriptions = append(descriptions, transaction.Description)
			}
			out := utils.Unique(descriptions)
			// NOTE: call reverse so last transactions will be suggested first.
			utils.Reverse(out)
			return out
		},
		SetSelectedFunc: func(s string) {
			state.InputMetadata.SetSelectedDescription(s)
		},
		GetInputFunc: func() string {
			return state.InputMetadata.DescriptionText()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build contextual list: %w", err)
	}
	state.AddOnChangeHook(list.Refresh)
	err = eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.description.listaction",
		Handler: list.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return list, nil
}
