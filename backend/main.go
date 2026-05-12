package main

import (
	"database/sql"
	"log"
	"morgans/articles"
	"morgans/repositories/sqlite"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//mux := http.NewServeMux()

	// Initialise the database
	db, err := sql.Open("sqlite3", "./morgans.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := sqlite.Init(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create the articles service
	service := articles.NewService(db)

	// Run FetchArticles immediately
	if err := service.FetchArticles(); err != nil {
		log.Printf("Error fetching articles: %v", err)
	} else {
		log.Println("Successfully fetched articles")
	}

	// Set up ticker to run FetchArticles every 4 hours
	ticker := time.NewTicker(4 * time.Hour)
	defer ticker.Stop()

	log.Println("Starting article fetcher (runs every 4 hours)")

	for range ticker.C {
		if err := service.FetchArticles(); err != nil {
			log.Printf("Error fetching articles: %v", err)
		} else {
			log.Println("Successfully fetched articles")
		}
	}
}
