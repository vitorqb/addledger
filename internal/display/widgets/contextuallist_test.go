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
		options        *ContextualListOptions
	}
	type testcase struct {
		name            string
		run             func(t *testing.T, c *testcontext)
		setupOptions    func(o *ContextualListOptions)
		buildErrorMatch string
	}
	testcases := []testcase{
		{
			name: "Fails to build if missing GetItemsFunc",
			setupOptions: func(o *ContextualListOptions) {
				o.GetInputFunc = nil
			},
			buildErrorMatch: "missing GetInputFunc",
		},
		{
			name: "Fails to build if missing GetInputFunc",
			setupOptions: func(o *ContextualListOptions) {
				o.GetInputFunc = nil
			},
			buildErrorMatch: "missing GetInputFunc",
		},
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
			name: "Refresh with empty list sets empty string",
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, c.contextualList.GetItemCount(), 3)
				c.contextualList.Refresh()
				assert.Equal(t, c.contextualList.GetItemCount(), 2)
				c.input = "adjsalkkjsd"
				c.contextualList.Refresh()
				assert.Equal(t, c.contextualList.GetItemCount(), 0)
				assert.Equal(t, "", c.selected)
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
		{
			name: "Default is printed first.",
			setupOptions: func(o *ContextualListOptions) {
				o.GetDefaultFunc = func() (string, bool) {
					return "FOO", true
				}
			},
			run: func(t *testing.T, c *testcontext) {
				assert.Equal(t, 4, c.contextualList.GetItemCount())
				assert.Equal(t, "FOO", c.selected)
				firstItem, _ := c.contextualList.GetItemText(0)
				assert.Equal(t, "FOO", firstItem)

				// Refreshes and keeps the same
				c.input = ""
				c.contextualList.Refresh()
				assert.Equal(t, 4, c.contextualList.GetItemCount())
				assert.Equal(t, "FOO", c.selected)
				firstItem, _ = c.contextualList.GetItemText(0)
				assert.Equal(t, "FOO", firstItem)
			},
		},
		{
			name: "Default is printed after writting something.",
			setupOptions: func(o *ContextualListOptions) {
				o.GetDefaultFunc = func() (string, bool) {
					return "FOO", true
				}
			},
			run: func(t *testing.T, c *testcontext) {
				c.input = "ONE"
				c.contextualList.Refresh()
				c.input = ""
				c.contextualList.Refresh()

				// Assert default is at pos 0
				assert.Equal(t, 4, c.contextualList.GetItemCount())
				assert.Equal(t, "FOO", c.selected)
				firstItem, _ := c.contextualList.GetItemText(0)
				assert.Equal(t, "FOO", firstItem)
			},
		},
		{
			name: "Sets selected item on Refresh",
			run: func(t *testing.T, c *testcontext) {
				c.selected = ""
				c.contextualList.Refresh()
				assert.Equal(t, "TWO", c.selected)
				// Note: run again because of cache
				c.selected = ""
				c.contextualList.Refresh()
				assert.Equal(t, "TWO", c.selected)
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			var err error
			c := new(testcontext)
			c.input = "T"
			c.options = &ContextualListOptions{
				GetItemsFunc: func() []string {
					return []string{"THREE", "TWO", "ONE"}
				},
				SetSelectedFunc: func(s string) {
					c.selected = s
				},
				GetInputFunc: func() string {
					return c.input
				},
			}
			if testcase.setupOptions != nil {
				testcase.setupOptions(c.options)
			}
			c.contextualList, err = NewContextualList(*c.options)
			if testcase.buildErrorMatch != "" {
				assert.ErrorContains(t, err, testcase.buildErrorMatch)
			} else {
				testcase.run(t, c)
			}
		})
	}
}
