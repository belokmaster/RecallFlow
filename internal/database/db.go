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

		CREATE TABLE IF NOT EXISTS succeeded_task (
			id SERIAL PRIMARY KEY,
			task_id INTEGER NOT NULL,
			title VARCHAR(250) NOT NULL,
			description TEXT,
			completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	return err
}

func UpdateTaskNextReviewDate(db *sql.DB, id int, newDate string) error {
	query := "UPDATE task SET next_review_date = $1 WHERE id = $2"
	_, err := db.Exec(query, newDate, id)
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

func CompleteTask(db *sql.DB, taskID int) error {
	fmt.Printf("CompleteTask - Starting for task ID: %d\n", taskID)

	// Сначала получаем задачу
	task, err := GetTaskByID(db, taskID)
	if err != nil {
		fmt.Printf("CompleteTask - Error getting task: %v\n", err)
		return err
	}
	fmt.Printf("CompleteTask - Task to complete: ID=%d, Title=%s\n", task.ID, task.Title)

	// Вставляем в succeeded_task
	query := `
        INSERT INTO succeeded_task (task_id, title, description) 
        VALUES ($1, $2, $3)
    `
	fmt.Printf("CompleteTask - Inserting into succeeded_task...\n")
	result, err := db.Exec(query, task.ID, task.Title, task.Description)
	if err != nil {
		fmt.Printf("CompleteTask - Error inserting into succeeded_task: %v\n", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("CompleteTask - Inserted into succeeded_task, rows affected: %d\n", rowsAffected)

	// Удаляем из task
	fmt.Printf("CompleteTask - Deleting from task table...\n")
	result, err = db.Exec("DELETE FROM task WHERE id = $1", taskID)
	if err != nil {
		fmt.Printf("CompleteTask - Error deleting from task: %v\n", err)
		return err
	}

	rowsAffected, _ = result.RowsAffected()
	fmt.Printf("CompleteTask - Deleted from task, rows affected: %d\n", rowsAffected)

	fmt.Printf("CompleteTask - Successfully completed task %d\n", taskID)
	return nil
}

func GetSucceededTasks(db *sql.DB) ([]SucceededTask, error) {
	fmt.Printf("GetSucceededTasks - Fetching succeeded tasks\n")

	query := "SELECT id, task_id, title, description, completed_at FROM succeeded_task ORDER BY completed_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("GetSucceededTasks - Error querying: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []SucceededTask
	for rows.Next() {
		var task SucceededTask
		var description sql.NullString

		err := rows.Scan(&task.ID, &task.TaskID, &task.Title, &description, &task.CompletedAt)
		if err != nil {
			fmt.Printf("GetSucceededTasks - Error scanning row: %v\n", err)
			return nil, err
		}

		if description.Valid {
			task.Description = &description.String
		}

		tasks = append(tasks, task)
	}

	fmt.Printf("GetSucceededTasks - Found %d succeeded tasks\n", len(tasks))
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
