package services

import (
	"errors"
	"go-crud/internal/models"
	"go-crud/internal/repositories"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Create a mock repository
	mockRepo := repositories.NewMockUserRepositoryInterface(ctrl)

	// Create a service instance with the mock repository
	service := NewUserService(mockRepo)

	// Test data
	user := models.User{
		Name:         "John",
		Email:        "john@gmail.com",
		PasswordHash: "password",
	}
	expectedId := 1

	// Mock repository behavior
	mockRepo.EXPECT().CreateUser(user).Return(expectedId, nil)

	//Call the method
	result, err := service.CreateUser(user)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedId, result.ID)
	assert.Equal(t, user.Name, result.Name)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.PasswordHash, result.PasswordHash)
}

// Test 1: Success Case for GetAllUsers
func TestGetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Create a mock repository
	mockRepo := repositories.NewMockUserRepositoryInterface(ctrl)

	//Create a service instance with the mock repository
	service := NewUserService(mockRepo)

	//Mock data
	expectedUsers := []models.User{
		{ID: 1, Name: "John", Email: "john@gmail.com", PasswordHash: "password"},
		{ID: 2, Name: "Jane", Email: "jane@gmail.com", PasswordHash: "password"},
	}

	// Mock repository behavior
	mockRepo.EXPECT().GetAllUsers().Return(expectedUsers, nil)

	// Call the method
	result, err := service.GetAllUsers()

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, result)
}

// Test 2: Empty List for GetAllUsers
func TestGetAllUsers_EmptyList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := repositories.NewMockUserRepositoryInterface(ctrl)

	// Create a service instance with the mock repository
	service := NewUserService(mockRepo)

	// Mock repository behavior
	mockRepo.EXPECT().GetAllUsers().Return([]models.User{}, nil)

	// Call the method
	result, err := service.GetAllUsers()

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, result)
}

// Test 3: Repository Error for GetAllUsers
func TestGetAllUsers_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := repositories.NewMockUserRepositoryInterface(ctrl)

	// Create a service instance with the mock repository
	service := NewUserService(mockRepo)

	// Mock repository behavior
	mockRepo.EXPECT().GetAllUsers().Return(nil, errors.New("database error"))

	// Call the method
	result, err := service.GetAllUsers()

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
}
