package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"os"
	"time"
)

type File struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PostID    uuid.UUID `json:"post_id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *File) AfterFind(tx *gorm.DB) error {
	f.Path = os.Getenv("APP_URL") + f.Path
	return nil
}
