package resources

// EventHandler is a function that processes a raw event payload.
type EventHandler func(subject string, payload []byte)

// EventSubscriber defines the contract for subscribing to domain events from a message broker.
type EventSubscriber interface {
	Subscribe(subject string, handler EventHandler) error
	Close()
}
