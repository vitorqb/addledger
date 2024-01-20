package dateguesser_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/dateguesser"
	"github.com/vitorqb/addledger/internal/finance"
	"github.com/vitorqb/addledger/internal/testutils"
	. "github.com/vitorqb/addledger/mocks/dateguesser"
)

func TestDateGuesser(t *testing.T) {
	type testcontext struct {
		clock   *MockIClock
		guesser *DateGuesser
	}
	type testcase struct {
		name string
		run  func(c *testcontext, t *testing.T)
	}
	var testcases = []testcase{
		{
			name: "No input suggests today",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t), guess)
			},
		},
		{
			name: "No input with statement entry suggests statement entry",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().AnyTimes().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("", finance.StatementEntry{Date: testutils.Date2(t)})
				assert.True(t, success)
				assert.Equal(t, testutils.Date2(t), guess)
			},
		},
		{
			name: "Valid user input (full date)",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("1993-11-23", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t), guess)
			},
		},
		{
			name: "Valid user input (day only)",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("23", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t), guess)
			},
		},
		{
			name: "Valid user input (day only, 1 digit month)",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().AnyTimes().Return(testutils.Date2(t))
				guess, success := c.guesser.Guess("01", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date2(t), guess)
			},
		},
		{
			name: "Valid user input (month and day)",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("11-23", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t), guess)
			},
		},
		{
			name: "Previous day",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("-1", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t).AddDate(0, 0, -1), guess)
			},
		},
		{
			name: "Two days ago",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				guess, success := c.guesser.Guess("-2", finance.StatementEntry{})
				assert.True(t, success)
				assert.Equal(t, testutils.Date1(t).AddDate(0, 0, -2), guess)
			},
		},
		{
			name: "Invalid user input",
			run: func(c *testcontext, t *testing.T) {
				c.clock.EXPECT().Now().Return(testutils.Date1(t))
				_, success := c.guesser.Guess("41", finance.StatementEntry{})
				assert.False(t, success)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := new(testcontext)
			c.clock = NewMockIClock(ctrl)
			c.guesser = &DateGuesser{c.clock}
			tc.run(c, t)
		})
	}
}
