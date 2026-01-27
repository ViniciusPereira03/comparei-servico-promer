package user_interface

import "main/internal/domain/user"

type UserRepository interface {
	CreateUser(user *user.User) error
	EditUser(user *user.User) error
	GetUser(id string) (*user.User, error)
}
