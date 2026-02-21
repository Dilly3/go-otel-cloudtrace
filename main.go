package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dilly3/go-otel-cloudtrace/api"
	"github.com/dilly3/go-otel-cloudtrace/datastore"
	"github.com/dilly3/go-otel-cloudtrace/tracer"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	tp := tracer.InitTracerMulti()
	defer func() { _ = tp.Shutdown(context.Background()) }()
	// 1. CONNECT DB (Instrumented!)
	db := datastore.ConnectDB()
	// 2. ROUTER
	r := mux.NewRouter()
	h := api.NewHandler(db)
	r.Use(otelhttp.NewMiddleware("otel-mart-server")) // Auto-instrument HTTP
	r.HandleFunc("/products", h.GetProducts).Methods("GET")
	r.HandleFunc("/order", h.CreateOrder).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
