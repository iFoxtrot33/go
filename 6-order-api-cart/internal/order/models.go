package order

import (
	"order-api/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Phone       string            `json:"phone"`
	Products    []product.Product `gorm:"many2many:order_products;"`
	Description string            `json:"description"`
}

type OrderProduct struct {
	OrderID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  uint `json:"quantity"`
}

type OrderProductItem struct {
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}
