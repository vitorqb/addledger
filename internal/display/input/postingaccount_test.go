package input_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/input"
	"github.com/vitorqb/addledger/internal/display/widgets"
	eventbusmod "github.com/vitorqb/addledger/internal/eventbus"
	"github.com/vitorqb/addledger/internal/listaction"
	"github.com/vitorqb/addledger/internal/userinput"
	. "github.com/vitorqb/addledger/mocks/controller"
)

func TestNewPostingAccountField(t *testing.T) {
	type testcontext struct {
		postingAccount *widgets.InputField
		controller     *MockIInputController
		eventbus       eventbusmod.IEventBus
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Sends next account when arrow down",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountListAcction(listaction.NEXT)
				event := tcell.NewEventKey(tcell.KeyDown, ' ', tcell.ModNone)
				c.postingAccount.InputHandler()(event, func(p tview.Primitive) {})
			},
		},
		{
			name: "Sends msg with current text to controller",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountChanged("FOO")
				c.postingAccount.SetText("FOO")
			},
		},
		{
			name: "Sends OnPostingAccountSelectedFromContext msg",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountChanged("FOO")
				c.controller.EXPECT().OnPostingAccountDone(userinput.Context)
				c.postingAccount.SetText("FOO")
				event := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				c.postingAccount.InputHandler()(event, func(tview.Primitive) {})
			},
		},
		{
			name: "Sends OnPostingAccountInsertFromContext msg",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountChanged("FOO")
				c.controller.EXPECT().OnPostingAccountInsertFromContext()
				c.postingAccount.SetText("FOO")
				event := tcell.NewEventKey(tcell.KeyTab, ' ', tcell.ModNone)
				c.postingAccount.InputHandler()(event, func(tview.Primitive) {})
			},
		},
		{
			name: "Set text from EventBus",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountChanged("FOO")
				event := eventbusmod.Event{
					Topic: "input.postingaccount.settext",
					Data:  "FOO",
				}
				err := c.eventbus.Send(event)
				assert.Nil(t, err)
				assert.Equal(t, "FOO", c.postingAccount.GetText())
			},
		},
		{
			name: "Calls OnPostingAccountDone on Ctrl+J",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnPostingAccountChanged("FOO")
				c.controller.EXPECT().OnPostingAccountDone(userinput.Input)
				c.postingAccount.SetText("FOO")
				event := tcell.NewEventKey(tcell.KeyCtrlJ, ' ', tcell.ModNone)
				c.postingAccount.InputHandler()(event, func(tview.Primitive) {})
			},
		},
		{
			name: "Dispatches to OnFinishPosting on ctrl+Enter",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnFinishPosting()
				event := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModAlt)
				c.postingAccount.InputHandler()(event, func(tview.Primitive) {})
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = NewMockIInputController(ctrl)
			c.eventbus = eventbusmod.New()
			c.postingAccount = NewPostingAccount(c.controller, c.eventbus)
			testcase.run(t, c)
		})
	}
}
