package main

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MyCustomClaims struct {
	*jwt.StandardClaims
	ID uuid.UUID `json:"id"`
	// Email string    `json:"email"`
}

// ? Keeping user email and otp in a block level struct
type OTPVerification struct {
	Email string
	OTP   string
}

type User struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name  string
	Email string
}
