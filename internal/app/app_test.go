package app_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
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

func TestConfigureLogger(t *testing.T) {

	t.Run("Logs to a file", func(t *testing.T) {
		logger := logrus.New()
		tempDir := t.TempDir()
		logFile := tempDir + "/log"
		ConfigureLogger(logger, logFile, "info")
		logger.Info("foo")

		file, err := os.Open(logFile)
		assert.Nil(t, err)
		defer file.Close()

		fileContentBytes, err := os.ReadFile(logFile)
		assert.Nil(t, err)
		fileContent := string(fileContentBytes)
		assert.Contains(t, fileContent, "foo")
	})

}
