package display_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display"
	"github.com/vitorqb/addledger/internal/finance"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
)

func TestStatementDisplay(t *testing.T) {
	type testcontext struct {
		statementDisplay *StatementDisplay
		state            *statemod.State
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "Displays statement when available",
			run: func(c *testcontext, t *testing.T) {
				assert.Equal(t, "", c.statementDisplay.GetText(false))
				c.state.SetStatementEntries([]statementloader.StatementEntry{
					{
						Account:     "ACC",
						Date:        time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
						Description: "FOO",
						Ammount: finance.Ammount{
							Commodity: "EUR",
							Quantity:  decimal.New(1221, -2),
						},
					},
				})
				assert.Equal(t, "2023/10/31 | FOO | ACC | EUR 12.21 | [1]", c.statementDisplay.GetText(false))
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.statementDisplay = NewStatementDisplay(c.state)
			tc.run(c, t)
		})
	}
}
