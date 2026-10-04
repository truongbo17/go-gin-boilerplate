package tracer

import (
	"context"
	"fmt"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"go.elastic.co/apm/v2"
	"go.elastic.co/apm/v2/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.30.0"
	"log"
	"net/url"
	"os"
	"time"
)

func IniTracerAPMConfig() {
	configTracer := config.EnvConfig.Tracer.APM

	_ = os.Setenv("ELASTIC_APM_SERVER_URL", configTracer.Url)
	_ = os.Setenv("ELASTIC_APM_SECRET_TOKEN", configTracer.SecretKey)
	_ = os.Setenv("ELASTIC_APM_SERVICE_NAME", configTracer.ServiceName)
	_ = os.Setenv("ELASTIC_APM_ACTIVE", "true")
}

var APMTracer *apm.Tracer

func InitTracer() *apm.Tracer {
	configTracer := config.EnvConfig.Tracer

	serverUrl, _ := url.Parse(configTracer.APM.Url)
	tr, err := transport.NewHTTPTransport(transport.HTTPTransportOptions{
		ServerURLs:    []*url.URL{serverUrl},
		SecretToken:   configTracer.APM.SecretKey,
		ServerTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error init tracer APM with HTTP Transport: %v", err)
	}

	tracer, err := apm.NewTracerOptions(apm.TracerOptions{
		ServiceName:    configTracer.APM.ServiceName,
		ServiceVersion: "1.0.0",
		Transport:      tr,
	})

	if err != nil {
		log.Fatalf("Error init Tracer APM: %v", err)
	}

	log.Println("Successfully init tracer APM!")
	APMTracer = tracer

	return tracer
}

var AMPTracerProvider *trace.TracerProvider

func InitTracerOTEL() {
	configTracer := config.EnvConfig.Tracer

	if !configTracer.Enable {
		return
	}
	ctx := context.Background()
	IniTracerAPMConfig()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(configTracer.Url),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Bearer " + configTracer.SecretKey,
		}),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	re := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(configTracer.ServiceName),
		attribute.String("deployment.environment", configTracer.Environment),
	)

	sampler := trace.WithSampler(trace.AlwaysSample())

	tp := trace.NewTracerProvider(
		sampler,
		trace.WithBatcher(exporter),
		trace.WithResource(re),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	AMPTracerProvider = tp

	fmt.Println("Success init tracer OTEL")
}

func DownAMPTracerProvider() {
	if AMPTracerProvider == nil {
		return
	}
	if err := AMPTracerProvider.Shutdown(context.Background()); err != nil {
		log.Fatal("Failed to shutdown TracerProvider: ", err)
	}
}
