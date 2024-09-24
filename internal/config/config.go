package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// PrinterConfig represents the value for configuring a printer.Printer.
type PrinterConfig struct {
	NumLineBreaksBefore int // Number of empty lines to print before a transaction.
	NumLineBreaksAfter  int // Number of empty lines to print after a transaction.
}

// Config is the root configuration for the entire app.
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
	// Configures the transaction printer
	PrinterConfig PrinterConfig
	// A initial file to load as a statement.
	CSVStatementFile string
	// A preset to use for the CSV statement.
	CSVStatementPreset string
	// Default file to load CSV sttatements from (interactively)
	DefaultCSVStatementFile string
}

func SetupFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP("destfile", "d", "", "Destination file (where we will write). Defaults to the ledger file.")
	flagSet.String("hledger-executable", "hledger", "Executable to use for HLedger")
	flagSet.String("ledger-file", "", "Ledger File to pass to HLedger commands. If empty let ledger executable find it.")
	flagSet.String("logfile", "", "File where to send log output. Empty for stderr.")
	flagSet.String("loglevel", "WARN", "Level of logger. Defaults to warning.")

	// Printer config
	flagSet.Int("printer-line-break-before", 1, "Number of line breaks to print before a transaction.")
	flagSet.Int("printer-line-break-after", 1, "Number of line breaks to print after a transaction.")

	// Statement Loader config
	flagSet.String("csv-statement-file", "", "CSV file to load as a statement.")
	flagSet.String("csv-statement-preset", "", "Preset to use for CSV statement. If a simple filename is given, it will be searched in ~/.config/addledger/presets (with a .json extension).")

	// Statement Modal config
	flagSet.String("default-csv-statement-file", "", "Default file to load statements from using the interactive modal.")
}

func Load(flagSet *pflag.FlagSet, args []string, loader ILoader) (*Config, error) {
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
	flagSet.VisitAll(func(f *pflag.Flag) {
		if err != nil {
			return
		}
		err = viper.BindEnv(f.Name)
	})
	if err != nil {
		return &Config{}, fmt.Errorf("failed to bind env: %w", err)
	}

	// Unpack
	config := &Config{
		DestFile:          viper.GetString("destfile"),
		HLedgerExecutable: viper.GetString("hledger-executable"),
		LedgerFile:        viper.GetString("ledger-file"),
		LogFile:           viper.GetString("logfile"),
		LogLevel:          viper.GetString("loglevel"),
		PrinterConfig: PrinterConfig{
			NumLineBreaksBefore: viper.GetInt("printer-line-break-before"),
			NumLineBreaksAfter:  viper.GetInt("printer-line-break-after"),
		},
		CSVStatementFile:        viper.GetString("csv-statement-file"),
		CSVStatementPreset:      viper.GetString("csv-statement-preset"),
		DefaultCSVStatementFile: viper.GetString("default-csv-statement-file"),
	}

	// Load dynamic values
	if config.DestFile == "" {
		config.DestFile = config.LedgerFile
	}
	if config.DestFile == "" {
		config.DestFile, err = loader.JournalFile(config.HLedgerExecutable)
	}

	// Validate
	if config.DestFile == "" {
		return config, fmt.Errorf("missing destination file!")
	}

	return config, nil
}

func LoadFromCommandLine() (*Config, error) {
	loader := NewLoader()
	SetupFlags(pflag.CommandLine)
	return Load(pflag.CommandLine, os.Args, loader)
}
