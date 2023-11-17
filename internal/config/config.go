package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vitorqb/addledger/internal/utils"
)

// PrinterConfig represents the value for configuring a printer.Printer.
type PrinterConfig struct {
	NumLineBreaksBefore int // Number of empty lines to print before a transaction.
	NumLineBreaksAfter  int // Number of empty lines to print after a transaction.
}

// CSVStatementLoaderConfig represents the value for configuring a statementloader.CSVLoader.
type CSVStatementLoaderConfig struct {
	// File to load.
	File string
	// Separator to use.
	Separator string `json:"separator"`
	// Default account to use for all entries.
	Account string `json:"account"`
	// Default commodity to use for all entries.
	Commodity string `json:"commodity"`
	// Index of the date field in the CSV file.
	DateFieldIndex int `json:"dateFieldIndex"`
	// Date format to use for parsing the date field.
	DateFormat string `json:"dateFormat"`
	// Index of the account field in the CSV file.
	AccountFieldIndex int `json:"accountFieldIndex"`
	// Index of the description field in the CSV file.
	DescriptionFieldIndex int `json:"descriptionFieldIndex"`
	// Index of the ammount field in the CSV file.
	AmmountFieldIndex int `json:"ammountFieldIndex"`
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
	// Configures a CSV statement loader
	CSVStatementLoaderConfig CSVStatementLoaderConfig
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

	// CSV Statement Loader config
	flagSet.String("csv-statement-file", "", "CSV file to load as a statement.")
	flagSet.String("csv-statement-preset", "", "Preset to use for CSV statement. If a simple filename is given, it will be searched in ~/.config/addledger/presets (with a .json extension).")
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

	// Loads the CSVStatementLoaderConfig
	csvStatementLoaderConfig, err := LoadCsvStatementLoaderConfig(
		viper.GetString("csv-statement-file"),
		viper.GetString("csv-statement-preset"),
	)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to load csv statement config: %w", err)
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
		CSVStatementLoaderConfig: csvStatementLoaderConfig,
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

func LoadCsvStatementLoaderConfig(file, preset string) (CSVStatementLoaderConfig, error) {
	if file == "" {
		return CSVStatementLoaderConfig{}, nil
	}
	if preset == "" {
		return CSVStatementLoaderConfig{}, fmt.Errorf("missing preset")
	}
	if !utils.LooksLikePath(preset) {
		preset = fmt.Sprintf("%s/.config/addledger/presets/%s", os.Getenv("HOME"), preset)
	}
	if filepath.Ext(preset) == "" {
		preset += ".json"
	}
	presetBytes, err := os.ReadFile(preset)
	if err != nil {
		return CSVStatementLoaderConfig{}, fmt.Errorf("failed to open preset file %s: %w", preset, err)
	}
	var config CSVStatementLoaderConfig
	config.AccountFieldIndex = -1
	config.AmmountFieldIndex = -1
	config.DateFieldIndex = -1
	config.DescriptionFieldIndex = -1
	config.DateFormat = "02/01/2006"
	err = json.Unmarshal(presetBytes, &config)
	if err != nil {
		return CSVStatementLoaderConfig{}, fmt.Errorf("failed to unmarshal preset file: %w", err)
	}
	config.File = file
	return config, nil
}
