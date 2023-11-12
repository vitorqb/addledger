package config_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/testutils"
)

type MockLoader struct{}

var _ ILoader = (*MockLoader)(nil)

func (l *MockLoader) JournalFile(_ string) (string, error) {
	return "/path/to/journal/from/mock", nil
}

func TestLoad(t *testing.T) {
	csvFile := testutils.TestDataPath(t, "statement.csv")
	fullCsvPresetFile := testutils.TestDataPath(t, "csv_preset_full.json")

	type testcontext struct {
		flagSet *pflag.FlagSet
		loader  ILoader
	}

	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "Minimal working example",
			run: func(t *testing.T, c *testcontext) {
				config, err := Load(c.flagSet, []string{"-dfoo"}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
				assert.Equal(t, config.HLedgerExecutable, "hledger")
				assert.Equal(t, config.LedgerFile, "")
			},
		},
		{
			name: "Working example from env",
			run: func(t *testing.T, c *testcontext) {
				cleanup := testutils.Setenv(t, "ADDLEDGER_DESTFILE", "foo")
				defer cleanup()
				config, err := Load(c.flagSet, []string{}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
			},
		},
		{
			name: "Full working example",
			run: func(t *testing.T, c *testcontext) {
				cleanup := testutils.Setenv(t, "ADDLEDGER_PRINTER_LINE_BREAK_AFTER", "4")
				defer cleanup()
				config, err := Load(c.flagSet, []string{
					"-dfoo",
					"--ledger-file=bar",
					"--hledger-executable=baz",
					"--printer-line-break-before=3",
				}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
				assert.Equal(t, config.HLedgerExecutable, "baz")
				assert.Equal(t, config.LedgerFile, "bar")
				assert.Equal(t, 3, config.PrinterConfig.NumLineBreaksBefore)
				assert.Equal(t, 4, config.PrinterConfig.NumLineBreaksAfter)
			},
		},
		{
			name: "Full working example from env",
			run: func(t *testing.T, c *testcontext) {
				cleanup := testutils.Setenvs(t,
					"ADDLEDGER_DESTFILE", "foo1",
					"ADDLEDGER_HLEDGER_EXECUTABLE", "foo2",
					"ADDLEDGER_LEDGER_FILE", "foo3",
				)
				defer cleanup()
				config, err := Load(c.flagSet, []string{}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, "foo1", config.DestFile)
				assert.Equal(t, "foo2", config.HLedgerExecutable)
				assert.Equal(t, "foo3", config.LedgerFile)
			},
		},
		{
			name: "With full csv statement config",
			run: func(t *testing.T, c *testcontext) {
				flags := []string{
					"--csv-statement-file=" + csvFile,
					"--csv-statement-preset=" + fullCsvPresetFile,
				}
				config, err := Load(c.flagSet, flags, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.CSVStatementLoaderConfig.File, csvFile)
				assert.Equal(t, config.CSVStatementLoaderConfig.Separator, ";")
				assert.Equal(t, config.CSVStatementLoaderConfig.Account, "acc")
				assert.Equal(t, config.CSVStatementLoaderConfig.DateFieldIndex, 0)
				assert.Equal(t, config.CSVStatementLoaderConfig.DateFormat, "01/02/2006")
				assert.Equal(t, config.CSVStatementLoaderConfig.DescriptionFieldIndex, 1)
				assert.Equal(t, config.CSVStatementLoaderConfig.AccountFieldIndex, 2)
				assert.Equal(t, config.CSVStatementLoaderConfig.AmmountFieldIndex, 3)
			},
		},
		{
			name: "Defaults DestFile to LedgerFile",
			run: func(t *testing.T, c *testcontext) {
				config, err := Load(c.flagSet, []string{"--ledger-file=foo"}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
			},
		},
		{
			name: "Defaults DestFile to LedgerFile from hledger executable",
			run: func(t *testing.T, c *testcontext) {
				config, err := Load(c.flagSet, []string{}, c.loader)
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "/path/to/journal/from/mock")
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			c := new(testcontext)
			cleanup := testutils.Unsetenvs(t,
				"ADDLEDGER_DESTFILE",
				"ADDLEDGER_HLEDGER_EXECUTABLE",
				"ADDLEDGER_LEDGER_FILE",
			)
			defer cleanup()
			c.flagSet = pflag.NewFlagSet("foo", pflag.ContinueOnError)
			c.loader = new(MockLoader)
			SetupFlags(c.flagSet)
			testcase.run(t, c)
		})
	}
}

func TestLoadCsvStatementLoaderConfig(t *testing.T) {
	csvFile := testutils.TestDataPath(t, "statement.csv")
	minPresetFile := testutils.TestDataPath(t, "csv_preset_min.json")
	fullPresetFile := testutils.TestDataPath(t, "csv_preset_full.json")

	t.Run("No file", func(t *testing.T) {
		config, err := LoadCsvStatementLoaderConfig("", "")
		assert.Equal(t, CSVStatementLoaderConfig{}, config)
		assert.NoError(t, err)
	})

	t.Run("No preset", func(t *testing.T) {
		config, err := LoadCsvStatementLoaderConfig(csvFile, "")
		assert.Equal(t, CSVStatementLoaderConfig{}, config)
		assert.ErrorContains(t, err, "missing preset")
	})

	t.Run("Preset not found", func(t *testing.T) {
		config, err := LoadCsvStatementLoaderConfig(csvFile, "foo")
		assert.Equal(t, CSVStatementLoaderConfig{}, config)
		assert.ErrorContains(t, err, "failed to open preset file")
	})

	t.Run("Minimal preset found", func(t *testing.T) {
		config, err := LoadCsvStatementLoaderConfig(csvFile, minPresetFile)
		assert.NoError(t, err)
		assert.Equal(t, CSVStatementLoaderConfig{
			File:                  csvFile,
			Separator:             "",
			Account:               "",
			Commodity:             "",
			DateFormat:            "02/01/2006",
			DateFieldIndex:        -1,
			DescriptionFieldIndex: -1,
			AccountFieldIndex:     -1,
			AmmountFieldIndex:     -1,
		}, config)
	})

	t.Run("Full preset found", func(t *testing.T) {
		config, err := LoadCsvStatementLoaderConfig(csvFile, fullPresetFile)
		assert.NoError(t, err)
		assert.Equal(t, CSVStatementLoaderConfig{
			File:                  csvFile,
			Separator:             ";",
			Account:               "acc",
			Commodity:             "com",
			DateFormat:            "01/02/2006",
			DateFieldIndex:        0,
			DescriptionFieldIndex: 1,
			AccountFieldIndex:     2,
			AmmountFieldIndex:     3,
		}, config)
	})
}
