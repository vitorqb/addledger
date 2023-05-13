package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ConfigValue represents a value for the configuration
type ConfigValue struct {
	flagName     string
	shorthand    string
	defaultValue string
	usage        string
}

type Config struct {
	// File to where we will write Journal Entries.
	DestFile string
	// LedgerFile to pass to `hledger` executable. Empty string means none.
	LedgerFile string
	// Executable path for hledger. Empty for "hledger".
	HLedgerExecutable string
	// File where to send log. Empty for stderr.
	LogFile string
	// Level for logging
	LogLevel string
}

// ConfigValues has all known config values
var ConfigValues = []ConfigValue{
	{"destfile", "d", "", "Destination file (where we will write)"},
	{"hledger-executable", "", "hledger", "Executable to use for HLedger"},
	{"ledger-file", "", "", "Ledger File to pass to HLedger commands"},
	{"logfile", "", "", "File where to send log output. Empty for stderr."},
	{"loglevel", "", "WARN", "Level of logger. Defaults to warning."},
}

func Setup(flagSet *pflag.FlagSet) {
	for _, configValue := range ConfigValues {
		if configValue.shorthand == "" {
			flagSet.String(configValue.flagName, configValue.defaultValue, configValue.usage)
		} else {
			flagSet.StringP(configValue.flagName, configValue.shorthand, configValue.defaultValue, configValue.usage)
		}
	}
}

func Load(flagSet *pflag.FlagSet, args []string) (*Config, error) {
	// Parse flags
	err := flagSet.Parse(args)
	if err != nil {
		return &Config{}, fmt.Errorf("error parsing flags: %w", err)
	}

	// Send to viper
	err = viper.BindPFlags(flagSet)
	if err != nil {
		return &Config{}, fmt.Errorf("error binding to viper: %w", err)
	}

	// Set reading from env
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("ADDLEDGER")
	for _, x := range ConfigValues {
		err := viper.BindEnv(x.flagName)
		if err != nil {
			return &Config{}, fmt.Errorf("failed to bind env: %w", err)
		}
	}

	// Unpack
	config := &Config{
		DestFile:          viper.GetString("destfile"),
		HLedgerExecutable: viper.GetString("hledger-executable"),
		LedgerFile:        viper.GetString("ledger-file"),
		LogFile:           viper.GetString("logfile"),
		LogLevel:          viper.GetString("loglevel"),
	}

	// Validate
	if config.DestFile == "" {
		return config, fmt.Errorf("missing destination file!")
	}

	return config, nil
}

func LoadFromCommandLine() (*Config, error) {
	Setup(pflag.CommandLine)
	return Load(pflag.CommandLine, os.Args)
}
