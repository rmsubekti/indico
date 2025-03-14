package service

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type orderService struct {
	repo          port.IOrderRepo
	detailService port.IOrderDetailService
	stockService  port.IStockService
}

func NewOrderService(repo port.IOrderRepo) port.IOrderService {
	return &orderService{repo: repo}
}

func (o *orderService) WithOrderDetailService(detServ port.IOrderDetailService) port.IOrderService {
	o.detailService = detServ
	return o
}

func (o *orderService) WithStockService(stockServ port.IStockService) port.IOrderService {
	o.stockService = stockServ
	return o
}

func (o *orderService) Add(ctx context.Context, in *domain.Order) (err error) {
	var details []domain.OrderDetail

	in.OrderStatus = domain.StatusOrderOpen
	if err = o.repo.Add(ctx, in); err != nil {
		return
	}

	for _, v := range *in.Details {
		v.OrderID = in.ID
		details = append(details, v)
	}

	if err = o.detailService.AddAll(ctx, &details); err != nil {
		return
	}
	in.Details = &details
	return
}

func (o *orderService) Update(ctx context.Context, in *domain.Order) error {
	return o.repo.Update(ctx, in)
}

func (o *orderService) ChangeStatus(ctx context.Context, orderID uint, status domain.OrderStatus) (err error) {
	order, _ := o.repo.GetByID(ctx, orderID)
	order.OrderStatus = status

	if status == domain.StatusOrderCompleted {
		err = o.detailService.GetAll(ctx, orderID, order.Details)
		for _, v := range *order.Details {
			if err != nil {
				break
			}
			if order.FromWarehouseID != nil {
				err = o.stockService.Update(ctx, v.ProductID, *order.FromWarehouseID, -int(v.Quantity))
			}
			if order.ToWarehouseID != nil {
				err = o.stockService.Update(ctx, v.ProductID, *order.FromWarehouseID, int(v.Quantity))
			}
		}

		if err != nil {
			return
		}

	}

	return o.repo.Update(ctx, &order)
}

func (o *orderService) Get(ctx context.Context, in *domain.Order) error {
	panic("unimplemented")
}

func (o *orderService) List(ctx context.Context, in *port.OrderList) error {
	panic("unimplemented")
}
