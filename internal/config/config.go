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
}

func Setup(flagSet *pflag.FlagSet) {
	flagSet.StringP("destfile", "d", "", "Destination file (where we will write)")
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
	err = viper.BindEnv("destfile")
	if err != nil {
		return &Config{}, fmt.Errorf("failed to bind env: %w", err)
	}

	// Unpack
	config := &Config{
		DestFile: viper.GetString("destfile"),
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
