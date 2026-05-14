package articles

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	log.Print("Getting articles for today")
	service := NewService(h.db)

	articles, err := service.FetchArticlesToday()

	if err != nil {
		http.Error(w, "Internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(articles)
}
