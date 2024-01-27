package context

import (
	"fmt"

	"time"

	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/utils"
	"github.com/vitorqb/addledger/pkg/delay"
)

// DescriptionPicker presents a list of known descriptions to the user,
// and allows it to pick one.
func NewDescriptionPicker(
	state *statemod.State,
	eventbus eventbusmod.IEventBus,
	app TviewApp,
) (*widgets.ContextualList, error) {
	list, err := widgets.NewContextualList(widgets.ContextualListOptions{
		GetItemsFunc: func() []string {
			descriptions := []string{}
			for _, transaction := range state.JournalMetadata.Transactions() {
				descriptions = append(descriptions, transaction.Description)
			}
			out := utils.Unique(descriptions)

			// NOTE: add current statement description if it exists.
			if sEntry, found := state.CurrentStatementEntry(); found {
				out = append(out, sEntry.Description)
			}

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
	// refreshes the list 500ms after the last change, and redraw the app/screen.
	refresher := delay.NewFunction(func() {
		app.QueueUpdateDraw(func() {
			list.Refresh()
		})
	}, 500*time.Millisecond)
	state.AddOnChangeHook(refresher.Schedule)
	err = eventbus.Subscribe(eventbusmod.Subscription{
		Topic:   "input.description.listaction",
		Handler: list.HandleActionFromEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to eventBus: %w", err)
	}
	return list, nil
}
