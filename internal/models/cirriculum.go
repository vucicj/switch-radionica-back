package models

import (
	"time"

	"github.com/google/uuid"
)

type Cirriculum struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Week      int       `json:"week"`
	Content   string    `json:"description"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
