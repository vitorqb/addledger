package eventbus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/vitorqb/addledger/internal/eventbus"
)

func TestEventBus(t *testing.T) {
	type testcontext struct {
		eventBus *EventBus
	}
	type testcase struct {
		name string
		run  func(t *testing.T, c *testcontext)
	}
	testcases := []testcase{
		{
			name: "Subscribe and receives event",
			run: func(t *testing.T, c *testcontext) {
				var events []Event
				err := c.eventBus.Subscribe(Subscription{
					Topic: "foo.bar",
					Handler: func(e Event) {
						events = append(events, e)
					},
				})
				assert.Nil(t, err)
				event := Event{"foo.bar", 100}
				err = c.eventBus.Send(event)
				assert.Nil(t, err)
				assert.Equal(t, []Event{event}, events)
			},
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			c := new(testcontext)
			c.eventBus = New()
			testcase.run(t, c)
		})
	}
}
