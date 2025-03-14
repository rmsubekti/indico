package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	IStockRepo interface {
		Add(ctx context.Context, in *domain.Stock) error
		UpdateQty(ctx context.Context, id, qty uint) error
		GetByWarehouseAndProductID(ctx context.Context, wid, productId uint) (domain.Stock, error)
	}
	IStockService interface {
		Update(ctx context.Context, productID, waID uint, qty int) error
	}
)
