package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"strconv"
	"time"
)

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

		task := database.Task{
			Title:          req.Title,
			Description:    req.Description,
			NextReviewDate: nextReviewDate,
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

		idStr := r.URL.Query().Get("id")
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
