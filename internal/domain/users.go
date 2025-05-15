package domain

import "github.com/go-playground/validator/v10"

const UserIDKey = "userId"

type User struct {
	Id       int    `json:"-" db:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var Validate = validator.New()
