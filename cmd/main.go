package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})

	// search trace read group
	// compatibility jaeger Read interface
	apiG := app.Group("/read")
	apiG.Get("/search")
	apiG.Get("/query/:traceID")
	apiG.Get("/service/:service")
	apiG.Get("/operation/:operation")

}
