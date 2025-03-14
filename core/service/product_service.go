package service

import (
	"context"
	"math"
	"time"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type productService struct {
	repo port.IProductRepo
}

func NewProductService(repo port.IProductRepo) port.IProductService {
	return &productService{repo: repo}
}

func (p *productService) Add(ctx context.Context, in *domain.Product) (err error) {
	return p.repo.Add(ctx, in)
}

func (p *productService) Delete(ctx context.Context, productId uint) (err error) {
	return p.repo.Delete(ctx, productId)
}

func (p *productService) Get(ctx context.Context, productId uint) (product domain.Product, err error) {
	return p.repo.Get(ctx, productId)
}

func (p *productService) Update(ctx context.Context, in *port.ProductUpdate) (err error) {
	var product domain.Product
	now := time.Now()

	if product, err = p.Get(ctx, in.ID); err != nil {
		return
	}

	if in.Name != nil {
		product.Name = *in.Name
	}
	if in.SKU != nil {
		product.SKU = *in.SKU
	}
	if in.Quantity != nil {
		product.Quantity = *in.Quantity
	}

	product.UpdatedAt = &now

	if err = p.repo.Update(ctx, &product); err != nil {
		return
	}
	return
}

func (p *productService) List(ctx context.Context, in *port.ProductList) (err error) {
	if in.Limit < 1 {
		in.Limit = 10
	}
	if in.Page < 1 {
		in.Page = 1
	}

	if err = p.repo.GetTotalRow(ctx, in); err != nil {
		return
	}

	in.Offset = (in.Page - 1) * in.Limit
	in.TotalPage = uint(math.Ceil(float64(in.TotalRow) / float64(in.Limit)))

	if err = p.repo.List(ctx, in); err != nil {
		return
	}
	return
}
