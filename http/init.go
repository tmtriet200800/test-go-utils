package pkgHttp

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// Adapter is http server app adapter
type Adapter struct {
	httpServer *fiber.App
}

// NewAdapter provides new primary HTTP adapter
func NewAdapter(httpServer *fiber.App) *Adapter {
	return &Adapter{
		httpServer: httpServer,
	}
}

// Start start http application adapter
func (adapter *Adapter) Start(ctx context.Context) error {
	return adapter.httpServer.Listen(":3000")
}

// Stop stops http application adapter
func (adapter *Adapter) Stop(ctx context.Context) error {
	return adapter.httpServer.Shutdown()
}
