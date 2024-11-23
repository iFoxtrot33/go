package test

import (
	"order-api/internal/product"
	"order-api/internal/user"

	"gorm.io/gorm"
)

func InitData(db *gorm.DB) {
	RemoveData(db)
	db.Create(&user.User{
		Phone:     "+79993333333",
		Name:      "Test User",
		SessionId: "0efb242e39666525383635acd2a127fc97f7968f60191ff7275f0a079c74d97b",
	})

	db.Create(&product.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       "100",
		Images:      []string{"https://example.com/image.jpg", "https://example.com/image2.jpg"},
	})

	db.Create(&product.Product{
		Name:        "Test Product 2",
		Description: "Test Description 2",
		Price:       "200",
		Images:      []string{"https://example.com/image.jpg", "https://example.com/image2.jpg"},
	})

}
