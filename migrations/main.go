package main

import (
	"bike/internal/addresses"
	"bike/internal/products"
	"bike/internal/users"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Ждем подключения к БД (важно в Docker)
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN environment variable is required")
	}

	var db *gorm.DB
	var err error

	// Пытаемся подключиться несколько раз
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		log.Fatal("Failed to connect to database after multiple attempts:", err)
	}

	// Выполняем миграции
	err = db.AutoMigrate(
		&products.Product{},
		&users.User{},
		&addresses.Address{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("✅ Database migrated successfully!")
}
