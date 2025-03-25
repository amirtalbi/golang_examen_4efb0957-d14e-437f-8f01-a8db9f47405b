package service

import (
	"github.com/amirtalbi/examen_go/internal/domain/models"
	"github.com/amirtalbi/examen_go/internal/domain/repositories"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
