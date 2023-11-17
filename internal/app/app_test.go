package app_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/app"
	"github.com/vitorqb/addledger/internal/config"
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

	t.Run("Fails to open log file", func(t *testing.T) {
		logger, logHook := logrusTest.NewNullLogger()
		logger.ExitFunc = func(int) {}
		logFile := "/foo/bar"
		ConfigureLogger(logger, logFile, "info")
		assert.Equal(t, 1, len(logHook.Entries))
		assert.Equal(t, logrus.FatalLevel, logHook.LastEntry().Level)
		assert.Contains(t, logHook.LastEntry().Message, "Failed to open log file")
	})

	t.Run("Fails to parse log level", func(t *testing.T) {
		logger, logHook := logrusTest.NewNullLogger()
		logger.ExitFunc = func(int) {}
		ConfigureLogger(logger, "", "foo")
		assert.Equal(t, 1, len(logHook.Entries))
		assert.Equal(t, logrus.FatalLevel, logHook.LastEntry().Level)
		assert.Contains(t, logHook.LastEntry().Message, "Failed to parse log level")
	})
}

func TestMaybeLoadCsvStatement(t *testing.T) {
	t.Run("No file", func(t *testing.T) {
		state := statemod.InitialState()
		err := MaybeLoadCsvStatement(config.CSVStatementLoaderConfig{}, state)
		assert.Nil(t, err)
		assert.Equal(t, []statementloader.StatementEntry{}, state.GetStatementEntries())
	})
	t.Run("Fails to load statement", func(t *testing.T) {
		state := statemod.InitialState()
		err := MaybeLoadCsvStatement(config.CSVStatementLoaderConfig{File: "dont-exist"}, state)
		assert.ErrorContains(t, err, "failed to load statement")
	})
}
