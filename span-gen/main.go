package main

import (
    "context"
    "errors"
    "fmt"
    "time"

    "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
    ctx := context.Background()
	otelShutdown, err := setupOTelSDK(ctx)
	if err != nil {
		return
	}

	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)
    tracer := tracerProvider.Tracer("github.com/danslimmon/clerk_go")

    ticker := time.Tick(1*time.Second)
    for {
        <-ticker
        _, span := tracer.Start(ctx, "blah-span")
        time.Sleep(20*time.Millisecond)
        span.End()
        fmt.Println("sent span")
    }
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider() (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(
        context.Background(),
        otlptracegrpc.WithInsecure(),
        otlptracegrpc.WithEndpoint("otel-collector:4317"),
    )
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}
