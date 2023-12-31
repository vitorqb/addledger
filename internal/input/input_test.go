package input_test

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/vitorqb/addledger/internal/finance"
	. "github.com/vitorqb/addledger/internal/input"
	"github.com/vitorqb/addledger/internal/journal"
	tu "github.com/vitorqb/addledger/internal/testutils"
)

func TestJournalEntryInput(t *testing.T) {

	type context struct {
		onChangeCalled    bool
		onChangeCallCount int
		input             *JournalEntryInput
	}

	type test struct {
		name string
		run  func(*testing.T, *context)
	}

	tests := []test{
		{
			name: "Date",
			run: func(t *testing.T, c *context) {
				_, found := c.input.GetDate()
				assert.False(t, found)
				c.input.SetDate(tu.Date1(t))
				date, found := c.input.GetDate()
				assert.True(t, found)
				assert.Equal(t, date, tu.Date1(t))
				assert.Equal(t, 1, c.onChangeCallCount)
				c.input.ClearDate()
				_, found = c.input.GetDate()
				assert.False(t, found)
				assert.Equal(t, 2, c.onChangeCallCount)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := NewJournalEntryInput()
			ctx := &context{false, 0, input}
			input.AddOnChangeHook(func() { ctx.onChangeCalled = true })
			input.AddOnChangeHook(func() { ctx.onChangeCallCount++ })
			tc.run(t, ctx)
		})
	}
}

func TestTextToAmmount(t *testing.T) {
	type testcase struct {
		text     string
		ammount  finance.Ammount
		errorMsg string
	}
	var testcases = []testcase{
		{
			text: "EUR 12.20",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text: "EUR 99999.99999",
			ammount: finance.Ammount{
				Commodity: "EUR",
				Quantity:  decimal.NewFromFloat(99999.99999),
			},
		},
		{
			text: "12.20",
			ammount: finance.Ammount{
				Commodity: "",
				Quantity:  decimal.New(1220, -2),
			},
		},
		{
			text:     "12,20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12 12",
			errorMsg: "invalid format",
		},
		{
			text:     "12 FOO",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR  12.20",
			errorMsg: "invalid format",
		},
		{
			text:     "EUR 12.20 ",
			errorMsg: "invalid format",
		},
		{
			text:     " EUR 12.20 ",
			errorMsg: "invalid format",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.text, func(t *testing.T) {
			result, err := TextToAmmount(tc.text)
			if tc.errorMsg == "" {
				assert.Nil(t, err)
				assert.Equal(t, tc.ammount, result)
			} else {
				assert.ErrorContains(t, err, tc.errorMsg)
			}
		})
	}
}

func TestTagTagToText(t *testing.T) {
	tag := journal.Tag{
		Name:  "foo",
		Value: "bar",
	}
	expected := "foo:bar"
	actual := TagToText(tag)
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestTagTextToTag__Good(t *testing.T) {
	type testCase struct {
		input    string
		expected journal.Tag
	}
	for _, tc := range []testCase{
		{
			input:    "foo:bar",
			expected: journal.Tag{Name: "foo", Value: "bar"},
		},
		{
			input:    " foo:bar ",
			expected: journal.Tag{Name: "foo", Value: "bar"},
		},
		{
			input:    "foo-bar:baz",
			expected: journal.Tag{Name: "foo-bar", Value: "baz"},
		},
		{
			input:    "foo_bar:baz",
			expected: journal.Tag{Name: "foo_bar", Value: "baz"},
		},
		{
			input:    "foo_bar:baz_123",
			expected: journal.Tag{Name: "foo_bar", Value: "baz_123"},
		},
		{
			input:    "foo_bar:baz-123",
			expected: journal.Tag{Name: "foo_bar", Value: "baz-123"},
		},
		{
			input:    "foo_bar:baz-123_abc",
			expected: journal.Tag{Name: "foo_bar", Value: "baz-123_abc"},
		},
	} {
		tag, err := TextToTag(tc.input)
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
		if tag != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, tag)
		}
	}
}

func TestTagTextToTag__Bad(t *testing.T) {
	for _, input := range []string{
		"foo",
		"foo:",
		"foo:bar:baz",
		"",
		"foo bar:baz",
		"foo:bar baz",
		"some word",
		"some word:foo",
	} {
		_, err := TextToTag(input)
		if err == nil {
			t.Errorf(fmt.Sprintf("Expected error, got none: %s", input))
		}

	}
}
