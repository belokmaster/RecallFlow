package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"reccal_flow/internal/database"
)

// Card request types
type CreateCardRequest struct {
	Word        string  `json:"word"`
	Translation string  `json:"translation"`
	Example     *string `json:"example"`
}

type CardsResponse struct {
	Cards []database.Card `json:"cards"`
	Error string          `json:"error,omitempty"`
}

type UpdateCardRequest struct {
	Word        string  `json:"word"`
	Translation string  `json:"translation"`
	Example     *string `json:"example"`
}

type RepeatCardRequest struct {
	Success bool `json:"success"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Printf("Error: code %d, message: %s", code, message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Critical error: failed to marshal JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
