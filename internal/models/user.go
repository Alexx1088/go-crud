package models

import "github.com/go-playground/validator/v10"

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name" validate:"required,min=2,max=20"`
	Email        string `json:"email" validate:"required,email"`
	PasswordHash string `json:"passwordHash" validate:"required,min=6"`
}

// Validate Global validator instance
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}
