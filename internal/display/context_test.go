package display_test

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/testutils"
)

var expectedDate1String = "1993-11-23\nTue, 23 Nov 1993"

func TestNewDateGuesser(t *testing.T) {
	type testcontext struct {
		state   *state.State
		guesser *tview.TextView
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "sets text from state",
			run: func(c *testcontext, t *testing.T) {
				c.state.InputMetadata.SetDateGuess(testutils.Date1(t))
				assert.Equal(t, expectedDate1String, c.guesser.GetText(true))
			},
		},
		{
			name: "clears when state clears",
			run: func(c *testcontext, t *testing.T) {
				c.state.InputMetadata.SetDateGuess(testutils.Date1(t))
				c.state.InputMetadata.ClearDateGuess()
				assert.Equal(t, "", c.guesser.GetText(true))
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			c := new(testcontext)
			c.state = state.InitialState()
			c.guesser, err = NewDateGuesser(c.state)
			if err != nil {
				t.Fatal(err)
			}
			tc.run(c, t)
		})
	}
}
