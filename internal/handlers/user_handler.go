package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "fmt"
	"github.com/go-playground/validator/v10"
	"go-crud/internal/repositories"
	"go-crud/middleware"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"go-crud/internal/models"
	"go-crud/internal/services"
)

// validateUser validates the User struct and returns an error message if validation fails.
func validateUser(user models.User) (string, bool) {
	if err := models.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		errorMessage := "Validation failed:"
		for _, e := range errors {
			errorMessage += fmt.Sprintf(" Field=%s, Tag=%s, Value=%v", e.Field(), e.Tag(), e.Value())
		}
		log.Println(errorMessage)
		return errorMessage, false
	}
	return "", true
}

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func RegisterUserRoutes(router *mux.Router, db *sql.DB, secretKey []byte) {
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)
	handler := NewUserHandler(service)

	// Apply AuthMiddleware to all /users routes
	protectedRouter := router.PathPrefix("/users").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(secretKey)) // Protect all /users routes

	protectedRouter.HandleFunc("", handler.GetUsers).Methods("GET")
	protectedRouter.HandleFunc("/{id}", handler.GetUser).Methods("GET")
	protectedRouter.HandleFunc("", handler.CreateUser).Methods("POST")
	protectedRouter.HandleFunc("/{id}", handler.UpdateUser).Methods("PUT")
	protectedRouter.HandleFunc("/{id}", handler.DeleteUser).Methods("DELETE")
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	user, err := h.Service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the User struct
	if errorMessage, isValid := validateUser(user); !isValid {
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	newUser, err := h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updateUserReq models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&updateUserReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validatePartialUpdate(updateUserReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.Service.UpdateUser(id, updateUserReq); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	if err := h.Service.DeleteUser(id); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully!"})
}

// validatePartialUpdate validates the fields provided in the UpdateUserRequest.
func validatePartialUpdate(req models.UpdateUserRequest) error {
	if req.Name != nil && len(*req.Name) < 2 {
		return fmt.Errorf("validation failed: Name must be at least 2 characters")
	}
	if req.Email != nil && !isValidEmail(*req.Email) {
		return fmt.Errorf("validation failed: Invalid email format")
	}
	if req.PasswordHash != nil && len(*req.PasswordHash) < 6 {
		return fmt.Errorf("validation failed: Password must be at least 6 characters")
	}
	return nil
}

// isValidEmail checks if an email address is valid.
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
