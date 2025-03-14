package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	IOrderDetailRepo interface {
		AddAll(ctx context.Context, details *[]domain.OrderDetail) error
		GetAll(ctx context.Context, order_id uint, details *[]domain.OrderDetail) error
	}
	IOrderDetailService interface {
		AddAll(ctx context.Context, details *[]domain.OrderDetail) error
		GetAll(ctx context.Context, order_id uint, details *[]domain.OrderDetail) error
	}
)
