package input_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	. "github.com/vitorqb/addledger/internal/display/input"
	"github.com/vitorqb/addledger/internal/listaction"
	. "github.com/vitorqb/addledger/mocks/controller"
)

func TestPostingAccountField(t *testing.T) {
	type testcontext struct {
		postingAccount *PostingAccountField
		controller     *MockIInputController
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
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			c := new(testcontext)
			c.controller = NewMockIInputController(ctrl)
			c.postingAccount = NewPostingAccount(c.controller)
			testcase.run(t, c)
		})
	}
}
