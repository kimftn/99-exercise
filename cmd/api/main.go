package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"property-api/internal/app"
	"property-api/internal/database/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file loaded: %v", err)
	}

	ctx := context.Background()
	pgConfig := postgres.LoadConfigFromEnv()
	if !pgConfig.IsConfigured() {
		log.Fatal("postgres configuration is required for user APIs; set DATABASE_URL or PGHOST/PGPORT/PGUSER/PGPASSWORD/PGDATABASE")
	}

	var pool *pgxpool.Pool
	var err error
	pool, err = postgres.NewPool(ctx, pgConfig)
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer pool.Close()

	log.Println("postgres connection established")

	listingServiceBaseURL := os.Getenv("LISTING_SERVICE_BASE_URL")
	if listingServiceBaseURL == "" {
		listingServiceBaseURL = "http://localhost:6000"
	}

	server := app.NewServerWithPool(pool, listingServiceBaseURL)

	log.Println("starting server on :3000")
	if err := server.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
