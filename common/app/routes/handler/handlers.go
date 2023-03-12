package handler

import (
	"context"

	"github.com/chazool/go-sample-app/common/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberOtel "github.com/psmarcin/fiber-opentelemetry/pkg/fiber-otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func SetTracing(ctx *fiber.Ctx, spanFuncName, requestId string) (trace.Span, context.Context) {

	fCtx := fiberOtel.FromCtx(ctx)
	parentSpan := trace.SpanFromContext(fCtx)
	parentSpan.SetAttributes(attribute.String("Request Id", requestId))

	span := &config.OpentelemetryParantCtx{
		ParentCtx:  fCtx,
		ParentSpan: parentSpan,
	}
	span.SetOpentelementryParentCtx()
 
	return parentSpan, fCtx
}

func UpdateSpanStatus(parentSpan *trace.Span, statusCode int) {
	span := *parentSpan
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(statusCode))
}

func GetRequestID(ctx *fiber.Ctx) string {
	return ctx.Locals(requestid.ConfigDefault.ContextKey).(string)
}
