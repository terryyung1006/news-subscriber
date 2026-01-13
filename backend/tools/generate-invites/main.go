package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"news-subscriber-core/src/config"
	"news-subscriber-core/src/models"
	"news-subscriber-core/src/repository/postgres"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()

	fmt.Println("Invite Code Generator")
	fmt.Println("====================")
	fmt.Println()

	// Generate 5 invite codes
	for i := 0; i < 5; i++ {
		code := generateCode(12)

		inviteCode := &models.InviteCode{
			Code: code,
		}

		if err := db.WithContext(context.Background()).Create(inviteCode).Error; err != nil {
			log.Printf("Failed to insert code %s: %v", code, err)
			continue
		}

		fmt.Printf("✅ Generated: %s\n", code)
	}

	fmt.Println()
	fmt.Println("Done! These codes can now be used for signup.")
}
