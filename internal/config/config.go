package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	// File to where we will write Journal Entries.
	DestFile string
	// LedgerFile to pass to `hledger` executable. Empty string means none.
	LedgerFile string
	// Executable path for hledger. Empty for "hledger".
	HLedgerExecutable string
}

func Setup(flagSet *pflag.FlagSet) {
	flagSet.StringP("destfile", "d", "", "Destination file (where we will write)")
	flagSet.String("hledger-executable", "hledger", "Hledger executable")
	flagSet.String("ledger-file", "", "Ledger File to pass to HLedger commands")
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
	viper.SetEnvPrefix("ADDLEDGER")
	for _, x := range []string{"destfile", "hledger-executable", "ledger-file"} {
		err := viper.BindEnv(x)
		if err != nil {
			return &Config{}, fmt.Errorf("failed to bind env: %w", err)
		}
	}

	// Unpack
	config := &Config{
		DestFile:          viper.GetString("destfile"),
		HLedgerExecutable: viper.GetString("hledger-executable"),
		LedgerFile:        viper.GetString("ledger-file"),
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
