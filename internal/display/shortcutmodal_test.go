package display_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	. "github.com/vitorqb/addledger/internal/display"
	. "github.com/vitorqb/addledger/mocks/display"
)

func TestShortcutModal(t *testing.T) {
	type testcase struct {
		name       string
		run        func(t *testing.T, tc *testcase)
		modal      *ShortcutModal
		controller *MockShortcutModalController
	}
	var testExitOnKey = func(k tcell.Key, r rune, mod tcell.ModMask) func(t *testing.T, tc *testcase) {
		return func(t *testing.T, tc *testcase) {
			event := tcell.NewEventKey(k, r, tcell.ModNone)
			setFocus := func(tview.Primitive) {}
			tc.controller.EXPECT().OnHideShortcutModal().Times(1)
			tc.modal.InputHandler()(event, setFocus)
		}
	}
	var testcases = []testcase{
		{
			name: "Calls controller exit on escape key",
			run:  testExitOnKey(tcell.KeyEscape, ' ', tcell.ModNone),
		},
		{
			name: "Calls controller exit on q",
			run:  testExitOnKey(tcell.KeyRune, 'q', tcell.ModNone),
		},
		{
			name: "Calls controller exit on ctrl+q",
			run:  testExitOnKey(tcell.KeyCtrlQ, 'q', tcell.ModNone),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tc.controller = NewMockShortcutModalController(ctrl)
			tc.modal = NewShortcutModal(tc.controller)
			tc.run(t, &tc)
		})
	}
}
