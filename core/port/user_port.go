package port

import (
	"context"

	"github.com/rmsubekti/indico/core/domain"
)

type (
	UserLogin struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	UserRegister struct {
		domain.User
		ConfirmPassword string `json:"confirm_password"`
	}
	UserList struct {
		domain.Pagination
		Rows *[]domain.User `json:"rows,omitempty"`
	}

	IUserRepo interface {
		Add(ctx context.Context, in *domain.User) error
		GetByEmail(ctx context.Context, email string) (domain.User, error)
		GetByID(ctx context.Context, uid uint) (domain.User, error)
		GetTotalRow(ctx context.Context, userList *UserList) error
		List(ctx context.Context, userList *UserList) error
	}

	IUserService interface {
		Register(ctx context.Context, in *UserRegister) error
		Login(ctx context.Context, login UserLogin) (domain.User, error)
		Get(ctx context.Context, user *domain.User) error
		List(ctx context.Context, userList *UserList) error
	}
)
