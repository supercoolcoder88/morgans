package main

import (
	"database/sql"
	"log"
	"morgans/articles"
	"morgans/repositories/sqlite"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	mux := http.NewServeMux()

	// Initialise the database
	db, err := sql.Open("sqlite3", "/data/morgans.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := sqlite.Init(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create the articles handler and register routes
	handler := articles.NewHandler(db)
	mux.HandleFunc("GET /articles", handler.GetArticles)

	// Create the articles service
	service := articles.NewService(db)

	// Run FetchArticles immediately
	if err := service.ReadArticlesFromRSSFeeds(); err != nil {
		log.Printf("Error fetching articles: %v", err)
	} else {
		log.Println("Successfully fetched articles")
	}

	// Start HTTP server in a goroutine
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	go func() {
		ticker := time.NewTicker(4 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			if err := service.ReadArticlesFromRSSFeeds(); err != nil {
				log.Printf("Error fetching articles: %v", err)
			} else {
				log.Println("Successfully fetched articles")
			}
		}
	}()

	log.Printf("Starting HTTP server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
