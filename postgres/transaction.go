package postgre

import (
	"database/sql"

	"github.com/rmsubekti/indico/core/port"
	"github.com/rmsubekti/indico/postgres/repo"
)

type (
	transaction struct {
		tx *sql.Tx
	}

	ITransaction interface {
		Rollback() error
		Commit() error
		UserRepo() port.IUserRepo
		WarehouseRepo() port.IWarehouseRepo
		ProductRepo() port.IProductRepo
	}
)

// Commit implements ITransaction.
func (tx *transaction) Commit() error {
	return tx.tx.Commit()
}

// Rollback implements ITransaction.
func (tx *transaction) Rollback() error {
	return tx.tx.Rollback()
}

func (tx transaction) UserRepo() port.IUserRepo {
	return &repo.UserRepo{Tx: tx.tx}
}

func (tx transaction) WarehouseRepo() port.IWarehouseRepo {
	return &repo.WarehouseRepo{Tx: tx.tx}
}

func (tx transaction) ProductRepo() port.IProductRepo {
	return &repo.ProductRepo{Tx: tx.tx}
}
