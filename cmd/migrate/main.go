package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"property-api/internal/database/migration"
	"property-api/internal/database/postgres"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file loaded: %v", err)
	}

	config := postgres.LoadConfigFromEnv()
	if !config.IsConfigured() {
		log.Fatal("database configuration is missing; set DATABASE_URL or PGHOST/PGPORT/PGUSER/PGPASSWORD/PGDATABASE")
	}

	runner, err := migration.NewRunner(config.DSN())
	if err != nil {
		log.Fatalf("create migration runner: %v", err)
	}
	defer func() {
		if err := runner.Close(); err != nil {
			log.Printf("close migration runner: %v", err)
		}
	}()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := runner.Up(); err != nil {
			log.Fatalf("run up migrations: %v", err)
		}
		log.Println("migrations applied successfully")
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			parsedSteps, err := strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("invalid down steps: %v", err)
			}
			steps = parsedSteps
		}

		if err := runner.Down(steps); err != nil {
			log.Fatalf("run down migrations: %v", err)
		}
		log.Printf("rolled back %d migration step(s)", steps)
	case "version":
		version, dirty, err := runner.Version()
		if err != nil {
			log.Fatalf("get migration version: %v", err)
		}
		log.Printf("migration version=%d dirty=%t", version, dirty)
	default:
		log.Fatalf("unsupported command %q. use: up | down [steps] | version", command)
	}

	fmt.Println("migration command completed")
}
