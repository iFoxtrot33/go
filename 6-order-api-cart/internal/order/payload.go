package order

type OrderCreateRequest struct {
	Products    []OrderProductItem `json:"products"`
	Description string             `json:"description"`
}
