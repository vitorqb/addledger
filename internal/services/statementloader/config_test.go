package statementloader_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/services/statementloader"
	"github.com/vitorqb/addledger/internal/testutils"
)

func TestLoadCsvStatementLoaderConfig(t *testing.T) {
	csvFile := testutils.TestDataPath(t, "statement.csv")
	minPresetFile := testutils.TestDataPath(t, "csv_preset_min.json")
	fullPresetFile := testutils.TestDataPath(t, "csv_preset_full.json")

	t.Run("No file", func(t *testing.T) {
		config, err := LoadConfig("", "")
		assert.Equal(t, Config{}, config)
		assert.NoError(t, err)
	})

	t.Run("No preset", func(t *testing.T) {
		config, err := LoadConfig(csvFile, "")
		assert.Equal(t, Config{}, config)
		assert.ErrorContains(t, err, "missing preset")
	})

	t.Run("Preset not found", func(t *testing.T) {
		config, err := LoadConfig(csvFile, "foo")
		assert.Equal(t, Config{}, config)
		assert.ErrorContains(t, err, "failed to open preset file")
	})

	t.Run("Preset as file name loads from config dir", func(t *testing.T) {
		t.Setenv("HOME", "/home/foo")
		_, err := LoadConfig(csvFile, "foo")
		assert.ErrorContains(t, err, "/home/foo/.config/addledger/presets/foo.json")
	})

	t.Run("Minimal preset found", func(t *testing.T) {
		config, err := LoadConfig(csvFile, minPresetFile)
		assert.NoError(t, err)
		assert.Equal(t, Config{
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
		config, err := LoadConfig(csvFile, fullPresetFile)
		assert.NoError(t, err)
		assert.Equal(t, Config{
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
