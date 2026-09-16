package telemetry

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const TracerName = "glovebox-go-service"

// InitTelemetry sets up structured slog logging and OpenTelemetry tracing.
func InitTelemetry(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, *slog.Logger, error) {
	// Structured logger setup
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// OpenTelemetry tracer provider setup
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, logger, nil
}

// TraceMiddleware creates an OpenTelemetry span and logs HTTP request metrics via slog.
func TraceMiddleware(next http.Handler) http.Handler {
	tracer := otel.GetTracerProvider().Tracer(TracerName)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.target", r.URL.Path),
			),
		)
		defer span.End()

		start := time.Now()
		ww := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r.WithContext(ctx))

		duration := time.Since(start)
		span.SetAttributes(
			attribute.Int("http.status_code", ww.statusCode),
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
		)

		slog.InfoContext(ctx, "HTTP request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.statusCode,
			"duration", duration.String(),
			"trace_id", span.SpanContext().TraceID().String(),
		)
	})
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
