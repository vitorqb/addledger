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
	var testExpectOnKey = func(k tcell.Key, r rune, mod tcell.ModMask, expect func(*MockShortcutModalController)) func(t *testing.T, tc *testcase) {
		return func(t *testing.T, tc *testcase) {
			event := tcell.NewEventKey(k, r, tcell.ModNone)
			setFocus := func(tview.Primitive) {}
			expect(tc.controller)
			tc.modal.InputHandler()(event, setFocus)
		}
	}
	var testcases = []testcase{
		{
			name: "Calls controller exit on escape key",
			run: testExpectOnKey(tcell.KeyEscape, ' ', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
		},
		{
			name: "Calls controller exit on q",
			run: testExpectOnKey(tcell.KeyRune, 'q', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
		},
		{
			name: "Calls controller exit on ctrl+q",
			run: testExpectOnKey(tcell.KeyCtrlQ, 'q', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
		},
		{
			name: "Calls discard statement on d",
			run: testExpectOnKey(tcell.KeyRune, 'd', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnDiscardStatement().Times(1)
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
		},
		{
			name: "Calls load statement on l",
			run: testExpectOnKey(tcell.KeyRune, 'l', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnLoadStatementRequest().Times(1)
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
		},
		{
			name: "Calls show statement modal",
			run: testExpectOnKey(tcell.KeyRune, 's', tcell.ModNone, func(c *MockShortcutModalController) {
				c.EXPECT().OnShowStatementModal().Times(1)
				c.EXPECT().OnHideShortcutModal().Times(1)
			}),
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
