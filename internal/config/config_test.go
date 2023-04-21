package config_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/config"
	"github.com/vitorqb/addledger/internal/testutils"
)

func TestLoad(t *testing.T) {
	t.Run("Missing destfile", func(t *testing.T) {
		cleanup := testutils.Unsetenv(t, "ADDLEDGER_DESTFILE")
		defer cleanup()
		flagSet := pflag.NewFlagSet("foo", pflag.ContinueOnError)
		Setup(flagSet)
		_, err := Load(flagSet, []string{})
		assert.ErrorContains(t, err, "missing destination file")
	})
	t.Run("Working example", func(t *testing.T) {
		cleanup := testutils.Unsetenv(t, "ADDLEDGER_DESTFILE")
		defer cleanup()
		flagSet := pflag.NewFlagSet("foo", pflag.ContinueOnError)
		Setup(flagSet)
		config, err := Load(flagSet, []string{"-dfoo"})
		assert.Nil(t, err)
		assert.Equal(t, config.DestFile, "foo")
	})
	t.Run("Working example from env", func(t *testing.T) {
		cleanup := testutils.Setenv(t, "ADDLEDGER_DESTFILE", "foo")
		defer cleanup()
		flagSet := pflag.NewFlagSet("foo", pflag.ContinueOnError)
		Setup(flagSet)
		config, err := Load(flagSet, []string{})
		assert.Nil(t, err)
		assert.Equal(t, config.DestFile, "foo")
	})
}
