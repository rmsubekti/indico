package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type OrderRepo struct {
	Tx *sql.Tx
}

func (o *OrderRepo) Add(ctx context.Context, in *domain.Order) (err error) {
	return o.Tx.QueryRowContext(ctx, `
		insert into order (from_warehouse_id,to_warehouse_id,type, status, note) values ($1,$2,$3,$4,$5) returning id 
	`, in.FromWarehouse, in.ToWarehouse, in.OrderType, in.OrderStatus, in.Note).Scan(&in.ID)
}

func (o *OrderRepo) GetByID(ctx context.Context, id uint) (out domain.Order, err error) {
	err = o.Tx.QueryRowContext(ctx, `
	select id, from_warehouse_id, to_warehouse_id, type,status, note, created_at, updated_at, deleted_at from order where id = $1 limit 1
`, id).Scan(&out.ID, &out.FromWarehouseID, &out.ToWarehouseID, &out.OrderType, &out.OrderType, &out.Note, &out.CreatedAt, &out.UpdatedAt, &out.DeletedAt)
	return
}

func (o *OrderRepo) GetTotalRow(ctx context.Context, in *port.OrderList) (err error) {
	var stmn strings.Builder
	stmn.WriteString("select count(id) from order where deleted_at is null ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and note ilike '%%%s%%'", in.Search))
	}
	return o.Tx.QueryRowContext(ctx, stmn.String()).Scan(&in.TotalRow)
}

func (o *OrderRepo) List(ctx context.Context, in *port.OrderList) (err error) {
	var (
		rows   *sql.Rows
		orders []domain.Order
		stmn   strings.Builder
	)
	stmn.WriteString("select id, from_warehouse_id,to_warehouse_id, type, status, note, created_at,updated_at,deleted_at from order where deleted_at is null ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and note ilike '%%%s%%'", in.Search))
	}
	stmn.WriteString(fmt.Sprintf("limit %d offset %d ", in.Limit, in.Offset))

	if rows, err = o.Tx.QueryContext(ctx, stmn.String()); err != nil {
		return
	}
	for rows.Next() {
		var order domain.Order
		if err = rows.Scan(&order.ID, &order.FromWarehouseID, &order.ToWarehouseID, &order.OrderType, &order.OrderStatus, &order.Note, &order.CreatedAt, &order.UpdatedAt, &order.DeletedAt); err != nil {
			return
		}
		orders = append(orders, order)
	}
	in.Rows = &orders
	return
}

func (o *OrderRepo) Update(ctx context.Context, in *domain.Order) (err error) {
	err = o.Tx.QueryRowContext(ctx, `
	update product set from_warehouse_id=$1, to_warehouse_id=$2, type=$3, status=$4, note=$5, updated_at=$6 where id=$7 
	`, in.FromWarehouseID, in.ToWarehouseID, in.OrderType, in.OrderStatus, in.Note, in.UpdatedAt, in.ID).Err()
	return
}
