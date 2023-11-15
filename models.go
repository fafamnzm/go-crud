package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID    uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name  string
	Email string
}
