package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/config"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

func TestJournalFileFinder(t *testing.T) {
	executable := tu.TestDataPath(t, "fake_hledger.sh")
	loader := NewLoader()
	result, err := loader.JournalFile(executable)
	assert.Nil(t, err)
	assert.Equal(t, result, "/path/to/hledger.journal")
}
