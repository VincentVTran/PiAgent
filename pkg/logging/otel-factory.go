package logging

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// InitTelemetry configures OpenTelemetry logging to emit logs to both
// an OTLP collector and the local console.
func InitTelemetry(ctx context.Context) (func(context.Context) error, error) {
	// 1) Set up OTLP Log exporter to send logs to collector
	logExp, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint("192.168.2.5:4318"),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// 2) Set up Console (stdout) Log exporter for local output
	consoleExp, err := stdoutlog.New(
		stdoutlog.WithPrettyPrint(),
	)
	if err != nil {
		return nil, err
	}

	// 3) Create LoggerProvider with two batch processors: one for OTLP, one for console
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExp)),
		log.WithProcessor(log.NewBatchProcessor(consoleExp)),
	)
	global.SetLoggerProvider(lp)

	// 4) (Optional) Set up a TracerProvider for tracing
	tp := sdktrace.NewTracerProvider(
	// TODO: add resource, sampler, and OTLP trace exporter similarly
	)
	otel.SetTracerProvider(tp)

	// 5) Return shutdown func
	return func(ctx context.Context) error {
		_ = lp.Shutdown(ctx)
		return tp.Shutdown(ctx)
	}, nil
}
