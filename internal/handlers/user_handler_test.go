package handlers

import (
	"github.com/golang/mock/gomock"
	"testing"

	"go-crud/internal/services"
)

func TestGetAllUsersHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Create a mock service
	mockService := services.NewMockUserServiceInterface(ctrl)
}
