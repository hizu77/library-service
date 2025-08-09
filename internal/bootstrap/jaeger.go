package bootstrap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

const FlushTimeout = time.Second * 5

func InitTracer(ctx context.Context, url string) (*trace.TracerProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	exp, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(url),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "init exporter")
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("library-service"),
			),
		),
	)

	otel.SetTracerProvider(tracerProvider)

	return tracerProvider, nil
}
