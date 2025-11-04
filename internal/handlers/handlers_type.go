package handlers

import "reccal_flow/internal/database"

type CreateTaskRequest struct {
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	NextReviewDate string  `json:"next_review_date"`
}

type CreateTaskResponse struct {
	Task  *database.Task `json:"task,omitempty"`
	Error string         `json:"error,omitempty"`
}

type TasksResponse struct {
	Tasks          []database.Task          `json:"tasks"`
	SucceededTasks []database.SucceededTask `json:"succeeded_tasks,omitempty"`
	Error          string                   `json:"error,omitempty"`
}

type UpdateTaskDateRequest struct {
	NewDate string `json:"new_date"`
}

type TaskActionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
