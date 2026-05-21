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

type updateIsReadRequest struct {
	Link   string `json:"link"`
	IsRead bool   `json:"isRead"`
}

func (h *Handler) UpdateIsRead(w http.ResponseWriter, r *http.Request) {
	log.Print("Updating article isRead")
	service := NewService(h.db)

	var req updateIsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Link == "" {
		http.Error(w, "Link is required", http.StatusBadRequest)
		return
	}

	if err := service.UpdateArticleIsRead(req.Link, req.IsRead); err != nil {
		http.Error(w, "Internal server error occurred", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
