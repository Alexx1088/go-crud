package services

import (
	_ "errors"
	"go-crud/internal/models"
	"go-crud/internal/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAllUsers()
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) CreateUser(user models.User) (models.User, error) {
	id, err := s.Repo.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}
	user.ID = id
	return user, nil
}

func (s *UserService) UpdateUser(id int, user models.User) error {
	return s.Repo.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.Repo.DeleteUser(id)
}
