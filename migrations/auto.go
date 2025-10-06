package main

import (
	"bike/internal/addresses"
	"bike/internal/products"
	"bike/internal/users"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&products.Product{}, &users.User{}, &addresses.Address{})
	fmt.Println("Database migrated")
}
