package services

import (
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
