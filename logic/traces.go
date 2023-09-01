package logic

import (
	"github.com/ShyunnY/cruise/pkg/reader"
	"github.com/gofiber/fiber/v2"
)

func SearchTracesService(svcCtx ServiceCtx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		// 1.条件查询 稍后再说

		return nil
	}
}

// QueryTracesService
// query trace info by trace id
func QueryTracesService(svcCtx ServiceCtx) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		// query Trace
		tid := TraceIDParam{}
		if err := ctx.ParamsParser(&tid); err != nil {
			return err
		}

		// check if empty
		if tid.Empty() {
			return fiber.NewError(fiber.StatusBadRequest, "trace id cannot empty,must provide no-empty value")
		}

		// check store
		trace := svcCtx.Store.GetTrace(tid.TraceID)
		if trace != nil {
			return ctx.Status(fiber.StatusOK).JSON(trace.GetResourceSpans())
		}

		// try to query reader trace
		queryTrace, err := svcCtx.Reader.QueryTrace(ctx.Context(), reader.QueryTraceRequest{TraceID: tid.TraceID})
		if err != nil {
			return err
		}

		return ctx.Status(fiber.StatusOK).JSON(queryTrace)
	}
}
