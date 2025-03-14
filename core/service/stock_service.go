package service

import (
	"context"
	"errors"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type stockService struct {
	repo port.IStockRepo
}

func NewStockService(repo port.IStockRepo) port.IStockService {
	return &stockService{repo: repo}
}

func (s *stockService) Update(ctx context.Context, productID, waID uint, qty int) (err error) {
	var stock domain.Stock

	stock, _ = s.repo.GetByWarehouseAndProductID(ctx, waID, productID)

	if stock.ID == 0 {
		if qty < 1 {
			return errors.New("kuantitas harus lebih besar dari nol")
		}
		stock.ProductID = productID
		stock.Quantity = uint(qty)
		stock.WarehouseID = waID
		if err = s.repo.Add(ctx, &stock); err != nil {
			return
		}
	} else {
		if (int(stock.Quantity) + qty) < 1 {
			return errors.New("kuantitas harus lebih besar dari nol")
		}
		stock.Quantity += uint(qty)
		if err = s.repo.UpdateQty(ctx, stock.ID, stock.Quantity); err != nil {
			return
		}
	}

	return
}
