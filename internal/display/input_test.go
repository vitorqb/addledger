package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	. "github.com/vitorqb/addledger/internal/display"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
	mock_controller "github.com/vitorqb/addledger/mocks/controller"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

func TestDescriptionField(t *testing.T) {
	type testcontext struct {
		controller *mock_controller.MockIInputController
		eventbus   *mock_eventbus.MockIEventBus
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Call controller on change",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionChanged("FOO")
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)
				field.SetText("FOO")
			},
		},
		{
			name: "Dispatches NEXT to context list",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionListAction(listaction.NEXT)
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
				field.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Dispatches PREV to context list",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionListAction(listaction.PREV)
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
				field.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Dispatches select from context to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionSelectedFromContext()
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				field.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Dispatches insert from context to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionInsertFromContext()
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyTAB, ' ', tcell.ModNone)
				field.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Dispatches done to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionDone()
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyCtrlJ, ' ', tcell.ModNone)
				field.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Set text from topic",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionChanged("FOO")
				var subscription eventbusmod.Subscription
				c.eventbus.
					EXPECT().
					Subscribe(gomock.Any()).
					Do(func(s eventbusmod.Subscription) {
						subscription = s
					})
				DescriptionField(c.controller, c.eventbus)
				event := eventbusmod.Event{Topic: "foo", Data: "FOO"}
				subscription.Handler(event)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.eventbus = mock_eventbus.NewMockIEventBus(ctrl)
			tc.run(c, t)
		})
	}
}
