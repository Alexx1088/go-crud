package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go-crud/internal/models"
	"go-crud/internal/repositories"
	"go-crud/internal/services"
	"go-crud/internal/utils"
	"net/http"
)

type AuthHandler struct {
	Service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{
		Service: service,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(user)
	_, err := h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user login and generates a JWT token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credential struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Service.GetUserByEmail(credential.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if err := utils.VerifyPassword(user.PasswordHash, credential.Password); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// RegisterAuthRoutes registers authentication-related routes.
func RegisterAuthRoutes(router *mux.Router, db *sql.DB) {
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)
	handler := NewAuthHandler(service)

	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/login", handler.Login).Methods("POST")
}
