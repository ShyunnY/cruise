package logic

import (
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/gofiber/fiber/v2"
)

func ListTraceSvcService(svcCtx ServiceCtx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		// list all service
		services := svcCtx.Store.ListServices()
		if services != nil && len(services) > 0 {

			return ctx.Status(fiber.StatusOK).JSON(services)
		}

		// try to query reader services
		queryServices, err := svcCtx.Reader.QueryServices(ctx.Context(), reader.QueryServicesRequest{})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(queryServices)
	}
}
