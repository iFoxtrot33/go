package order

import (
	"fmt"
	"order-api/pkg/db"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Database *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{
		Database: database,
	}
}

func (repo *OrderRepository) Create(order *Order, products []OrderProductItem) (*Order, error) {
	err := repo.Database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, item := range products {
			orderProduct := OrderProduct{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
			}

			if err := tx.Create(&orderProduct).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	err = repo.Database.DB.Preload("Products").First(order, order.ID).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (repo *OrderRepository) GetById(id uint) (*Order, error) {
	var order Order
	result := repo.Database.DB.Preload("Products").First(&order, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order with id %d not found", id)
		}
		return nil, result.Error
	}

	return &order, nil
}

func (repo *OrderRepository) GetByPhone(phone string) ([]Order, error) {
	var orders []Order
	result := repo.Database.DB.Preload("Products").
		Where("phone = ?", phone).
		Order("created_at DESC").
		Find(&orders)

	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func (repo *OrderRepository) IsOrderOwner(orderID uint, phone string) bool {
	var count int64
	repo.Database.DB.Model(&Order{}).
		Where("id = ? AND phone = ?", orderID, phone).
		Count(&count)

	return count > 0
}

func (repo *OrderRepository) GetOrderDetails(id uint) (*Order, error) {
	var order Order
	result := repo.Database.DB.
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Select("products.*, order_products.quantity")
		}).
		First(&order, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order with id %d not found", id)
		}
		return nil, result.Error
	}

	return &order, nil
}
