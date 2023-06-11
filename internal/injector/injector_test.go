package injector_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/ammountguesser"
	. "github.com/vitorqb/addledger/internal/injector"
	"github.com/vitorqb/addledger/internal/journal"
	statemod "github.com/vitorqb/addledger/internal/state"
)

func TestAmmountGuesserEngine(t *testing.T) {
	state := statemod.InitialState()
	_ = AmmountGuesserEngine(state)

	// At the beggining, default guess
	guess, found := state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	assert.Equal(t, ammountguesser.DefaultGuess, guess)

	// On new input for ammount guesser text, updates guess
	state.InputMetadata.SetPostingAmmountText("99.99")
	guess, found = state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	expectedGuess := journal.Ammount{
		Commodity: ammountguesser.DefaultCommodity,
		Quantity:  decimal.New(9999, -2),
	}
	assert.Equal(t, expectedGuess, guess)

	// On invalid input, defaults to default guess
	state.InputMetadata.SetPostingAmmountText("aaaa")
	guess, found = state.InputMetadata.GetPostingAmmountGuess()
	assert.True(t, found)
	assert.Equal(t, ammountguesser.DefaultGuess, guess)
}
