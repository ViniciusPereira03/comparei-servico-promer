package app

import (
	"main/internal/domain/user"
	user_interface "main/internal/domain/user/interface"
)

type UserService struct {
	mysqlRepo user_interface.UserRepository
}

func NewUserService(mysqlRepo user_interface.UserRepository) *UserService {
	return &UserService{mysqlRepo: mysqlRepo}
}

func (s *UserService) CreateUser(user *user.User) error {
	err := s.mysqlRepo.CreateUser(user)
	return err
}

func (s *UserService) EditUser(user *user.User) error {
	err := s.mysqlRepo.EditUser(user)
	return err
}

func (s *UserService) GetUser(id string) (*user.User, error) {
	return s.mysqlRepo.GetUser(id)
}
