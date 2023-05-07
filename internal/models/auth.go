package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Credentials struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserID  uuid.UUID
	TokenID uuid.UUID
	jwt.RegisteredClaims
}

type Token struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	IsValid   bool      `gorm:"default:true"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
