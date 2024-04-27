package statementloader_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	. "github.com/vitorqb/addledger/internal/services/statementloader"
	statemod "github.com/vitorqb/addledger/internal/state"
	"github.com/vitorqb/addledger/internal/statementreader"
	"github.com/vitorqb/addledger/internal/testutils"
	statementreader_mock "github.com/vitorqb/addledger/mocks/statementreader"
)

func TestStatementLoaderSvc(t *testing.T) {
	statement := testutils.TestDataPath(t, "statement.csv")
	type testcontext struct {
		state   *statemod.State
		reader  *statementreader_mock.MockIStatementReader
		service *Service
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Fail to read file",
			run: func(t *testing.T, c *testcontext) {
				config := Config{File: "not-a-file"}
				err := c.service.Load(config)
				assert.ErrorContains(t, err, "failed to open file")
			},
		},
		{
			name: "Fail to load statement",
			run: func(t *testing.T, c *testcontext) {
				config := Config{File: statement}
				c.reader.EXPECT().Read(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				err := c.service.Load(config)
				assert.ErrorContains(t, err, "failed to load statement")
			},
		},
		{
			name: "Success",
			run: func(t *testing.T, c *testcontext) {
				entries := []finance.StatementEntry{{Account: "ACC"}}
				config := Config{File: statement}
				c.reader.EXPECT().Read(gomock.Any(), gomock.Any()).Return(entries, nil)
				err := c.service.Load(config)
				assert.Nil(t, err)
				assert.Equal(t, entries, c.state.GetStatementEntries())
			},
		},
		{
			name: "LoadFromFiles Success",
			run: func(t *testing.T, c *testcontext) {
				entries := []finance.StatementEntry{{Account: "ACC"}}
				presetFile := testutils.TestDataPath(t, "csv_preset_full.json")
				c.reader.EXPECT().Read(gomock.Any(), gomock.Any()).Return(entries, nil)
				err := c.service.LoadFromFiles(statement, presetFile)
				assert.Nil(t, err)
				assert.Equal(t, entries, c.state.GetStatementEntries())
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.state = statemod.InitialState()
			c.reader = statementreader_mock.NewMockIStatementReader(ctrl)
			c.service = New(c.state, c.reader)
			tc.run(t, c)
		})
	}
}

func TestParseStatementLoaderConfig(t *testing.T) {
	type testcase struct {
		name            string
		config          Config
		expectedOptions []statementreader.Option
		expectedError   string
	}
	testcases := []testcase{
		{
			name: "empty",
			config: Config{
				DateFieldIndex:        -1,
				DescriptionFieldIndex: -1,
				AccountFieldIndex:     -1,
				AmmountFieldIndex:     -1,
			},
			expectedOptions: []statementreader.Option{
				statementreader.WithLoaderMapping([]statementreader.CSVColumnMapping{}),
			},
		},
		{
			name: "full",
			config: Config{
				Separator:             ";",
				Account:               "acc",
				Commodity:             "com",
				DateFieldIndex:        0,
				DateFormat:            "01/02/2006",
				DescriptionFieldIndex: 1,
				AccountFieldIndex:     2,
				AmmountFieldIndex:     3,
				SortBy:                "date",
			},
			expectedOptions: []statementreader.Option{
				statementreader.WithSeparator(';'),
				statementreader.WithAccountName("acc"),
				statementreader.WithDefaultCommodity("com"),
				statementreader.WithSortStrategy(statementreader.SortByDate{}),
				statementreader.WithLoaderMapping([]statementreader.CSVColumnMapping{
					{Column: 0, Importer: statementreader.DateImporter{Format: "01/02/2006"}},
					{Column: 1, Importer: statementreader.DescriptionImporter{}},
					{Column: 2, Importer: statementreader.AccountImporter{}},
					{Column: 3, Importer: statementreader.AmmountImporter{}},
				}),
			},
		},
		{
			name: "invalid sortBy",
			config: Config{
				SortBy: "invalid",
			},
			expectedError: "invalid SortBy: invalid",
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			actualConfig := statementreader.Config{}
			expectedConfig := statementreader.Config{}
			options, err := ParseConfig(testcase.config)
			if expError := testcase.expectedError; expError != "" {
				assert.ErrorContains(t, err, expError)
				return
			}
			assert.Nil(t, err)
			for _, option := range options {
				option(&actualConfig)
			}
			for _, option := range testcase.expectedOptions {
				option(&expectedConfig)
			}
			assert.Equal(t, expectedConfig, actualConfig)
		})
	}
}
