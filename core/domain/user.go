package domain

import (
	"errors"

	"github.com/pyfal/gook"
)

type User struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"`
	UserRole UserRole `json:"user_role"`
}

type UserRole string

const UserAdmin UserRole = "Admin"
const UserStaff UserRole = "Staff"

func (u *User) Valid() (err error) {
	if len(u.Name) < 2 {
		return errors.New("nama harus lebih dari 2 karakter")
	}
	if err = gook.Email(u.Email).Error(); err != nil {
		return
	}
	if err = gook.Password(u.Password).Error(); err != nil {
		return
	}
	return
}
