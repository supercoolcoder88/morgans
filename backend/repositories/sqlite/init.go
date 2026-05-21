package sqlite

import (
	"database/sql"
	"log"
)

// Create the SQL tables
func Init(db *sql.DB) error {
	createArticlesTable := `
	CREATE TABLE IF NOT EXISTS Article (
		link TEXT PRIMARY KEY,
		source TEXT,
		title TEXT,
		description TEXT,
		addedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
		isRead BOOLEAN DEFAULT 0
	);`

	_, err := db.Exec(createArticlesTable)
	if err != nil {
		log.Printf("Error creating Article table: %v", err)
		return err
	}

	log.Println("Article table created successfully")
	return nil
}
