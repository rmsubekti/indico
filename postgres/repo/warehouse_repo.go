package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type WarehouseRepo struct {
	Tx *sql.Tx
}

func (w *WarehouseRepo) Add(ctx context.Context, in *domain.Warehouse) error {
	return w.Tx.QueryRowContext(ctx, `
		insert into warehouse (name,address,capacity) values ($1,$2,$3) returning id 
	`, in.Name, in.Address, in.Capacity).Scan(&in.ID)
}

func (w *WarehouseRepo) GetTotalRow(ctx context.Context, in *port.WarehouseList) error {
	var stmn strings.Builder
	stmn.WriteString("select count(id) from warehouse where 1=1 ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	return w.Tx.QueryRowContext(ctx, stmn.String()).Scan(&in.TotalRow)
}

func (w *WarehouseRepo) List(ctx context.Context, in *port.WarehouseList) (err error) {
	var (
		rows       *sql.Rows
		warehouses []domain.Warehouse
		stmn       strings.Builder
	)

	stmn.WriteString("select id, name,address,capacity from warehouse where 1=1 ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	stmn.WriteString(fmt.Sprintf("limit %d offset %d ", in.Limit, in.Offset))

	if rows, err = w.Tx.QueryContext(ctx, stmn.String()); err != nil {
		return
	}
	for rows.Next() {
		var warehouse domain.Warehouse
		if err = rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Address, &warehouse.Capacity); err != nil {
			return
		}
		warehouses = append(warehouses, warehouse)
	}

	in.Rows = &warehouses
	return
}
