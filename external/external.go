package external

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
)

func ChildFunction(ctx context.Context) {
	// We do NOT create a new tracer provider. We just call otel.Tracer.
	tracer := otel.Tracer("demo-app")
	// 3. Create a Child Span
	// Because we passed 'ctx' from main, OTel knows this relies on "main-operation"
	_, span := tracer.Start(ctx, "child-operation")
	defer span.End()
	fmt.Println("Doing work in child function...")
	time.Sleep(50 * time.Millisecond)
}
