package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"strconv"
	"time"
)

func getPriorityLvl(priorityStr string) (database.PriorityLvl, error) {
	var priority database.PriorityLvl

	switch priorityStr {
	case "":
		priority = database.None
	case "Low":
		priority = database.Low
	case "Medium":
		priority = database.Medium
	case "High":
		priority = database.High
	default:
		return 0, fmt.Errorf("неверный приоритет приоретатов: %s", priorityStr)
	}

	return priority, nil
}

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateTaskHandler: Начало обработки")

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		var req CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		log.Printf("CreateTaskHandler: Получены данные: %+v", req)

		if req.Title == "" || len(req.Title) > 250 {
			respondWithError(w, http.StatusBadRequest, "Заголовок обязателен и не должен превышать 250 символов")
			return
		}

		nextReviewDate, err := time.Parse("2006-01-02T15:04:05", req.NextReviewDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат даты. Ожидается: YYYY-MM-DDTHH:MM:SS")
			return
		}

		priority, err := getPriorityLvl(req.Priority)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный уровень приоритета. Допустимые значения: Low, Medium, High")
			return
		}

		task := database.Task{
			Title:          req.Title,
			Description:    req.Description,
			NextReviewDate: nextReviewDate,
			Priority:       priority,
		}

		createdTask, err := database.AddNewTask(db, task)
		if err != nil {
			log.Printf("CreateTaskHandler: Ошибка при добавлении в БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Не удалось создать задачу в базе данных")
			return
		}

		log.Printf("CreateTaskHandler: Задача успешно создана с ID: %d", createdTask.ID)
		respondWithJSON(w, http.StatusCreated, createdTask)
	}
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetTasksHandler: Начало обработки")

		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		tasks, err := database.GetAllTasks(db)
		if err != nil {
			log.Printf("GetTasksHandler: Ошибка получения активных задач: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Не удалось получить список задач")
			return
		}

		succeededTasks, err := database.GetSucceededTasks(db)
		if err != nil {
			log.Printf("GetTasksHandler: Ошибка получения выполненных задач: %v", err)
			// Не считаем это критической ошибкой, просто вернем пустой список
			succeededTasks = []database.SucceededTask{}
		}

		response := TasksResponse{
			Tasks:          tasks,
			SucceededTasks: succeededTasks,
		}

		log.Printf("GetTasksHandler: Отправка %d активных и %d выполненных задач", len(tasks), len(succeededTasks))
		respondWithJSON(w, http.StatusOK, response)
	}
}

func UpdateTaskDateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateTaskDateHandler: Начало обработки")

		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")

		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный или отсутствующий ID задачи в URL")
			return
		}

		var req UpdateTaskDateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		log.Printf("UpdateTaskDateHandler: Обновление даты для задачи ID %d на %s", taskID, req.NewDate)

		newDate, err := time.Parse("2006-01-02T15:04:05", req.NewDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат даты. Ожидается: YYYY-MM-DDTHH:MM:SS")
			return
		}

		if err := database.UpdateTaskNextReviewDate(db, taskID, newDate); err != nil {
			log.Printf("UpdateTaskDateHandler: Ошибка обновления в БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Не удалось обновить задачу")
			return
		}

		log.Printf("UpdateTaskDateHandler: Задача %d успешно обновлена", taskID)
		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func CompleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CompleteTaskHandler: Начало обработки")

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный или отсутствующий ID задачи в URL")
			return
		}

		log.Printf("CompleteTaskHandler: Завершение задачи с ID: %d", taskID)

		if err := database.CompleteTask(db, taskID); err != nil {
			log.Printf("CompleteTaskHandler: Ошибка завершения задачи в БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Не удалось завершить задачу")
			return
		}

		log.Printf("CompleteTaskHandler: Задача %d успешно завершена", taskID)
		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID")
			return
		}

		if err := database.DeleteTask(db, taskID); err != nil {
			log.Printf("DeleteTaskHandler: Ошибка удаления: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка удаления задачи")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func DeleteSucceededTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID")
			return
		}

		if err := database.DeleteSucceededTask(db, taskID); err != nil {
			log.Printf("DeleteSucceededTaskHandler: Ошибка удаления: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка удаления выполненной задачи")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func EditTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("EditTaskHandler: Начало обработки")

		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID задачи")
			return
		}

		var req EditTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		if req.Title == "" || len(req.Title) > 250 {
			respondWithError(w, http.StatusBadRequest, "Заголовок обязателен")
			return
		}

		newCreatedAt, err := time.Parse("2006-01-02T15:04:05", req.NewCreatedAt)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}

		nextReviewDate, err := time.Parse("2006-01-02T15:04:05", req.NextReviewDate)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат даты")
			return
		}

		priority, err := getPriorityLvl(req.Priority)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный уровень приоритета. Допустимые значения: Low, Medium, High")
			return
		}

		updatedTask := database.Task{
			ID:             taskID,
			Title:          req.Title,
			Description:    req.Description,
			CreatedAt:      newCreatedAt,
			NextReviewDate: nextReviewDate,
			Priority:       priority,
		}

		if err := database.RedactTask(db, updatedTask); err != nil {
			log.Printf("EditTaskHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func EditSucceededTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("EditSucceededTaskHandler: Начало обработки")

		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID задачи")
			return
		}

		var req EditSucceededTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		if req.Title == "" || len(req.Title) > 250 {
			respondWithError(w, http.StatusBadRequest, "Заголовок обязателен")
			return
		}

		priority, err := getPriorityLvl(req.Priority)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный уровень приоритета. Допустимые значения: Low, Medium, High")
			return
		}

		updatedSucceededTask := database.SucceededTask{
			ID:          taskID,
			Title:       req.Title,
			Description: req.Description,
			Priority:    priority,
		}

		if err := database.RedactSucceededTask(db, updatedSucceededTask); err != nil {
			log.Printf("EditSucceededTask: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении задачи")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}
