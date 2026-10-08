package tracer

import (
	"context"
	"fmt"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.30.0"
	"time"
)

var AMPTracerProvider *trace.TracerProvider

func InitTracerOTEL() error {
	configTracer := config.EnvConfig.Tracer

	if !configTracer.Enable {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(configTracer.Url),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Bearer " + configTracer.SecretKey,
		}),
	)
	if err != nil {
		return fmt.Errorf("create trace exporter: %w", err)
	}

	re := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(configTracer.ServiceName),
		attribute.String("deployment.environment", configTracer.Environment),
	)

	sampler := trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(configTracer.SampleRatio)))

	tp := trace.NewTracerProvider(
		sampler,
		trace.WithBatcher(exporter),
		trace.WithResource(re),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	AMPTracerProvider = tp

	fmt.Println("Success init tracer OTEL")
	return nil
}

func Shutdown(ctx context.Context) error {
	if AMPTracerProvider == nil {
		return nil
	}
	return AMPTracerProvider.Shutdown(ctx)
}
