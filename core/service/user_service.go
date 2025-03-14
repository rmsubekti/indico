package service

import (
	"context"
	"errors"
	"math"
	"strings"

	"github.com/pyfal/gook"
	"github.com/rmsubekti/indico/core/domain"
	"github.com/rmsubekti/indico/core/port"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo port.IUserRepo
}

func NewUserService(repo port.IUserRepo) port.IUserService {
	return &userService{repo: repo}
}

func (u *userService) Get(ctx context.Context, user *domain.User) (err error) {
	if len(user.Email) > 1 {
		if *user, err = u.repo.GetByEmail(ctx, user.Email); err != nil {
			return
		}
	} else {
		if *user, err = u.repo.GetByID(ctx, user.ID); err != nil {
			return
		}
	}
	return
}

func (u *userService) List(ctx context.Context, in *port.UserList) (err error) {
	if in.Limit < 1 {
		in.Limit = 10
	}
	if in.Page < 1 {
		in.Page = 1
	}

	if err = u.repo.GetTotalRow(ctx, in); err != nil {
		return
	}

	in.Offset = (in.Page - 1) * in.Limit
	in.TotalPage = uint(math.Ceil(float64(in.TotalRow) / float64(in.Limit)))

	if err = u.repo.List(ctx, in); err != nil {
		return
	}
	return
}

func (u *userService) Register(ctx context.Context, in *port.UserRegister) (err error) {
	var hash []byte
	if !strings.EqualFold(in.Password, in.ConfirmPassword) {
		return errors.New("konfirmasi password salah")
	}

	if err = in.Valid(); err != nil {
		return
	}

	if hash, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost); err != nil {
		return
	}

	in.User.Password = string(hash)

	if err = u.repo.Add(ctx, &in.User); err != nil {
		return
	}
	return
}

func (u *userService) Login(ctx context.Context, login port.UserLogin) (out domain.User, err error) {
	if err = gook.Email(login.Email).Error(); err != nil {
		return
	}

	if out, err = u.repo.GetByEmail(ctx, login.Email); err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(out.Password), []byte(login.Password)); err != nil {
		return out, bcrypt.ErrMismatchedHashAndPassword
	}

	out.Password = ""
	return
}
