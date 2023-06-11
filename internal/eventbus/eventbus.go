package eventbus

//go:generate $MOCKGEN --source=eventbus.go --destination=../../mocks/eventbus/eventbus_mock.go

// Event represents an event
type Event struct {
	Topic string
	Data  interface{}
}

// Subscription represents a subscriptions to a topic
type Subscription struct {
	Topic   string
	Handler func(e Event)
}

// EventBus routes events
type IEventBus interface {
	Send(e Event) error
	Subscribe(s Subscription) error
}

type EventBus struct {
	subscriptions []Subscription
}

// New creates a new EventBus
func New() *EventBus {
	return &EventBus{}
}

func (eb *EventBus) Send(e Event) error {
	for _, subscription := range eb.subscriptions {
		if subscription.Topic == e.Topic {
			subscription.Handler(e)
		}
	}
	return nil
}

func (eb *EventBus) Subscribe(s Subscription) error {
	eb.subscriptions = append(eb.subscriptions, s)
	return nil
}
