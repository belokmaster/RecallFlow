package main

import (
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"reccal_flow/internal/handlers"

	_ "github.com/lib/pq"
)

func main() {
	path := "text.txt"

	db, err := database.InitDB(path)
	if err != nil {
		log.Fatalf("Ошибка при инициализации БД: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", serveHTML)
	mux.HandleFunc("POST /tasks", handlers.CreateTaskHandler(db))
	mux.HandleFunc("GET /tasks", handlers.GetTasksHandler(db))
	mux.HandleFunc("PUT /tasks/{id}", handlers.EditTaskHandler(db))
	mux.HandleFunc("POST /tasks/{id}/complete", handlers.CompleteTaskHandler(db))
	mux.HandleFunc("DELETE /tasks/{id}", handlers.DeleteTaskHandler(db))
	mux.HandleFunc("DELETE /tasks/succeeded/{id}", handlers.DeleteSucceededTaskHandler(db))
	mux.HandleFunc("PUT /tasks/succeeded/{id}", handlers.EditSucceededTaskHandler(db))

	log.Println("Сервер запущен на http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

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
