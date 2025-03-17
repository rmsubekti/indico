package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type orderService struct {
	repo          port.IOrderRepo
	detailService port.IOrderDetailService
	stockService  port.IStockService
	whService     port.IWarehouseService
	prodService   port.IProductService
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

func (o *orderService) WithProductService(prodServ port.IProductService) port.IOrderService {
	o.prodService = prodServ
	return o
}

func (o *orderService) WithWarehouseService(whService port.IWarehouseService) port.IOrderService {
	o.whService = whService
	return o
}

func (o *orderService) Add(ctx context.Context, in *domain.Order) (err error) {
	if o.detailService == nil || o.stockService == nil {
		return errors.New("perlu detail service ")
	}
	if err = in.Valid(); err != nil {
		return
	}
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

func (o *orderService) Update(ctx context.Context, in *domain.Order) (err error) {
	if err = in.Valid(); err != nil {
		return
	}
	now := time.Now()
	in.UpdatedAt = &now
	return o.repo.Update(ctx, in)
}

func (o *orderService) ChangeStatus(ctx context.Context, orderID uint, status domain.OrderStatus) (err error) {
	if o.detailService == nil || o.stockService == nil {
		return errors.New("perlu stock service & detail service ")
	}

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
			order.Note = fmt.Sprintf("terjadi kesalahan. err: %v", err)
			return
		}

	}

	now := time.Now()
	order.UpdatedAt = &now
	return o.repo.Update(ctx, &order)
}

func (o *orderService) GetByID(ctx context.Context, id uint) (out domain.Order, err error) {
	if o.whService == nil || o.detailService == nil || o.prodService == nil {
		err = errors.New("perlu detail service, warehouse service dan product sercive")
		return

	}
	if out, err = o.repo.GetByID(ctx, id); err != nil {
		return
	}

	if out.FromWarehouseID != nil {
		var wh domain.Warehouse
		if wh, err = o.whService.GetByID(ctx, *out.FromWarehouseID); err != nil {
			return
		}
		out.FromWarehouse = &wh

	}
	if out.ToWarehouseID != nil {
		var wh domain.Warehouse
		if wh, err = o.whService.GetByID(ctx, *out.ToWarehouseID); err != nil {
			return
		}
		out.FromWarehouse = &wh
	}

	if err = o.detailService.GetAll(ctx, id, out.Details); err != nil {
		return
	}

	var details []domain.OrderDetail
	for _, v := range *out.Details {
		var product domain.Product
		if product, err = o.prodService.Get(ctx, v.ProductID); err != nil {
			break
		}
		v.Product = product
		details = append(details, v)
	}
	out.Details = &details
	return
}

func (o *orderService) List(ctx context.Context, in *port.OrderList) (err error) {
	if o.whService == nil {
		return errors.New("perlu warehouse servuce")

	}
	if in.Limit < 1 {
		in.Limit = 10
	}

	if in.Page < 1 {
		in.Page = 1
	}

	if err = o.repo.GetTotalRow(ctx, in); err != nil {
		return
	}

	in.Offset = (in.Page - 1) * in.Limit
	in.TotalPage = uint(math.Ceil(float64(in.TotalRow) / float64(in.Limit)))

	if err = o.repo.List(ctx, in); err != nil {
		return
	}

	var rows []domain.Order
	for _, v := range *in.Rows {
		if v.FromWarehouseID != nil {
			var wh domain.Warehouse
			if wh, err = o.whService.GetByID(ctx, *v.FromWarehouseID); err != nil {
				break
			}
			v.FromWarehouse = &wh
			rows = append(rows, v)
		}
		if v.ToWarehouseID != nil {
			var wh domain.Warehouse
			if wh, err = o.whService.GetByID(ctx, *v.ToWarehouseID); err != nil {
				break
			}

			v.FromWarehouse = &wh
			rows = append(rows, v)
		}
	}

	in.Rows = &rows
	return
}
