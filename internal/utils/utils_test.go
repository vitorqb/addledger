package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/utils"
)

func TestSplitArray(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		next, err := SplitArray(2, []int{})
		assert.Nil(t, err)
		first, err := next()
		assert.ErrorIs(t, &StopSplitArray{}, err)
		assert.Nil(t, first)
	})
	t.Run("Two long", func(t *testing.T) {
		next, err := SplitArray(2, []int{1, 2, 3, 4})
		assert.Nil(t, err)
		first, err := next()
		assert.Nil(t, err)
		assert.Equal(t, []int{1, 2}, first)
		second, err := next()
		assert.Nil(t, err)
		assert.Equal(t, []int{3, 4}, second)
		_, err = next()
		assert.ErrorIs(t, &StopSplitArray{}, err)
	})
}
