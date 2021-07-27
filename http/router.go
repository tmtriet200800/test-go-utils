package pkgHttp

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	healthRouter *fiber.App
	readinessRouter *fiber.App

	apiRouter *fiber.App
}

// New provides new router
func NewRouter(healthRouter *fiber.App, readinessRouter *fiber.App, apiRouter *fiber.App) *Router {
	return &Router{
		healthRouter: healthRouter,
		readinessRouter: readinessRouter,
		apiRouter: apiRouter,
	}
}


func (router *Router) AddApp(app *fiber.App) *fiber.App{
	app.Mount("/health", router.healthRouter)
	app.Mount("/readiness", router.readinessRouter)
	app.Mount("/api", router.apiRouter)

	return app
}