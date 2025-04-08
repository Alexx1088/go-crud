package repositories

import (
	"database/sql"
	"errors"
	"go-crud/internal/models"
	"go-crud/internal/utils"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *UserRepository) CreateUser(user models.User) (int, error) {
	hashedPassword, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.DB.QueryRow("INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Email, hashedPassword).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) UpdateUser(id int, user models.User) error {
	query := `
        UPDATE users
        SET name = $1, email = $2, password_hash = $3
        WHERE id = $4
    `
	_, err := r.DB.Exec(query, user.Name, user.Email, user.PasswordHash, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	_, err := r.DB.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT id, name, email, password_hash FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}
