package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/context"
	"github.com/vitorqb/addledger/internal/display/widgets"
	"github.com/vitorqb/addledger/internal/eventbus"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/journal"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
)

func TestDescriptionPicker(t *testing.T) {
	type testcontext struct {
		descPicker *widgets.ContextualList
		state      *statemod.State
		eventbus   *eventbusmod.EventBus
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Loads descriptions from state",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, 2, c.descPicker.GetItemCount())
				assert.Equal(t, "Description Two", c.state.InputMetadata.SelectedDescription())
				err := c.eventbus.Send(eventbus.Event{
					Topic: "input.description.listaction",
					Data:  listaction.NEXT,
				})
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, "Description One", c.state.InputMetadata.SelectedDescription())
			},
		},
		{
			name: "Load description from CurrentStatement if found",
			run: func(t *testing.T, c *testcontext) {
				sEntries := []statementreader.StatementEntry{{Description: "Statement Description"}}
				c.state.SetStatementEntries(sEntries)
				c.descPicker.Refresh()
				assert.Equal(t, 3, c.descPicker.GetItemCount())
				assert.Equal(t, "Statement Description", c.state.InputMetadata.SelectedDescription())
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var err error
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.state.JournalMetadata.SetTransactions([]journal.Transaction{
				{Description: "Description One"},
				{Description: "Description Two"},
			})
			c.eventbus = eventbusmod.New()
			c.descPicker, err = NewDescriptionPicker(c.state, c.eventbus)
			if err != nil {
				t.Fatal(err)
			}
			testcase.run(t, c)
		})
	}
}
