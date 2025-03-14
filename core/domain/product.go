package domain

import "time"

type Product struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	SKU       string     `json:"sku"`
	Quantity  uint       `json:"quantity"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
