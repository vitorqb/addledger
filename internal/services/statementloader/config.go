package statementloader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vitorqb/addledger/internal/utils"
)

func expandUserHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(os.Getenv("HOME"), path[2:])
	}
	return path
}

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

type ConfigLoader struct {
	PresetsDir string
}

func (cf *ConfigLoader) Load(file, preset string) (Config, error) {
	if file == "" {
		return Config{}, nil
	}
	if preset == "" {
		defaultPresetFile := filepath.Join(cf.PresetsDir, "default.json")
		if _, err := os.Stat(defaultPresetFile); err != nil {
			return Config{}, fmt.Errorf("missing preset (and no default defined)")
		}
		preset = defaultPresetFile
	}
	if !utils.LooksLikePath(preset) {
		preset = filepath.Join(cf.PresetsDir, preset)
	}
	if filepath.Ext(preset) == "" {
		preset += ".json"
	}
	preset = expandUserHome(preset)
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
	config.File = expandUserHome(file)
	return config, nil
}

func LoadConfig(file, preset string) (Config, error) {
	presetsDir := filepath.Join(os.Getenv("HOME"), ".config/addledger/presets")
	loader := ConfigLoader{presetsDir}
	return loader.Load(file, preset)
}
