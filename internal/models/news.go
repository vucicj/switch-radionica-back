package models

import (
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImagePath string    `json:"image_path"`
	Category  string    `json:"category"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
