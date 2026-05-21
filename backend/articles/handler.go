package articles

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
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
	start := time.Now()
	timeFilter := r.URL.Query().Get("time")
	todayOnly := timeFilter == "today"

	log.Printf("[http] GET /articles?time=%s from %s", timeFilter, r.RemoteAddr)

	service := NewService(h.db)
	articles, err := service.FetchArticles(todayOnly)
	if err != nil {
		log.Printf("[http] GET /articles ERROR: %v (%v)", err, time.Since(start))
		http.Error(w, "Internal server error occurred", http.StatusInternalServerError)
		return
	}

	log.Printf("[http] GET /articles returning %d articles (%v)", len(articles), time.Since(start))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(articles)
}

type updateIsReadRequest struct {
	Link   string `json:"link"`
	IsRead bool   `json:"isRead"`
}

func (h *Handler) UpdateIsRead(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("[http] PATCH /articles from %s", r.RemoteAddr)

	service := NewService(h.db)

	var req updateIsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[http] PATCH /articles bad request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Link == "" {
		log.Printf("[http] PATCH /articles missing link field")
		http.Error(w, "Link is required", http.StatusBadRequest)
		return
	}

	if err := service.UpdateArticleIsRead(req.Link, req.IsRead); err != nil {
		log.Printf("[http] PATCH /articles ERROR updating %q: %v", req.Link, err)
		http.Error(w, "Internal server error occurred", http.StatusInternalServerError)
		return
	}

	log.Printf("[http] PATCH /articles updated isRead=%v for %q (%v)", req.IsRead, req.Link, time.Since(start))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
