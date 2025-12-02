package database

import (
	"time"
)

type PriorityLvl int

const (
	None PriorityLvl = iota
	Low
	Medium
	High
)

type Task struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Description    *string     `json:"description"`
	CreatedAt      time.Time   `json:"created_at"`
	NextReviewDate time.Time   `json:"next_review_date"`
	Priority       PriorityLvl `json:"priority"`
}

type SucceededTask struct {
	ID          int         `json:"id"`
	TaskID      int         `json:"task_id"`
	Title       string      `json:"title"`
	Description *string     `json:"description"`
	CompletedAt time.Time   `json:"completed_at"`
	Priority    PriorityLvl `json:"priority"`
}
