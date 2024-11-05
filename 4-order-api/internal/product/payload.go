package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Price       string         `json:"price" validate:"required"`
	Images      pq.StringArray `json:"images" validate:"required"`
}
