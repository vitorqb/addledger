package statementloader

import (
	"encoding/json"
	"fmt"
	"github.com/vitorqb/addledger/internal/utils"
	"os"
	"path/filepath"
)

type Config struct {
	// File to load statement from.
	File string
	// Separator to use.
	Separator string `json:"separator"`
	// Default account to use for all entries.
	Account string `json:"account"`
	// Default commodity to use for all entries.
	Commodity string `json:"commodity"`
	// SortBy defines a stratgy for sorting. As of now either empty (no sorting)
	// or date are supported.
	SortBy string `json:"sortBy"`
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

func LoadConfig(file, preset string) (Config, error) {
	if file == "" {
		return Config{}, nil
	}
	if preset == "" {
		return Config{}, fmt.Errorf("missing preset")
	}
	if !utils.LooksLikePath(preset) {
		preset = fmt.Sprintf("%s/.config/addledger/presets/%s", os.Getenv("HOME"), preset)
	}
	if filepath.Ext(preset) == "" {
		preset += ".json"
	}
	presetBytes, err := os.ReadFile(preset)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open preset file %s: %w", preset, err)
	}
	var config Config
	config.AccountFieldIndex = -1
	config.AmmountFieldIndex = -1
	config.DateFieldIndex = -1
	config.DescriptionFieldIndex = -1
	config.DateFormat = "02/01/2006"
	err = json.Unmarshal(presetBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal preset file: %w", err)
	}
	config.File = file
	return config, nil
}
