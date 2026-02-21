package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// The HTTP Middleware already started a span. We get the context.
	ctx := r.Context()

	// Manual child span for specific logic
	tr := otel.Tracer("handlers")
	ctx, span := tr.Start(ctx, "fetch_products_logic")
	defer span.End()
	rows, err := h.db.QueryContext(ctx, "SELECT id, name FROM products")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "db query failed")
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	// ... (Simulate processing rows) ...
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[{"id":1, "name":"Otel T-Shirt"}]`))
}
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("handlers")

	// Start a span
	ctx, span := tr.Start(ctx, "process_order")
	defer span.End()
	// Simulate Order Processing
	h.ProcessOrder(ctx, "order-123") // Call internal function
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "created"}`))
}
func (h *Handler) ProcessOrder(ctx context.Context, orderID string) {
	tr := otel.Tracer("core-logic")
	_, span := tr.Start(ctx, "payment_gateway_call")
	defer span.End()
	span.SetAttributes(attribute.String("order.id", orderID))
	// Simulate a slow payment
	time.Sleep(200 * time.Millisecond)

	// Simulate an error 10% of the time to test observability
	if time.Now().Unix()%10 == 0 {
		err := fmt.Errorf("payment gateway timeout")
		span.RecordError(err)
		span.SetStatus(codes.Error, "payment failed")
	}
}
