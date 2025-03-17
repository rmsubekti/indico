package domain

import (
	"errors"
	"time"
)

type Order struct {
	ID              uint           `json:"id"`
	FromWarehouseID *uint          `json:"from_warehouse_id"`
	FromWarehouse   *Warehouse     `json:"from_warehouse"`
	ToWarehouseID   *uint          `json:"to_warehouse_id"`
	ToWarehouse     *Warehouse     `json:"to_warehouse"`
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

func (o Order) Valid() (err error) {
	if o.FromWarehouseID == nil && o.OrderType == TypeShipping {
		return errors.New("warehouse asal tidak boleh kosong untuk shipping")
	}
	if o.ToWarehouseID == nil && o.OrderType == TypeReceive {
		return errors.New("warehouse tujuan tidak boleh kosong untuk receive")
	}
	return
}
