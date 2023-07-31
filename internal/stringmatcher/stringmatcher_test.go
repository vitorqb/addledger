package stringmatcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/stringmatcher"
)

func TestStringMatcher(t *testing.T) {
	{
		type testcontext struct{}
		type testcase struct {
			name string
			run  func(t *testing.T, c *testcontext)
		}
		var testcases = []testcase{
			{
				name: "Returns a proper distance",
				run: func(t *testing.T, c *testcontext) {
					matcher, err := New(&Options{})
					assert.Nil(t, err)
					result := matcher.Distance("foo", "foa")
					assert.Equal(t, 1, result)
					result = matcher.Distance("foa", "foo")
					assert.Equal(t, 1, result)
				},
			},
			{
				name: "Uses cache",
				run: func(t *testing.T, c *testcontext) {
					cache := NewCache()
					cache.Set("foo", "foa", 99)
					cache.Set("foa", "foo", 88)
					matcher, err := New(&Options{Cache: cache})
					assert.Nil(t, err)
					result := matcher.Distance("foo", "foa")
					assert.Equal(t, 99, result)
					result = matcher.Distance("foa", "foo")
					assert.Equal(t, 88, result)
				},
			},
		}
		for _, tc := range testcases {
			t.Run(tc.name, func(t *testing.T) {
				c := new(testcontext)
				tc.run(t, c)
			})
		}
	}
}
