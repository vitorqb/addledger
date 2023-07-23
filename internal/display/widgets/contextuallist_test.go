package widgets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/display/widgets"
	"github.com/vitorqb/addledger/internal/listaction"
)

func TestContextualList(t *testing.T) {
	type testcontext struct {
		contextualList *ContextualList
		selected       string
		input          string
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Fills list with items",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, 3, c.contextualList.GetItemCount())
				text, _ := c.contextualList.GetItemText(0)
				assert.Equal(t, "THREE", text)
			},
		},
		{
			name: "Stores selected item",
			run: func(t *testing.T, c *testcontext) {
				c.contextualList.SetCurrentItem(2)
				assert.Equal(t, "ONE", c.selected)
			},
		},
		{
			name: "Handle list actions",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, c.selected, "THREE")
				c.contextualList.HandleAction(listaction.NEXT)
				assert.Equal(t, c.selected, "TWO")
			},
		},
		{
			name: "Refresh",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, c.contextualList.GetItemCount(), 3)
				c.contextualList.Refresh()
				assert.Equal(t, c.contextualList.GetItemCount(), 2)
				c.input = "THREE"
				c.contextualList.Refresh()
				assert.Equal(t, c.contextualList.GetItemCount(), 1)
				assert.Equal(t, c.selected, "THREE")
			},
		},
		{
			name: "Refresh preserves order from getItemsFunc",
			run: func(t *testing.T, c *testcontext) {
				// Writes something
				c.input = "THREE"
				c.contextualList.Refresh()

				// Resets
				c.input = ""
				c.contextualList.Refresh()
				assert.Equal(t, c.contextualList.GetItemCount(), 3)
				assert.Equal(t, "THREE", c.selected)
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			c := new(testcontext)
			c.input = "T"
			c.contextualList = NewContextualList(
				func() []string {
					return []string{"THREE", "TWO", "ONE"}
				},
				func(s string) {
					c.selected = s
				},
				func() string {
					return c.input
				},
			)
			testcase.run(t, c)
		})
	}
}
