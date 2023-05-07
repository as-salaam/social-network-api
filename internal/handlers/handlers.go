package handlers

import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		DB: db,
	}
}
