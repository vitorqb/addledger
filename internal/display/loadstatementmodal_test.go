package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	. "github.com/vitorqb/addledger/internal/display"
	. "github.com/vitorqb/addledger/mocks/display"
)

func TestLoadStatementModal(t *testing.T) {
	type testcontext struct {
		loadStatementModal *LoadStatementModal
		controller         *MockLoadStatementModalController
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	var testcases = []testcase{
		{
			name: "Calls controller on load",
			run: func(t *testing.T, c *testcontext) {
				c.controller.EXPECT().OnLoadStatement("file", "preset").Times(1)
				csvFileField := c.loadStatementModal.GetCsvInput()
				csvFileField.SetText("file")
				presetField := c.loadStatementModal.GetPresetInput()
				presetField.SetText("preset")
				enterEvent := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				c.loadStatementModal.GetButton(0).InputHandler()(enterEvent, func(tview.Primitive) {})

			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = NewMockLoadStatementModalController(ctrl)
			c.loadStatementModal = NewLoadStatementModal(c.controller)
			tc.run(t, c)
		})
	}
}
