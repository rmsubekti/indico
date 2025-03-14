package service

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type orderDetailService struct {
	repo port.IOrderDetailRepo
}

func NewOrderDetailService(repo port.IOrderDetailRepo) port.IOrderDetailService {
	return &orderDetailService{repo: repo}
}

func (o *orderDetailService) AddAll(ctx context.Context, details *[]domain.OrderDetail) error {
	return o.repo.AddAll(ctx, details)
}

func (o *orderDetailService) GetAll(ctx context.Context, order_id uint, details *[]domain.OrderDetail) error {
	return o.repo.GetAll(ctx, order_id, details)
}
