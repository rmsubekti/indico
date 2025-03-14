package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	WarehouseList struct {
		domain.Pagination
		Rows *[]domain.Warehouse `json:"rows,omitempty"`
	}

	IWarehouseRepo interface {
		Add(ctx context.Context, in *domain.Warehouse) error
		GetTotalRow(ctx context.Context, in *WarehouseList) error
		List(ctx context.Context, in *WarehouseList) error
	}
	IWarehouseService interface {
		Add(ctx context.Context, in *domain.Warehouse) error
		List(ctx context.Context, in *WarehouseList) error
	}
)
