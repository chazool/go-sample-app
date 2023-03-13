package handler

import (
	responsebuilder "github.com/chazool/go-sample-app/app/routes/responseBuilder"
	"github.com/chazool/go-sample-app/common/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func Lives(ctx *fiber.Ctx) (err error) {
	var (
		requestID     = GetRequestID(ctx)
		parentSpan, _ = SetTracing(ctx, LivezMethod, requestID)
	)

	defer utils.CloseSpan(parentSpan)

	responseBuilder := responsebuilder.APIResponse{
		ctx:ctx,
	}

	return nil
}
