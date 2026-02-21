package datastore

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	//cloudsqlconn "github.com/GoogleCloudPlatform/cloudsql-proxy/cloudsqlconn"
	cloudsql "cloud.google.com/go/cloudsqlconn"
	otelsql "github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func ConnectDB() *sql.DB {
	ctx := context.Background()
	// Create a Cloud SQL Connector dialer
	d, err := cloudsql.NewDialer(ctx, cloudsql.WithIAMAuthN())
	if err != nil {
		log.Fatalf("failed to create dialer: %v", err)
	}
	defer d.Close()
	// Instance connection name: PROJECT:REGION:INSTANCE
	instanceConnName := "open-telemetry-1000:us-central1:open-tel-db"
	// Build a DSN string (no password needed if using IAM)
	dsn := fmt.Sprintf("dbname=open-tel-db sslmode=disable")
	// Instead of cloudsql.Driver, you pass the dialer into pgx
	connStr := fmt.Sprintf("host=%s user=postgres %s", instanceConnName, dsn)
	// Open DB with otelsql wrapping the driver
	db, err := otelsql.Open("pgx", connStr, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		log.Fatal(err)
	}
	_, err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		log.Printf("Could not register db metrics: %v", err)
	}
	return db
}
