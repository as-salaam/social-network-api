package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Profile  *Profile  `json:"profile,omitempty"`
	Posts    []Post    `json:"posts,omitempty"`
	Comments []Comment `json:"comments,omitempty"`
	Stories  []Story   `json:"stories,omitempty"`
}
