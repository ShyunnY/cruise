package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// search trace api group
	apiG := app.Group("/api")
	apiG.Get("/search")
	apiG.Get("/query/:traceID")
	apiG.Get("/service/:service")
	apiG.Get("/operation/:operation")

}
