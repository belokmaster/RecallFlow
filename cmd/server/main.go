package main

import (
	"log"
	"net/http"
	"os"
	"reccal_flow/internal/database"
	"reccal_flow/internal/handlers"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Problem with inisialization BD: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /cards", enableCORS(handlers.CreateCardHandler(db)))
	mux.HandleFunc("GET /cards", enableCORS(handlers.GetCardsHandler(db)))
	mux.HandleFunc("PUT /cards/{id}", enableCORS(handlers.UpdateCardHandler(db)))
	mux.HandleFunc("DELETE /cards/{id}", enableCORS(handlers.DeleteCardHandler(db)))
	mux.HandleFunc("GET /repeat", enableCORS(handlers.GetCardsForReviewHandler(db)))
	mux.HandleFunc("POST /repeat/{id}", enableCORS(handlers.RepeatCardHandler(db)))

	// Serve frontend
	mux.HandleFunc("/", serveFrontend)

	log.Println("Server start at http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Problem with start server: %v", err)
	}
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	path := "web/frontend/dist" + r.URL.Path
	
	if _, err := os.Stat(path); err == nil {
		http.ServeFile(w, r, path)
		return
	}

	http.ServeFile(w, r, "web/frontend/dist/index.html")
}
