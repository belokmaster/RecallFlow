package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

func CreateTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS task (
			id SERIAL PRIMARY KEY,
			title VARCHAR(250) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			next_review_date TIMESTAMP NOT NULL
		);
	`

	_, err := db.Exec(query)
	return err
}

func InitDB(path string) error {
	config, err := ReadConfig(path)
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

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %v", filename, err)
	}
	defer file.Close()

	config := &Config{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// пустрая строка и комменты пропускаются
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// парсим на ключ и value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "host":
			config.Host = value
		case "port":
			config.Port = value
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "dbname":
			config.DBName = value
		case "sslmode":
			config.SSLMode = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("problem with reading file: %v", err)
	}

	return config, nil
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}
