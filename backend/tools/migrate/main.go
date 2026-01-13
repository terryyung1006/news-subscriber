package main

import (
	"flag"
	"fmt"
	"log"
	"news-subscriber-core/src/config"
	"news-subscriber-core/src/repository/postgres"
)

func main() {
	action := flag.String("action", "migrate", "Action to perform: migrate, reset, or seed")
	flag.Parse()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	switch *action {
	case "migrate":
		fmt.Println("Running database migrations...")
		if err := db.Migrate(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✓ Migrations completed successfully")

	case "reset":
		fmt.Println("Resetting database (dropping all tables)...")
		if err := db.Reset(); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		fmt.Println("✓ Database reset completed successfully")

	case "seed":
		fmt.Println("Seeding initial data...")
		if err := db.Seed(); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("✓ Seeding completed successfully")

	case "fresh":
		fmt.Println("Running fresh setup (reset + migrate + seed)...")
		if err := db.Reset(); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		if err := db.Seed(); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		fmt.Println("✓ Fresh setup completed successfully")

	default:
		log.Fatalf("Unknown action: %s. Use: migrate, reset, seed, or fresh", *action)
	}
}

