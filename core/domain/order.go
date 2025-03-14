package domain

import "time"

type Order struct {
	ID              uint           `json:"id"`
	FromWarehouseID *uint          `json:"from_warehouse_id"`
	ToWarehouseID   *uint          `json:"to_warehouse_id"`
	OrderType       OrderType      `json:"order_type"`
	OrderStatus     OrderStatus    `json:"order_status"`
	Note            string         `json:"note"`
	Details         *[]OrderDetail `json:"details"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at"`
}

type OrderType string

const TypeReceive OrderType = "Receive"
const TypeShipping OrderType = "Shipping"

type OrderStatus string

const StatusOrderOpen OrderStatus = "OPEN"
const StatusOrderProcessing OrderStatus = "PROCESSING"
const StatusOrderCompleted OrderStatus = "COMPLETED"
const StatusOrderCanceled OrderStatus = "CANCELED"
