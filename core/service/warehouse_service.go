package service

import (
	"context"
	"math"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type warehouseService struct {
	repo port.IWarehouseRepo
}

func NewWarehouseService(repo port.IWarehouseRepo) port.IWarehouseService {
	return &warehouseService{repo: repo}
}

func (w *warehouseService) Add(ctx context.Context, in *domain.Warehouse) (err error) {
	if err = w.repo.Add(ctx, in); err != nil {
		return
	}
	return
}

func (w *warehouseService) List(ctx context.Context, in *port.WarehouseList) (err error) {
	if in.Limit < 1 {
		in.Limit = 10
	}
	if in.Page < 1 {
		in.Page = 1
	}

	if err = w.repo.GetTotalRow(ctx, in); err != nil {
		return
	}

	in.Offset = (in.Page - 1) * in.Limit
	in.TotalPage = uint(math.Ceil(float64(in.TotalRow) / float64(in.Limit)))

	if err = w.repo.List(ctx, in); err != nil {
		return
	}
	return
}
