package database

import (
	"database/sql"
	"time"
)

type Task struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    *string   `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	NextReviewDate time.Time `json:"next_review_date"`
}

var DB *sql.DB
