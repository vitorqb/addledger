package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	statemod "github.com/vitorqb/addledger/internal/state"
	mock_accountguesser "github.com/vitorqb/addledger/mocks/accountguesser"
	mock_controller "github.com/vitorqb/addledger/mocks/controller"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

func TestNewLayout(t *testing.T) {
	type testcontext struct {
		controller     *mock_controller.MockIInputController
		state          *statemod.State
		eventbus       *mock_eventbus.MockIEventBus
		accountGuesser *mock_accountguesser.MockIAccountGuesser
		layout         *Layout
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Handles CTRL+Z",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnUndo().Times(1)
				key := tcell.KeyCtrlZ
				event := tcell.NewEventKey(key, 'z', tcell.ModCtrl)
				setFocus := func(tview.Primitive) {}
				c.layout.InputHandler()(event, setFocus)
			},
		},
		{
			name: "Displays the tag picker",
			run: func(c *testcontext, t *testing.T) {
				c.state.SetPhase(statemod.InputTags)
				_, page := c.layout.GetItem(3).(*tview.Pages).GetFrontPage()
				assert.IsType(t, &TagsPicker{}, page)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.eventbus = mock_eventbus.NewMockIEventBus(ctrl)
			c.accountGuesser = mock_accountguesser.NewMockIAccountGuesser(ctrl)
			c.state = statemod.InitialState()
			// Subscribe is called multiple times for each subscription
			// that happens in the entire layout.
			c.eventbus.EXPECT().Subscribe(gomock.Any()).AnyTimes()
			// Some controller methods are called on startup
			c.controller.EXPECT().OnDateChanged("")
			// AccountGuesser is called on startup
			c.accountGuesser.EXPECT().Guess().AnyTimes()
			c.layout, err = NewLayout(c.controller, c.state, c.eventbus, c.accountGuesser)
			if err != nil {
				t.Fatalf("Failed to create layout: %s", err)
			}
			tc.run(c, t)
		})
	}
}
