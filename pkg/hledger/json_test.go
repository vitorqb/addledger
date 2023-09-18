package hledger_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/pkg/hledger"
)

func TestJsonTagUnmarshall(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		source := `["tag1", "value1"]`
		var tag JSONTag
		err := json.Unmarshal([]byte(source), &tag)
		assert.Nil(t, err)
		assert.Equal(t, tag.Name, "tag1")
		assert.Equal(t, tag.Value, "value1")
	})
	t.Run("Invalid", func(t *testing.T) {
		source := `["tag1", "value1", "value2"]`
		var tag JSONTag
		err := json.Unmarshal([]byte(source), &tag)
		assert.NotNil(t, err)
	})
}
