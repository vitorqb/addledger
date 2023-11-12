package context_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/context"
	"github.com/vitorqb/addledger/internal/finance"
	statemod "github.com/vitorqb/addledger/internal/state"
)

var anAmmount = finance.Ammount{Commodity: "EUR", Quantity: decimal.New(1220, -2)}

func TestAmmountGuesser(t *testing.T) {
	type testcontext struct {
		ammountGuesser *AmmountGuesser
		state          *statemod.State
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Loads guess from state",
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.SetPostingAmmountGuess(anAmmount)
				text := c.ammountGuesser.GetText(true)
				assert.Equal(t, "EUR 12.2", text)
			},
		},
		{
			name: "Loads guess from state",
			run: func(t *testing.T, c *testcontext) {
				c.state.InputMetadata.ClearPostingAmmountGuess()
				text := c.ammountGuesser.GetText(true)
				assert.Equal(t, "", text)
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var err error
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.ammountGuesser, err = NewAmmountGuesser(c.state)
			if err != nil {
				t.Fatal(err)
			}
			testcase.run(t, c)
		})
	}
}
