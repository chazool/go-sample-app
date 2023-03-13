package utils

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
)

// CreateChildSpan create a child span based on the parent span context
func CreateChildSpan(ctx context.Context, spanName string) (trace.Span, context.Context) {
	ctx, span := otel.GetTracerProvider().Tracer("componet-main").Start(ctx, spanName)
	return span, ctx
}

// closeSpan is used to close a span
func CloseSpan(span trace.Span) {
	span.End()
}
