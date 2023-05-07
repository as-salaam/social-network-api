package models

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id"`
	Avatar    string    `json:"avatar"`
	Email     string    `json:"email"`
	Type      string    `json:"type"`
	Bio       string    `json:"bio"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
