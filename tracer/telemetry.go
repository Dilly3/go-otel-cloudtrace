// package tracer
package tracer

import (
	"context"
	"log"
	"os"

	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitTracerJaeger() *trace.TracerProvider {
	grpcURL := "localhost:4317"
	ctx := context.Background()
	// 1. Create the OTLP Exporter (gRPC)
	// We point it to localhost:4317 where Jaeger is listening
	// Insecure is needed because we aren't using SSL on localhost
	// Port 4317 is for gRPC (Binary, fast, persistent connection). Use this for Go backend services.
	// Port 4318 is for HTTP (JSON/Protobuf over HTTP). Use this for JavaScript/Frontend apps that can't do gRPC easily
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(grpcURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 2. Define the Resource (Who am I?)
	// This adds labels to every trace like "service.name = ecommerce-api
	// If you don't do this, Jaeger will just say "unknown_service".
	// When you have 50 microservices, you need to know which one generated the error.
	// We use semconv (Semantic Conventions) to ensure we name things the standard way.
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("ecommerce-api-local"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 3. Create the Provider with the Exporter and Resource
	tp := trace.NewTracerProvider(
		// The "Batcher": we used trace.WithBatcher(exporter). This is crucial.
		// It means the app doesn't send a network request to Jaeger for every single span.
		// It collects them in memory for a few seconds and sends them in a bundle.
		// This prevents the tracing tool from slowing down your APP
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp
}

func InitTracer() *trace.TracerProvider {
	// 1. Create an exporter that prints to terminal
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	// 2. Create the TracerProvider (The Factory)
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	// 3. Set the global provider so we can call otel.Tracer anywhere
	otel.SetTracerProvider(tp)

	return tp
}

func InitTracerMulti() *trace.TracerProvider {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	env := os.Getenv("APP_ENV") // "dev" or "prod"
	var exporter trace.SpanExporter
	var err error
	// 1. CHOOSE YOUR EXPORTER
	if env == "prod" {
		// --- GOOGLE CLOUD TRACE ---
		exporter, err = texporter.New(texporter.WithProjectID(projectID))
		if err != nil {
			log.Fatalf("Failed to create GCP exporter: %v", err)
		}
	} else {
		// --- JAEGER (Local) ---
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint("localhost:4317"),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			log.Fatalf("Failed to create Jaeger exporter: %v", err)
		}
	}
	// 2. RESOURCE IDENTIFICATION
	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("otel-mart-api"),
			semconv.DeploymentEnvironmentName(env),
		),
	)
	// 3. CREATE PROVIDER
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		// Sample 100% in local, 10% in prod
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(
			func() float64 {
				if env == "production" {
					return 0.1
				}
				return 1.0
			}(),
		))),
	)
	// 4. SET GLOBALS
	otel.SetTracerProvider(tp)

	// Enables passing trace context in headers (W3C standard)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp
}
