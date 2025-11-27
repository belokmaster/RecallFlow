package database

import (
	"database/sql"
	"fmt"
	"log"
	"reccal_flow/internal/config"
	"time"
)

func InitDB(path string) (*sql.DB, error) {
	log.Println("InitDB: Start inizialization of DB")
	config, err := config.ReadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("problem with gettig config: %v", err)
	}
	connStr := config.ConnectionString()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("problem with connecting to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("problem with ping db: %v", err)
	}

	if err := CreateTables(db); err != nil {
		return nil, fmt.Errorf("problem with creating db: %v", err)
	}

	log.Println("InitDB: DB sucessfully inizializated")
	return db, nil
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

func UpdateTaskNextReviewDate(db *sql.DB, id int, newDate time.Time) error {
	query := "UPDATE task SET next_review_date = $1 WHERE id = $2"
	_, err := db.Exec(query, newDate, id)
	return err
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	query := "SELECT id, title, description, created_at, next_review_date FROM task ORDER BY next_review_date ASC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var description sql.NullString
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.CreatedAt,
			&task.NextReviewDate,
		)

		if err != nil {
			return nil, err
		}

		if description.Valid {
			task.Description = &description.String
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func CompleteTask(db *sql.DB, taskID int) error {
	log.Printf("CompleteTask: Starting for task ID: %d\n", taskID)
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback() // Откатит транзакцию, если она не была закоммичена

	task, err := GetTaskByID(db, taskID)
	if err != nil {
		return err
	}

	insertQuery := `
        INSERT INTO succeeded_task (task_id, title, description) 
        VALUES ($1, $2, $3)
    `

	if _, err := tx.Exec(insertQuery, task.ID, task.Title, task.Description); err != nil {
		return fmt.Errorf("failed to insert into succeeded_task: %v", err)
	}

	deleteQuery := "DELETE FROM task WHERE id = $1"
	if _, err := tx.Exec(deleteQuery, taskID); err != nil {
		return fmt.Errorf("failed to delete from task: %v", err)
	}

	return tx.Commit()
}

func GetSucceededTasks(db *sql.DB) ([]SucceededTask, error) {
	log.Println("GetSucceededTasks: Fetching succeeded tasks")

	query := "SELECT id, task_id, title, description, completed_at FROM succeeded_task ORDER BY completed_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("GetSucceededTasks: Error querying: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []SucceededTask
	for rows.Next() {
		var task SucceededTask
		var description sql.NullString

		err := rows.Scan(&task.ID,
			&task.TaskID,
			&task.Title,
			&description,
			&task.CompletedAt,
		)

		if err != nil {
			log.Printf("GetSucceededTasks: Error scanning row: %v\n", err)
			return nil, err
		}

		if description.Valid {
			task.Description = &description.String
		}

		tasks = append(tasks, task)
	}

	log.Printf("GetSucceededTasks: Found %d succeeded tasks\n", len(tasks))
	return tasks, nil
}

func GetTaskByID(db *sql.DB, id int) (*Task, error) {
	var task Task
	query := "SELECT id, title, description, created_at, next_review_date FROM task WHERE id = $1"

	err := db.QueryRow(
		query,
		id,
	).Scan(&task.ID,
		&task.Title,
		&task.Description,
		&task.CreatedAt,
		&task.NextReviewDate,
	)

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

func DeleteTask(db *sql.DB, id int) error {
	query := "DELETE FROM task WHERE id = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete task: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("task with id=%d not found", id)
	}

	return nil
}

func DeleteSucceededTask(db *sql.DB, id int) error {
	query := "DELETE FROM succeeded_task WHERE id = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete succeeded_task: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("succeeded_task with id=%d not found", id)
	}

	return nil
}

func RedactTask(db *sql.DB, task Task) error {
	query := "UPDATE task SET title = $1, description = $2, created_at = $3, next_review_date = $4 WHERE id = $5"
	res, err := db.Exec(query, task.Title, task.Description, task.CreatedAt, task.NextReviewDate, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("succeeded task with id=%d not found", task.ID)
	}

	return err
}

func RedactSucceededTask(db *sql.DB, task SucceededTask) error {
	query := "UPDATE succeeded_task SET title = $1, description = $2 WHERE id = $3"
	res, err := db.Exec(query, task.Title, task.Description, task.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("succeeded task with id=%d not found", task.ID)
	}

	return err
}
