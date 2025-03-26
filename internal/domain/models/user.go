package models

import (
	"time"
)

type User struct {
	ID                string     `json:"id" db:"id"`
	Name              string     `json:"name" db:"name"`
	Email             string     `json:"email" db:"email"`
	Password          string     `json:"-" db:"password"`
	ResetToken        *string    `json:"-" db:"reset_token"`
	ResetTokenExpires *time.Time `json:"-" db:"reset_token_expires"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}
