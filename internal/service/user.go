package service

import (
    "github.com/YoubaImkf/go-auth-api/internal/model"
    "github.com/YoubaImkf/go-auth-api/internal/repository"
)

type UserService interface {
    GetAllUsers() ([]model.User, error)
}

type userService struct {
    userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
    return &userService{
        userRepository: userRepo,
    }
}

func (s *userService) GetAllUsers() ([]model.User, error) {
    return s.userRepository.GetAll()
}