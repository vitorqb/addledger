package config_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/testutils"
)

func TestLoad(t *testing.T) {

	type testcontext struct {
		flagSet *pflag.FlagSet
	}

	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}

	testcases := []testcase{
		{
			name: "Missing destfile",
			run: func(t *testing.T, c *testcontext) {
				_, err := Load(c.flagSet, []string{})
				assert.ErrorContains(t, err, "missing destination file")
			},
		},
		{
			name: "Minimal working example",
			run: func(t *testing.T, c *testcontext) {
				config, err := Load(c.flagSet, []string{"-dfoo"})
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
				config, err := Load(c.flagSet, []string{})
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
			},
		},
		{
			name: "Full working example",
			run: func(t *testing.T, c *testcontext) {
				config, err := Load(c.flagSet, []string{
					"-dfoo",
					"--ledger-file=bar",
					"--hledger-executable=baz",
				})
				assert.Nil(t, err)
				assert.Equal(t, config.DestFile, "foo")
				assert.Equal(t, config.HLedgerExecutable, "baz")
				assert.Equal(t, config.LedgerFile, "bar")
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
				config, err := Load(c.flagSet, []string{})
				assert.Nil(t, err)
				assert.Equal(t, "foo1", config.DestFile)
				assert.Equal(t, "foo2", config.HLedgerExecutable)
				assert.Equal(t, "foo3", config.LedgerFile)
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
			Setup(c.flagSet)
			testcase.run(t, c)
		})
	}
}
