package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reccal_flow/internal/database"
)

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

type EditTaskRequest struct {
	Title          string  `json:"title"`
	Description    *string `json:"description"`
	NewCreatedAt   string  `json:"created_at"`
	NextReviewDate string  `json:"next_review_date"`
}

type EditSucceededTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Printf("Отправка ошибки: код %d, сообщение: %s", code, message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Критическая ошибка: не удалось закодировать JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
