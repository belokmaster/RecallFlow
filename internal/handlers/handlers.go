package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"reccal_flow/internal/database"
	"strconv"
)

// Card Handlers
func CreateCardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateCardHandler: Начало обработки")

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		var req CreateCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		if req.Word == "" || len(req.Word) > 250 {
			respondWithError(w, http.StatusBadRequest, "Слово обязательно и не должно превышать 250 символов")
			return
		}

		if req.Translation == "" || len(req.Translation) > 250 {
			respondWithError(w, http.StatusBadRequest, "Перевод обязателен и не должен превышать 250 символов")
			return
		}

		card, err := database.CreateCard(db, req.Word, req.Translation, req.Example)
		if err != nil {
			log.Printf("CreateCardHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при создании карточки")
			return
		}

		respondWithJSON(w, http.StatusCreated, card)
	}
}

func GetCardsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetCardsHandler: Начало обработки")

		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		cards, err := database.GetAllCards(db)
		if err != nil {
			log.Printf("GetCardsHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении карточек")
			return
		}

		respondWithJSON(w, http.StatusOK, CardsResponse{Cards: cards})
	}
}

func DeleteCardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DeleteCardHandler: Начало обработки")

		if r.Method != http.MethodDelete {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		cardID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID карточки")
			return
		}

		if err := database.DeleteCard(db, cardID); err != nil {
			log.Printf("DeleteCardHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при удалении карточки")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func UpdateCardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateCardHandler: Начало обработки")

		if r.Method != http.MethodPut {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		cardID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID карточки")
			return
		}

		var req UpdateCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		if req.Word == "" || len(req.Word) > 250 {
			respondWithError(w, http.StatusBadRequest, "Слово обязательно и не должно превышать 250 символов")
			return
		}

		if req.Translation == "" || len(req.Translation) > 250 {
			respondWithError(w, http.StatusBadRequest, "Перевод обязателен и не должен превышать 250 символов")
			return
		}

		if err := database.UpdateCard(db, cardID, req.Word, req.Translation, req.Example); err != nil {
			log.Printf("UpdateCardHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении карточки")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}

func GetCardsForReviewHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetCardsForReviewHandler: Начало обработки")

		if r.Method != http.MethodGet {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		cards, err := database.GetCardsForReview(db)
		if err != nil {
			log.Printf("GetCardsForReviewHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении карточек для повторения")
			return
		}

		respondWithJSON(w, http.StatusOK, CardsResponse{Cards: cards})
	}
}

func RepeatCardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("RepeatCardHandler: Начало обработки")

		if r.Method != http.MethodPost {
			respondWithError(w, http.StatusMethodNotAllowed, "Метод не разрешен")
			return
		}

		idStr := r.PathValue("id")
		cardID, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный ID карточки")
			return
		}

		var req RepeatCardRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Неверный формат JSON")
			return
		}

		if err := database.RepeatCard(db, cardID, req.Success); err != nil {
			log.Printf("RepeatCardHandler: Ошибка БД: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении карточки")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
	}
}
