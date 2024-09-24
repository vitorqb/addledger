package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	. "github.com/vitorqb/addledger/mocks/display"
)

func TestLoadStatementModal(t *testing.T) {
	type testcontext struct {
		controller *MockLoadStatementModalController
		state      *MockState
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
				c.state.EXPECT().DefaultCsvFile().Return("/foo")
				loadStatementModal := NewLoadStatementModal(c.controller, c.state)
				csvFileField := loadStatementModal.GetCsvInput()
				csvFileField.SetText("file")
				presetField := loadStatementModal.GetPresetInput()
				presetField.SetText("preset")
				enterEvent := tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone)
				loadStatementModal.GetButton(0).InputHandler()(enterEvent, func(tview.Primitive) {})

			},
		},
		{
			name: "Uses default path from state",
			run: func(t *testing.T, c *testcontext) {
				c.state.EXPECT().DefaultCsvFile().Return("/foo")
				loadStatementModal := NewLoadStatementModal(c.controller, c.state)
				csvFileField := loadStatementModal.GetCsvInput()
				assert.Equal(t, csvFileField.GetText(), "/foo")
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.controller = NewMockLoadStatementModalController(ctrl)
			c.state = NewMockState(ctrl)
			tc.run(t, c)
		})
	}
}
