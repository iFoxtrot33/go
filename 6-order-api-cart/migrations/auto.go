package main

import (
	"os"

	"order-api/internal/order"
	"order-api/internal/product"
	"order-api/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db.AutoMigrate(&product.Product{}, &user.User{}, &order.Order{}, &order.OrderProduct{})
}
