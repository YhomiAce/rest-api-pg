package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(payload *CreateUserPayload) (error)
}

type User struct {
	ID       int
	Username string
	Email    string 
	Password string
	CreatedAt time.Time 
	UpdatedAt time.Time
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateUserPayload struct {
	Username string 
	Email    string 
	Password string
}