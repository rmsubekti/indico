package domain

type OrderDetail struct {
	ID        uint `json:"id"`
	OrderID   uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"quantity"`
}
