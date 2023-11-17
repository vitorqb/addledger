package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/listaction"
	statemod "github.com/vitorqb/addledger/internal/state"
	mock_controller "github.com/vitorqb/addledger/mocks/controller"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

var enterEventKey = tcell.NewEventKey(tcell.KeyEnter, 'e', tcell.ModNone)
var ctrlJEventKey = tcell.NewEventKey(tcell.KeyCtrlJ, 'j', tcell.ModCtrl)
var fakeSetFocus = func(p tview.Primitive) {}

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
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Dispatches PREV to context list",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionListAction(listaction.PREV)
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Dispatches select from context to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionDone(input.Context)
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Dispatches insert from context to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionInsertFromContext()
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyTAB, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Dispatches done to controller",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDescriptionDone(input.Input)
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				field := DescriptionField(c.controller, c.eventbus)

				event := tcell.NewEventKey(tcell.KeyCtrlJ, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
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

func TestPostingAmmountField(t *testing.T) {
	type testcontext struct {
		controller *mock_controller.MockIInputController
		eventbus   *mock_eventbus.MockIEventBus
		state      *statemod.State
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Calls controller when done (enter)",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnPostingAmmountChanged("EUR 12.20")
				c.controller.EXPECT().OnPostingAmmountDone(input.Context)
				field := PostingAmmountField(c.controller)
				field.SetText("EUR 12.20")
				field.InputHandler()(enterEventKey, fakeSetFocus)
			},
		},
		{
			name: "Calls controller when done (ctrl+j)",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnPostingAmmountChanged("EUR 12.20")
				c.controller.EXPECT().OnPostingAmmountDone(input.Input)
				field := PostingAmmountField(c.controller)
				field.SetText("EUR 12.20")
				field.InputHandler()(ctrlJEventKey, fakeSetFocus)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.eventbus = mock_eventbus.NewMockIEventBus(ctrl)
			c.state = statemod.InitialState()
			tc.run(c, t)
		})
	}
}

func TestDateField(t *testing.T) {
	type testcontext struct {
		controller *mock_controller.MockIInputController
		state      *statemod.State
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Call controller when done",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDateChanged("1993-11-23")
				c.controller.EXPECT().OnDateDone()
				dateField := DateField(c.controller)
				dateField.SetText("1993-11-23")
				dateField.InputHandler()(enterEventKey, fakeSetFocus)
			},
		},
		{
			name: "Call controller when change",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnDateChanged("1993-11-23")
				dateField := DateField(c.controller)
				dateField.SetText("1993-11-23")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.state = statemod.InitialState()
			tc.run(c, t)
		})
	}
}

func TestInput(t *testing.T) {
	type testcontext struct {
		controller *mock_controller.MockIInputController
		eventbus   *mock_eventbus.MockIEventBus
		state      *statemod.State
		input      *Input
	}
	type testcase struct {
		name  string
		setup func(c *testcontext, t *testing.T)
		run   func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Set up a page for Tags",
			run: func(c *testcontext, t *testing.T) {
				c.state.SetPhase(statemod.InputTags)
				_, page := c.input.GetContent().(*tview.Pages).GetFrontPage()
				field, ok := page.(*widgets.InputField)
				assert.True(t, ok)
				assert.Equal(t, "Tags: ", field.GetLabel())
			},
		},
		{
			name: "Tags come after description",
			run: func(c *testcontext, t *testing.T) {
				c.state.SetPhase(statemod.InputDescription)
				c.state.NextPhase()
				_, page := c.input.GetContent().(*tview.Pages).GetFrontPage()
				field, ok := page.(*widgets.InputField)
				assert.True(t, ok)
				assert.Equal(t, "Tags: ", field.GetLabel())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.eventbus = mock_eventbus.NewMockIEventBus(ctrl)
			c.state = statemod.InitialState()

			// Called during init
			c.controller.EXPECT().OnDateChanged(gomock.Any()).AnyTimes()
			c.eventbus.EXPECT().Subscribe(gomock.Any()).AnyTimes()

			if tc.setup != nil {
				tc.setup(c, t)
			}

			c.input = NewInput(c.controller, c.state, c.eventbus)

			tc.run(c, t)
		})
	}
}

func TestTagsField(t *testing.T) {
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
			name: "Moves contextual list down",
			run: func(c *testcontext, t *testing.T) {
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				c.controller.EXPECT().OnTagListAction(listaction.NEXT)
				field := NewTagsField(c.controller, c.eventbus)
				event := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Moves contextual list up",
			run: func(c *testcontext, t *testing.T) {
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				c.controller.EXPECT().OnTagListAction(listaction.PREV)
				field := NewTagsField(c.controller, c.eventbus)
				event := tcell.NewEventKey(tcell.KeyUp, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Calls controller when done (enter)",
			run: func(c *testcontext, t *testing.T) {
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				c.controller.EXPECT().OnTagDone(input.Context)
				field := NewTagsField(c.controller, c.eventbus)
				event := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Calls controller when done (ctrl+j)",
			run: func(c *testcontext, t *testing.T) {
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				c.controller.EXPECT().OnTagDone(input.Input)
				field := NewTagsField(c.controller, c.eventbus)
				event := tcell.NewEventKey(tcell.KeyCtrlJ, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Calls controller when insert from context",
			run: func(c *testcontext, t *testing.T) {
				c.eventbus.EXPECT().Subscribe(gomock.Any())
				c.controller.EXPECT().OnTagInsertFromContext()
				field := NewTagsField(c.controller, c.eventbus)
				event := tcell.NewEventKey(tcell.KeyTAB, ' ', tcell.ModNone)
				field.InputHandler()(event, fakeSetFocus)
			},
		},
		{
			name: "Set text from event",
			run: func(c *testcontext, t *testing.T) {
				c.controller.EXPECT().OnTagChanged("FOO")
				var subscription eventbusmod.Subscription
				c.eventbus.
					EXPECT().
					Subscribe(gomock.Any()).
					Do(func(s eventbusmod.Subscription) {
						subscription = s
					})
				field := NewTagsField(c.controller, c.eventbus)
				event := eventbusmod.Event{Topic: "foo", Data: "FOO"}
				subscription.Handler(event)
				assert.Equal(t, "FOO", field.GetText())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = mock_controller.NewMockIInputController(ctrl)
			c.eventbus = mock_eventbus.NewMockIEventBus(ctrl)
			tc.run(c, t)
		})
	}
}
