package pkgCommandBus

import (
	"context"

	pkgDomain "github.com/tmtriet200800/test-go-utils/domain"
)

// CommandHandler function
type CommandHandler func(ctx context.Context, command *pkgDomain.Command) error

// CommandBus allows to subscribe/dispatch commands
// Subscribing to the same command twice will unsubscribe previous handler
// command handler should be one to one
type CommandBusInterface interface {
	Publish(ctx context.Context, command *pkgDomain.Command) error
	Subscribe(ctx context.Context, commandName string, fn CommandHandler) error
		// Unsubscribe(ctx context.Context, eventType string, fn EventHandler) error

	// PublishAndAcknowledge blocks and returns grouped error after all handlers are executed
	// PublishAndAcknowledge(parentCtx context.Context, event *domain.Event) error
}
