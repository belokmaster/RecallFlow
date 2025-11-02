package handlers

import (
	"encoding/json"
	"net/http"
	"reccal_flow/internal/database"
	"time"
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

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := CreateTaskResponse{Error: "Invalid JSON format"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация
	if req.Title == "" {
		response := CreateTaskResponse{Error: "Title is required"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(req.Title) > 250 {
		response := CreateTaskResponse{Error: "Title too long (max 250 characters)"}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Парсим дату
	nextReviewDate, err := time.Parse("2006-01-02T15:04", req.NextReviewDate)
	if err != nil {
		nextReviewDate, err = time.Parse("2006-01-02T15:04:05", req.NextReviewDate)
		if err != nil {
			response := CreateTaskResponse{Error: "Invalid date format. Expected: YYYY-MM-DDTHH:MM. Got: " + req.NextReviewDate}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	// Создаем задачу
	task := database.Task{
		Title:          req.Title,
		Description:    req.Description,
		NextReviewDate: nextReviewDate,
	}

	createdTask, err := database.AddNewTask(database.DB, task)
	if err != nil {
		response := CreateTaskResponse{Error: "Failed to create task: " + err.Error()}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CreateTaskResponse{Task: createdTask}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
