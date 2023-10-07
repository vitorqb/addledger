package widgets_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/listaction"
	mock_eventbus "github.com/vitorqb/addledger/mocks/eventbus"
)

var enterEventKey = tcell.NewEventKey(tcell.KeyEnter, 'e', tcell.ModNone)
var ctrlJEventKey = tcell.NewEventKey(tcell.KeyCtrlJ, 'j', tcell.ModCtrl)
var downEventKey = tcell.NewEventKey(tcell.KeyDown, 'd', tcell.ModNone)
var upEventKey = tcell.NewEventKey(tcell.KeyUp, 'u', tcell.ModNone)
var fakeSetFocus = func(p tview.Primitive) {}

// Mock for ContextualListLinkOpts
type ContextualListLinkMock struct {
	onListActionsCalls       []listaction.ListAction
	onDoneCalls              []input.DoneSource
	onInsertFromContextCalls []struct{}
}

func (c *ContextualListLinkMock) GetLinkOpts() ContextualListLinkOpts {
	return ContextualListLinkOpts{
		InputName: "test",
		OnListAction: func(la listaction.ListAction) {
			c.onListActionsCalls = append(c.onListActionsCalls, la)
		},
		OnDone: func(ds input.DoneSource) {
			c.onDoneCalls = append(c.onDoneCalls, ds)
		},
		OnInsertFromContext: func() {
			c.onInsertFromContextCalls = append(c.onInsertFromContextCalls, struct{}{})
		},
	}
}

func TestLinkContextualList(t *testing.T) {
	type testcontext struct {
		t                  *testing.T
		inputField         *InputField
		contextualLinkMock *ContextualListLinkMock
		eventbus           *mock_eventbus.MockIEventBus
	}
	type testcase struct {
		name  string
		run   func(*testcontext)
		setup func(*testcontext)
	}
	var testcases = []testcase{
		{
			name: "should call OnListAction when the user presses key up",
			run: func(c *testcontext) {
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				c.inputField.InputHandler()(upEventKey, fakeSetFocus)
				assert.Equal(c.t, 1, len(c.contextualLinkMock.onListActionsCalls))
				assert.Equal(c.t, listaction.PREV, c.contextualLinkMock.onListActionsCalls[0])
			},
		},
		{
			name: "should call OnListAction when the user presses key down",
			run: func(c *testcontext) {
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				c.inputField.InputHandler()(downEventKey, fakeSetFocus)
				assert.Equal(c.t, 1, len(c.contextualLinkMock.onListActionsCalls))
				assert.Equal(c.t, listaction.NEXT, c.contextualLinkMock.onListActionsCalls[0])
			},
		},
		{
			name: "should call OnDone when the user presses enter",
			run: func(c *testcontext) {
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				c.inputField.InputHandler()(enterEventKey, fakeSetFocus)
				assert.Equal(c.t, 1, len(c.contextualLinkMock.onDoneCalls))
				assert.Equal(c.t, input.Context, c.contextualLinkMock.onDoneCalls[0])
			},
		},
		{
			name: "should call OnDone when the user presses ctrl+j",
			run: func(c *testcontext) {
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				c.inputField.InputHandler()(ctrlJEventKey, fakeSetFocus)
				assert.Equal(c.t, 1, len(c.contextualLinkMock.onDoneCalls))
				assert.Equal(c.t, input.Input, c.contextualLinkMock.onDoneCalls[0])
			},
		},
		{
			name: "should call OnInsertFromContext when the user presses tab",
			run: func(c *testcontext) {
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				c.inputField.InputHandler()(tcell.NewEventKey(tcell.KeyTab, 't', tcell.ModNone), fakeSetFocus)
				assert.Equal(c.t, 1, len(c.contextualLinkMock.onInsertFromContextCalls))
			},
		},
		{
			name:  "should subscribe to eventbus topic and set text when event is received",
			setup: func(c *testcontext) {}, // Override default setup
			run: func(c *testcontext) {
				var subscription eventbusmod.Subscription
				c.eventbus.EXPECT().Subscribe(gomock.Any()).Do(func(s eventbusmod.Subscription) {
					subscription = s
				})
				listLinkOpts := c.contextualLinkMock.GetLinkOpts()
				c.inputField.LinkContextualList(c.eventbus, listLinkOpts)
				assert.Equal(c.t, "input.test.settext", subscription.Topic)
				subscription.Handler(eventbusmod.Event{Data: "FOO"})
				assert.Equal(c.t, "FOO", c.inputField.GetText())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tviewField := tview.NewInputField()
			inputField := &InputField{InputField: tviewField}
			contextualLinkMock := &ContextualListLinkMock{}
			eventbus := mock_eventbus.NewMockIEventBus(ctrl)
			c := &testcontext{
				t:                  t,
				inputField:         inputField,
				contextualLinkMock: contextualLinkMock,
				eventbus:           eventbus,
			}

			if tc.setup != nil {
				tc.setup(c)
			} else {
				eventbus.EXPECT().Subscribe(gomock.Any()).Times(1)
			}
			tc.run(c)
		})
	}
}
