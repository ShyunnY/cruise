package logic

import (
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/gofiber/fiber/v2"
)

func ListOperationsService(svcCtx ServiceCtx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		svcN := ServiceNameParam{}
		err := ctx.ParamsParser(&svcN)
		if err != nil {
			return err
		}
		if svcN.Empty() {
			return fiber.NewError(fiber.StatusBadRequest, "service cannot empty,must provide no-empty value")
		}

		// list all operation for target service
		operations := svcCtx.Store.ListOperations(svcN.ServiceName)
		if operations != nil && len(operations) > 0 {
			return ctx.Status(fiber.StatusOK).JSON(operations)
		}

		// try to query reader operations
		queryOperations, err := svcCtx.Reader.QueryOperations(ctx.Context(), reader.QueryOperationsRequest{
			Service: svcN.ServiceName,
		})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(queryOperations)
	}
}
