package database

import (
	"time"
)

type Card struct {
	ID           int        `json:"id"`
	Word         string     `json:"word"`
	Translation  string     `json:"translation"`
	Example      *string    `json:"example"`
	CreatedAt    time.Time  `json:"created_at"`
	LastReviewed *time.Time `json:"last_reviewed"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Attempts     int        `json:"attempts"`
	Successes    int        `json:"successes"`
}
