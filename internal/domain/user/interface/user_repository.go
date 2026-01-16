package user_interface

import "main/internal/domain/user"

type UserRepository interface {
	CreateUser(user *user.User) error
}
