package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	OrderList struct {
		domain.Pagination
		Rows *[]domain.Order `json:"rows,omitempty"`
	}
	IOrderRepo interface {
		Add(ctx context.Context, in *domain.Order) error
		GetByID(ctx context.Context, id uint) (domain.Order, error)
		Update(ctx context.Context, in *domain.Order) error
		GetTotalRow(ctx context.Context, in *OrderList) error
		List(ctx context.Context, in *OrderList) error
	}
	IOrderService interface {
		WithWarehouseService(whServ IWarehouseService) IOrderService
		WithProductService(prodServ IProductService) IOrderService
		WithOrderDetailService(detServ IOrderDetailService) IOrderService
		WithStockService(stockServ IStockService) IOrderService
		Add(ctx context.Context, in *domain.Order) error
		ChangeStatus(ctx context.Context, orderID uint, status domain.OrderStatus) error
		GetByID(ctx context.Context, id uint) (domain.Order, error)
		Update(ctx context.Context, in *domain.Order) error
		List(ctx context.Context, in *OrderList) error
	}
)
