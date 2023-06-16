package input_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
)

var aDecimal, _ = decimal.NewFromString("2.20")
var anAmmount = journal.Ammount{Commodity: "EUR", Quantity: aDecimal}
var anotherAmmount = journal.Ammount{Commodity: "EUR", Quantity: decimal.New(-22, -1)}

func TestPostingInput(t *testing.T) {

	type context struct {
		onChangeCallCount int
		postingInput      *PostingInput
	}

	type test struct {
		name string
		run  func(t *testing.T, c *context)
	}

	testcases := []test{
		{
			"Set and clear Account",
			func(t *testing.T, c *context) {
				_, found := c.postingInput.GetAccount()
				assert.False(t, found)
				c.postingInput.SetAccount("FOO")
				account, found := c.postingInput.GetAccount()
				assert.True(t, found)
				assert.Equal(t, "FOO", account)
				assert.Equal(t, 1, c.onChangeCallCount)
				c.postingInput.ClearAccount()
				_, found = c.postingInput.GetAccount()
				assert.False(t, found)
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
		{
			"Set and clear Ammount",
			func(t *testing.T, c *context) {
				_, found := c.postingInput.GetAmmount()
				assert.False(t, found)

				c.postingInput.SetAmmount(anAmmount)

				value, found := c.postingInput.GetAmmount()
				assert.True(t, found)
				assert.Equal(t, anAmmount, value)
				assert.Equal(t, 1, c.onChangeCallCount)

				c.postingInput.ClearAmmount()

				_, found = c.postingInput.GetAmmount()
				assert.False(t, found)
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			postingInput := NewPostingInput()
			c := &context{
				onChangeCallCount: 0,
				postingInput:      postingInput,
			}
			postingInput.AddOnChangeHook(func() { c.onChangeCallCount++ })
			tc.run(t, c)
		})
	}
}
