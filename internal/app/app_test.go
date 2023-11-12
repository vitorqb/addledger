package app_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/app"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementloader"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/statementloader"
)

func TestLoadStatement(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		loader := NewMockStatementLoader(ctrl)
		statementEntries := []statementloader.StatementEntry{
			{Account: "ACC", Description: "FOO"},
			{Account: "ACC", Description: "BAR"},
		}
		loader.EXPECT().Load(gomock.Any()).Return(statementEntries, nil)
		filePath := testutils.TestDataPath(t, "file")
		state := statemod.InitialState()
		err := LoadStatement(loader, filePath, state)
		assert.Nil(t, err)
		assert.Equal(t, statementEntries, state.GetStatementEntries())
	})

}
