package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reccal_flow/internal/database"
	"strconv"
	"time"
)

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

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("GetTasksHandler - Fetching all tasks\n")

	// Получаем активные задачи
	tasks, err := database.GetAllTasks(database.DB)
	if err != nil {
		fmt.Printf("GetTasksHandler - Error getting active tasks: %v\n", err)
		response := TasksResponse{Error: "Failed to get tasks: " + err.Error()}
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	fmt.Printf("GetTasksHandler - Found %d active tasks\n", len(tasks))

	// Получаем выполненные задачи
	succeededTasks, err := database.GetSucceededTasks(database.DB)
	if err != nil {
		fmt.Printf("GetTasksHandler - Error getting succeeded tasks: %v\n", err)
		succeededTasks = []database.SucceededTask{}
	}
	fmt.Printf("GetTasksHandler - Found %d succeeded tasks\n", len(succeededTasks))

	response := TasksResponse{
		Tasks:          tasks,
		SucceededTasks: succeededTasks,
	}

	sendJSONResponse(w, http.StatusOK, response)
}

func UpdateTaskDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем ID задачи из URL
	path := r.URL.Path
	taskIDStr := path[len("/tasks/"):]
	fmt.Printf("UpdateTaskDateHandler - Path: %s, TaskID: %s\n", path, taskIDStr)

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		fmt.Printf("Error converting task ID: %v\n", err)
		response := TaskActionResponse{Error: "Invalid task ID: " + taskIDStr}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var req UpdateTaskDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		response := TaskActionResponse{Error: "Invalid JSON format"}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	fmt.Printf("Updating task %d with date: %s\n", taskID, req.NewDate)

	// Валидация даты
	newDate, err := time.Parse("2006-01-02T15:04", req.NewDate)
	if err != nil {
		newDate, err = time.Parse("2006-01-02T15:04:05", req.NewDate)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			response := TaskActionResponse{Error: "Invalid date format. Expected: YYYY-MM-DDTHH:MM"}
			sendJSONResponse(w, http.StatusBadRequest, response)
			return
		}
	}

	// Используем тот же формат, что и при создании задачи
	formattedDate := newDate.Format("2006-01-02 15:04:05")
	fmt.Printf("Formatted date for DB: %s\n", formattedDate)

	err = database.UpdateTaskNextReviewDate(database.DB, taskID, formattedDate)
	if err != nil {
		fmt.Printf("Error updating task in DB: %v\n", err)
		response := TaskActionResponse{Error: "Failed to update task: " + err.Error()}
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	fmt.Printf("Task %d successfully updated\n", taskID)
	response := TaskActionResponse{Success: true}
	sendJSONResponse(w, http.StatusOK, response)
}

func CompleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	taskIDStr := path[len("/tasks/complete/"):]
	fmt.Printf("CompleteTaskHandler - Path: %s, TaskID: %s\n", path, taskIDStr)

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		fmt.Printf("Error converting task ID: %v\n", err)
		response := TaskActionResponse{Error: "Invalid task ID: " + taskIDStr}
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	fmt.Printf("Completing task %d\n", taskID)

	err = database.CompleteTask(database.DB, taskID)
	if err != nil {
		fmt.Printf("Error completing task: %v\n", err)
		response := TaskActionResponse{Error: "Failed to complete task: " + err.Error()}
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	fmt.Printf("Task %d successfully completed\n", taskID)
	response := TaskActionResponse{Success: true}
	sendJSONResponse(w, http.StatusOK, response)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
