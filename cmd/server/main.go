package main

import (
	"fmt"
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"reccal_flow/internal/handlers"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	path := "text.txt"

	err := database.InitDB(path)
	if err != nil {
		log.Fatalf("problem with init db: %v", err)
	}
	defer database.DB.Close()

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", taskByIdHandler)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// Обработчик для /tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.CreateTaskHandler(w, r)
	case http.MethodGet:
		handlers.GetTasksHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskByIdHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Printf("taskByIdHandler - Path: %s, Method: %s\n", path, r.Method)

	// Проверяем сначала на завершение задачи
	if strings.HasPrefix(path, "/tasks/complete/") {
		// запрос на завершение задачи
		if r.Method == http.MethodPost {
			fmt.Println("Routing to CompleteTaskHandler")
			handlers.CompleteTaskHandler(w, r)
		} else {
			fmt.Printf("Method not allowed for completion: %s\n", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasPrefix(path, "/tasks/") && len(path) > len("/tasks/") {
		// запрос на обновление задачи
		if r.Method == http.MethodPut {
			fmt.Println("Routing to UpdateTaskDateHandler")
			handlers.UpdateTaskDateHandler(w, r)
		} else {
			fmt.Printf("Method not allowed for update: %s\n", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		fmt.Println("Path not found")
		http.NotFound(w, r)
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
