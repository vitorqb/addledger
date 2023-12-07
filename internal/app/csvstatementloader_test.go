package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/app"
	"github.com/vitorqb/addledger/internal/config"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
)

func TestMaybeLoadCsvStatement(t *testing.T) {
	t.Run("No file", func(t *testing.T) {
		state := statemod.InitialState()
		loader := NewCSVStatementLoader(state)
		err := loader.Load(config.CSVStatementLoaderConfig{})
		assert.Nil(t, err)
		assert.Equal(t, []statementreader.StatementEntry{}, state.GetStatementEntries())
	})
	t.Run("Fails to load statement", func(t *testing.T) {
		state := statemod.InitialState()
		loader := NewCSVStatementLoader(state)
		err := loader.Load(config.CSVStatementLoaderConfig{File: "dont-exist"})
		assert.ErrorContains(t, err, "failed to open file")
	})
}
