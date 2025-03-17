package repo

import (
	"context"
	"database/sql"

	"github.com/rmsubekti/indico/core/domain"
)

type StockRepo struct {
	Tx *sql.Tx
}

func (s *StockRepo) Add(ctx context.Context, in *domain.Stock) (err error) {
	return s.Tx.QueryRowContext(ctx, `
	insert into stock (warehouse_id,product_id,quantity) values ($1,$2,$3) returning id 
`, in.WarehouseID, in.ProductID, in.Quantity).Scan(&in.ID)
}

func (s *StockRepo) GetByWarehouseAndProductID(ctx context.Context, wid uint, productId uint) (out domain.Stock, err error) {
	err = s.Tx.QueryRowContext(ctx, `
	select id, warehouse_id, product_id, quantity from stock where warehouse_id = $1 and product_id=$2 limit 1
`, wid, productId).Scan(&out.ID, &out.WarehouseID, &out.ProductID, &out.Quantity)
	return
}

func (s *StockRepo) UpdateQty(ctx context.Context, id uint, qty uint) (err error) {
	return s.Tx.QueryRowContext(ctx, `
	update stock set quantity=$1 where id=$2 
	`, qty, id).Err()
}
