package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	ProductUpdate struct {
		ID       uint    `json:"id,omitempty"`
		Name     *string `json:"name,omitempty"`
		SKU      *string `json:"sku,omitempty"`
		Quantity *uint   `json:"quantity,omitempty"`
	}
	ProductList struct {
		domain.Pagination
		Rows *[]domain.Product `json:"rows,omitempty"`
	}
	IProductRepo interface {
		Add(ctx context.Context, in *domain.Product) error
		Update(ctx context.Context, in *domain.Product) error
		Get(ctx context.Context, productId uint) (domain.Product, error)
		Delete(ctx context.Context, productId uint) error
		GetTotalRow(ctx context.Context, in *ProductList) error
		List(ctx context.Context, in *ProductList) error
	}

	IProductService interface {
		Add(ctx context.Context, in *domain.Product) error
		Update(ctx context.Context, in *ProductUpdate) error
		Delete(ctx context.Context, productId uint) error
		Get(ctx context.Context, productId uint) (domain.Product, error)
		List(ctx context.Context, in *ProductList) error
	}
)
