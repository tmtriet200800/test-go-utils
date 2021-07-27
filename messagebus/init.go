package pkgMessageBus

import (
	"context"
)

type Message interface {}

// MessageHandler function
type MessageHandler func(ctx context.Context, message *Message) error

// EventBus interface
// event bus is different to command bus by allowing multiple handlers to the same topic
type MessageBusInterface interface {
	Publish(ctx context.Context, message *Message) error
	Subscribe(ctx context.Context, eventType string, fn MessageHandler) error
	// Unsubscribe(ctx context.Context, eventType string, fn EventHandler) error

	// PublishAndAcknowledge blocks and returns grouped error after all handlers are executed
	// PublishAndAcknowledge(parentCtx context.Context, event *domain.Event) error
}
