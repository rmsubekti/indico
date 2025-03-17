package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
)

type UserRepo struct {
	Tx *sql.Tx
}

func (u *UserRepo) Add(ctx context.Context, in *domain.User) error {
	return u.Tx.QueryRowContext(ctx, `
		insert into user (name,email,pass,role) values ($1,$2,$3,$4) returning id 
	`, in.Name, in.Email, in.Password, in.UserRole).Scan(&in.ID)
}

func (u *UserRepo) GetByEmail(ctx context.Context, email string) (out domain.User, err error) {
	err = u.Tx.QueryRowContext(ctx, `
		select id, name, email, pass from user where email = $1 limit 1
	`, email).Scan(&out.ID, &out.Name, &out.Email, &out.Password, &out.UserRole)
	return
}

func (u *UserRepo) GetByID(ctx context.Context, uid uint) (out domain.User, err error) {
	err = u.Tx.QueryRowContext(ctx, `
		select id, name, email, pass from user where id = $1 limit 1
	`, uid).Scan(&out.ID, &out.Name, &out.Email, &out.Password, &out.UserRole)
	return
}

func (u *UserRepo) GetTotalRow(ctx context.Context, in *port.UserList) error {
	var stmn strings.Builder
	stmn.WriteString("select count(id) from user where 1=1 ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	return u.Tx.QueryRowContext(ctx, stmn.String()).Scan(&in.TotalRow)
}
func (u *UserRepo) List(ctx context.Context, in *port.UserList) (err error) {
	var (
		rows  *sql.Rows
		users []domain.User
		stmn  strings.Builder
	)
	stmn.WriteString("select id, name,email,role from user where 1=1 ")

	if len(in.Search) > 0 {
		stmn.WriteString(fmt.Sprintf(" and name ilike '%%%s%%'", in.Search))
	}
	stmn.WriteString(fmt.Sprintf("limit %d offset %d ", in.Limit, in.Offset))

	if rows, err = u.Tx.QueryContext(ctx, stmn.String()); err != nil {
		return
	}
	for rows.Next() {
		var user domain.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.UserRole); err != nil {
			return
		}
		users = append(users, user)
	}
	in.Rows = &users
	return
}
