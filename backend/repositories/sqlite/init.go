package sqlite

import (
	"database/sql"
	"log"
)

// Create the SQL tables
func Init(db *sql.DB) error {
	createArticlesTable := `
	CREATE TABLE IF NOT EXISTS Article (
		id TEXT,
		source TEXT,
		title TEXT,
		link TEXT,
		description TEXT,
		addedDate DATETIME DEFAULT CURRENT_TIMESTAMP,
		isRead BOOLEAN DEFAULT 0,
		PRIMARY KEY (id, source)
	);`

	_, err := db.Exec(createArticlesTable)
	if err != nil {
		log.Printf("Error creating Article table: %v", err)
		return err
	}

	log.Println("Article table created successfully")
	return nil
}
