package main

import (
	"log"
	"reccal_flow/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	path := "text.txt"
	err := database.InitDB(path)
	if err != nil {
		log.Fatalf("problem with init db: %v", err)
	}
	defer database.DB.Close()
}
