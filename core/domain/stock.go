package domain

type Stock struct {
	ID          uint `json:"id"`
	WarehouseID uint `json:"warehouse_id"`
	ProductID   uint `json:"product_id"`
	Quantity    uint `json:"quantity"`
}
