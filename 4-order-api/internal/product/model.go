package product

import (
	"math/rand"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Price       string         `json:"price" gorm:"type:varchar(50);not null"`
	Images      pq.StringArray `json:"images" gorm:"type:text[];not null"`
}

func NewProduct(req *ProductCreateRequest) *Product {
	return &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Images:      req.Images,
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
