package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type ProductRepo struct {
	Tx *sql.Tx
}

func (p *ProductRepo) Add(ctx context.Context, in *domain.Product) error {
	return p.Tx.QueryRowContext(ctx, `
		insert into product (name,sku,qty) values ($1,$2,$3) returning id 
	`, in.Name, in.SKU, in.Quantity).Scan(&in.ID)
}

func (p *ProductRepo) Get(ctx context.Context, productId uint) (out domain.Product, err error) {
	err = p.Tx.QueryRowContext(ctx, `
		select id, name, sku, qty, created_at, updated_at, deleted_at from product where id = $1 limit 1
	`, productId).Scan(&out.ID, &out.Name, &out.SKU, &out.Quantity, &out.CreatedAt, &out.UpdatedAt, &out.DeletedAt)
	return
}

func (p *ProductRepo) GetTotalRow(ctx context.Context, in *port.ProductList) error {
	var stmn strings.Builder
	stmn.WriteString("select count(id) from product where deleted_at is null ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	return p.Tx.QueryRowContext(ctx, stmn.String()).Scan(&in.TotalRow)
}

func (p *ProductRepo) List(ctx context.Context, in *port.ProductList) (err error) {
	var (
		rows     *sql.Rows
		products []domain.Product
		stmn     strings.Builder
	)
	stmn.WriteString("select id, name,sku from product where deleted_at is null ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	stmn.WriteString(fmt.Sprintf("limit %d offset %d ", in.Limit, in.Offset))

	if rows, err = p.Tx.QueryContext(ctx, stmn.String()); err != nil {
		return
	}
	for rows.Next() {
		var product domain.Product
		if err = rows.Scan(&product.ID, &product.Name, &product.SKU); err != nil {
			return
		}
		products = append(products, product)
	}
	in.Rows = &products
	return
}

func (p *ProductRepo) Update(ctx context.Context, in *domain.Product) (err error) {
	err = p.Tx.QueryRowContext(ctx, `
	update product set name=$1, sku=$2, qty=$3, updated_at=$4 where id=$5 
	`, in.Name, in.SKU, in.Quantity, in.UpdatedAt, in.ID).Err()
	return
}

func (p *ProductRepo) Delete(ctx context.Context, productId uint) error {
	return p.Tx.QueryRowContext(ctx, `
	update product set updated_at=now() where id=$1 
	`, productId).Err()
}
