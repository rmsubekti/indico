package repo

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"

	"github.com/rmsubekti/indico/core/domain"
)

type OrderDetailRepo struct {
	Tx *sql.Tx
}

func (o *OrderDetailRepo) AddAll(ctx context.Context, details *[]domain.OrderDetail) (err error) {
	var (
		stmn bytes.Buffer
		vArg = []any{}
	)
	stmn.WriteString("insert into order_detail(order_id,product_id,quantity) values ")
	for i, v := range *details {
		vArg = append(vArg, v.OrderID, v.ProductID, v.Quantity)

		numField := 3
		n := i * numField
		stmn.WriteString("(")
		for j := 0; j < numField; j++ {
			stmn.WriteString("$")
			stmn.WriteString(strconv.Itoa(n + j + 1))
			stmn.WriteString(",")
		}
		stmn.Truncate(stmn.Len() - 1)
		stmn.WriteString("),")
	}
	stmn.Truncate(stmn.Len() - 1)
	_, err = o.Tx.ExecContext(ctx, stmn.String(), vArg...)
	return
}

func (o *OrderDetailRepo) GetAll(ctx context.Context, order_id uint, details *[]domain.OrderDetail) (err error) {
	var (
		rows *sql.Rows
		dts  []domain.OrderDetail
	)
	if rows, err = o.Tx.QueryContext(ctx, "select id, order_id,product_id, quantity from order_detail where order_id=$1", order_id); err != nil {
		return
	}

	for rows.Next() {
		var dt domain.OrderDetail
		if err = rows.Scan(&dt.ID, &dt.OrderID, &dt.ProductID, &dt.Quantity); err != nil {
			return
		}
		dts = append(dts, dt)
	}
	details = &dts
	return
}
