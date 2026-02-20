package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func GetDatabasePath() string {
	if path := os.Getenv("RECALL_FLOW_DB"); path != "" {
		return path
	}

	return "./recall_flow.db"
}

func InitDB() (*sql.DB, error) {
	log.Println("InitDB: Starting database initialization")
	
	dbPath := GetDatabasePath()
	log.Printf("InitDB: Using database path: %s", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("mistake to connection to DB: %v", err)
	}

	// on foreign keys for SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("mistake to enable foreign keys: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("mistake to ping DB: %v", err)
	}

	if err := CreateTables(db); err != nil {
		return nil, fmt.Errorf("mistake to create tables: %v", err)
	}

	log.Println("InitDB: DB initialization completed successfully")
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS card (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word TEXT NOT NULL,
			translation TEXT NOT NULL,
			example TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			last_reviewed DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			attempts INTEGER DEFAULT 0,
			successes INTEGER DEFAULT 0
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	checkColumn := `PRAGMA table_info(card)`
	rows, err := db.Query(checkColumn)
	if err != nil {
		return err
	}
	defer rows.Close()

	hasUpdatedAt := false
	for rows.Next() {
		var cid int
		var name string
		var colType string
		var notnull int
		var dfltValue interface{}
		var pk int

		err := rows.Scan(&cid, &name, &colType, &notnull, &dfltValue, &pk)
		if err != nil {
			continue
		}

		if name == "updated_at" {
			hasUpdatedAt = true
			break
		}
	}
	if !hasUpdatedAt {
		altQuery := `ALTER TABLE card ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP`
		_, err := db.Exec(altQuery)
		if err != nil {
			return fmt.Errorf("failed to add updated_at column: %v", err)
		}
	}

	return nil
}

// Card operations
func CreateCard(db *sql.DB, word, translation string, example *string) (*Card, error) {
	query := `
		INSERT INTO card (word, translation, example, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, word, translation, example, created_at, last_reviewed, updated_at, attempts, successes
	`

	var card Card
	var exampleNull sql.NullString
	var lastReviewedNull sql.NullTime

	err := db.QueryRow(query, word, translation, example).Scan(
		&card.ID,
		&card.Word,
		&card.Translation,
		&exampleNull,
		&card.CreatedAt,
		&lastReviewedNull,
		&card.UpdatedAt,
		&card.Attempts,
		&card.Successes,
	)

	if err != nil {
		return nil, fmt.Errorf("create card: %v", err)
	}

	if exampleNull.Valid {
		card.Example = &exampleNull.String
	}
	if lastReviewedNull.Valid {
		card.LastReviewed = &lastReviewedNull.Time
	}

	return &card, nil
}

func GetAllCards(db *sql.DB) ([]Card, error) {
	query := `
		SELECT id, word, translation, example, created_at, last_reviewed, updated_at, attempts, successes
		FROM card
		ORDER BY updated_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query cards: %v", err)
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		var exampleNull sql.NullString
		var lastReviewedNull sql.NullTime

		err := rows.Scan(
			&card.ID,
			&card.Word,
			&card.Translation,
			&exampleNull,
			&card.CreatedAt,
			&lastReviewedNull,
			&card.UpdatedAt,
			&card.Attempts,
			&card.Successes,
		)

		if err != nil {
			return nil, fmt.Errorf("scan card: %v", err)
		}

		if exampleNull.Valid {
			card.Example = &exampleNull.String
		}

		if lastReviewedNull.Valid {
			card.LastReviewed = &lastReviewedNull.Time
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func GetCardByID(db *sql.DB, id int) (*Card, error) {
	query := `
		SELECT id, word, translation, example, created_at, last_reviewed, updated_at, attempts, successes
		FROM card
		WHERE id = $1
	`

	var card Card
	var exampleNull sql.NullString
	var lastReviewedNull sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&card.ID,
		&card.Word,
		&card.Translation,
		&exampleNull,
		&card.CreatedAt,
		&lastReviewedNull,
		&card.UpdatedAt,
		&card.Attempts,
		&card.Successes,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("card with id=%d not found", id)
		}

		return nil, fmt.Errorf("query card: %v", err)
	}

	if exampleNull.Valid {
		card.Example = &exampleNull.String
	}

	if lastReviewedNull.Valid {
		card.LastReviewed = &lastReviewedNull.Time
	}

	return &card, nil
}

func DeleteCard(db *sql.DB, id int) error {
	query := "DELETE FROM card WHERE id = $1"
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete card: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("card with id=%d not found", id)
	}

	return nil
}

func UpdateCard(db *sql.DB, id int, word, translation string, example *string) error {
	query := `
		UPDATE card
		SET word = $1, translation = $2, example = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	res, err := db.Exec(query, word, translation, example, id)
	if err != nil {
		return fmt.Errorf("update card: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("card with id=%d not found", id)
	}

	return nil
}

// GetCardsForReview returns cards sorted by last_reviewed (oldest first) for spaced repetition
func GetCardsForReview(db *sql.DB) ([]Card, error) {
	query := `
		SELECT id, word, translation, example, created_at, last_reviewed, updated_at, attempts, successes
		FROM card
		ORDER BY 
			CASE WHEN last_reviewed IS NULL THEN 0 ELSE 1 END,
			last_reviewed ASC
		LIMIT 10
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query cards for review: %v", err)
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		var exampleNull sql.NullString
		var lastReviewedNull sql.NullTime

		err := rows.Scan(
			&card.ID,
			&card.Word,
			&card.Translation,
			&exampleNull,
			&card.CreatedAt,
			&lastReviewedNull,
			&card.UpdatedAt,
			&card.Attempts,
			&card.Successes,
		)

		if err != nil {
			return nil, fmt.Errorf("scan card: %v", err)
		}

		if exampleNull.Valid {
			card.Example = &exampleNull.String
		}

		if lastReviewedNull.Valid {
			card.LastReviewed = &lastReviewedNull.Time
		}

		cards = append(cards, card)
	}

	return cards, nil
}

// RepeatCard marks a card as reviewed and updates its statistics
func RepeatCard(db *sql.DB, id int, success bool) error {
	query := `
		UPDATE card
		SET last_reviewed = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP, attempts = attempts + 1
	`

	if success {
		query += `, successes = successes + 1`
	}

	query += ` WHERE id = $1`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("repeat card: %v", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}

	if rows == 0 {
		return fmt.Errorf("card with id=%d not found", id)
	}

	return nil
}
