package services

import (
	_ "errors"
	"go-crud/internal/models"
	"go-crud/internal/repositories"
	"go-crud/internal/utils"
)

type NewMockUserServiceInterface interface {
}

type UserService struct {
	Repo repositories.UserRepositoryInterface
}

func NewUserService(repo repositories.UserRepositoryInterface) *UserService {
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

func (s *UserService) UpdateUser(id int, req models.UpdateUserRequest) error {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.PasswordHash != nil {
		hashedPassword, err := utils.HashPassword(*req.PasswordHash)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword
	}

	return s.Repo.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id int) error {
	return s.Repo.DeleteUser(id)
}

func (s *UserService) GetUserByEmail(email string) (models.User, error) {
	return s.Repo.GetUserByEmail(email)
}
