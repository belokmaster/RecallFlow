package database

import (
	"database/sql"
	"fmt"
	"log"
	"reccal_flow/internal/config"
)

func InitDB(path string) error {
	config, err := config.ReadConfig(path)
	if err != nil {
		return fmt.Errorf("problem with gettig config: %v", err)
	}
	connStr := config.ConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("problem with connecting to db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("problem with ping db: %v", err)
	}

	err = CreateTables(db)
	if err != nil {
		return fmt.Errorf("problem with creating db: %v", err)
	}

	DB = db // глобальная переменная конечно не круто но че поделать
	log.Println("DB sucessfully inizializated")
	return nil
}

func CreateTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS task (
			id SERIAL PRIMARY KEY,
			title VARCHAR(250) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			next_review_date TIMESTAMP NOT NULL
		);
	`

	_, err := db.Exec(query)
	return err
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	query := "SELECT id, title, description, created_at, next_review_date FROM task"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.NextReviewDate)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskByID(db *sql.DB, id int) (*Task, error) {
	var task Task
	query := "SELECT id, title, description, created_at, next_review_date FROM task WHERE id = $1"

	err := db.QueryRow(
		query,
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.NextReviewDate)

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// просроченные таски
func GetOverdueTasks(db *sql.DB) ([]Task, error) {
	query := "SELECT id, title, description, created_at, next_review_date FROM task WHERE next_review_date < NOW() ORDER BY next_review_date ASC"
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.CreatedAt, &task.NextReviewDate)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func AddNewTask(db *sql.DB, task Task) (*Task, error) {
	query := `
		INSERT INTO task (title, description, next_review_date) 
		VALUES ($1, $2, $3) 
		RETURNING id, title, description, created_at, next_review_date
	`

	var createdTask Task
	var description sql.NullString

	err := db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.NextReviewDate,
	).Scan(&createdTask.ID, &createdTask.Title, &description, &createdTask.CreatedAt, &createdTask.NextReviewDate)

	if err != nil {
		return nil, fmt.Errorf("problem with creating task: %v", err)
	}

	// Конвертируем sql.NullString в *string
	if description.Valid {
		createdTask.Description = &description.String
	}

	return &createdTask, nil
}
