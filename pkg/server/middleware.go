package server

import (
	"github.com/gofiber/fiber/v2"
	midlog "github.com/gofiber/fiber/v2/middleware/logger"
)

// SetServerMiddleware
// TODO: maybe we'll need to configure middlware via config in the future
func SetServerMiddleware(app *fiber.App) {

	app.Use(
		recodeMiddleware(),
	).Name("middleware")

}

// use logger recode HTTP info
func recodeMiddleware() fiber.Handler {
	return midlog.New(midlog.ConfigDefault)
}
