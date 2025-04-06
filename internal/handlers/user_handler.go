package handlers

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"go-crud/internal/repositories"
	"go-crud/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go-crud/internal/models"
	"go-crud/internal/services"
)

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
	newUser, err := h.Service.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateUser(id, user); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	user.ID = id
	json.NewEncoder(w).Encode(user)
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
