package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
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

func (p *Profile) AfterFind(tx *gorm.DB) error {
	p.Avatar = os.Getenv("APP_URL") + p.Avatar
	return nil
}
