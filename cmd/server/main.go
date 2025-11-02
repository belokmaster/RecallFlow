package main

import (
	"fmt"
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"reccal_flow/internal/handlers"

	_ "github.com/lib/pq"
)

func main() {
	path := "text.txt"

	err := database.InitDB(path)
	if err != nil {
		log.Fatalf("problem with init db: %v", err)
	}
	defer database.DB.Close()

	v, _ := database.GetAllTasks(database.DB)
	fmt.Println(v)

	// Настройка маршрутов ДО запуска сервера
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/tasks", handlers.CreateTaskHandler)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// Обработчик для главной страницы
func serveHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "web/templates/index.html")
}
