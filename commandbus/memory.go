package pkgCommandBus

// import (
// 	"context"
// 	"fmt"
// 	pkgApplication "go_utils/application"
// 	pkgDomain "go_utils/domain"
// 	pkgError "go_utils/errors"
// 	pkgMessageBus "go_utils/messagebus"

// 	"github.com/sirupsen/logrus"
// )

// // CommandBusMemory allows to subscribe/dispatch commands
// // Subscribing to the same command twice will unsubscribe previous handler
// // command handler should be one to one
// type CommandBusMemory struct {
// 	logger *logrus.Logger
// 	messageBus pkgMessageBus.MessageBus
// }

// func NewMemory(maxConcurrentCalls int) CommandBusInterface {
// 	return &CommandBusMemory{
// 		logger: logrus.New(),
// 		messageBus: pkgMessageBus.New(maxConcurrentCalls),
// 	}
// }

// func (bus *CommandBusMemory) Publish(ctx context.Context, command pkgDomain.Command) error{
// 	out := make(chan error, 1)
// 	defer close(out)

// 	bus.logger.Info(fmt.Sprintf("[CommandBus] Publish: %s %+v", command.GetName(), command))
// 	bus.messageBus.Publish(command.GetName(), ctx, command, out)

// 	ctxDoneCh := ctx.Done()

// 	select {
// 	case <-ctxDoneCh:
// 		return pkgError.Wrap(fmt.Errorf("%w: %s", pkgApplication.ErrTimeout, ctx.Err()))
// 	case err := <-out:
// 		if err != nil {
// 			return pkgError.Wrap(fmt.Errorf("create client failed: %w", err))
// 		}
// 		return nil
// 	}
// }

// func (bus *CommandBusMemory) Subscribe(ctx context.Context, commandName string, fn CommandHandler) error{
// 	bus.logger.Info(fmt.Sprintf("[CommandBus] Subscribe: %s", commandName))

// 	// unsubscribe all other handlers
// 	bus.messageBus.Close(commandName)

// 	return bus.messageBus.Subscribe(commandName, func(ctx context.Context, command pkgDomain.Command, out chan<- error) {
// 		out <- fn(ctx, command)
// 	})
// }

// func (bus *CommandBusMemory) Unsubscribe(ctx context.Context, commandName string) error {
// 	bus.logger.Info(nil, "[CommandBus] Unsubscribe: %s", commandName)
// 	bus.messageBus.Close(commandName)

// 	return nil
// }


