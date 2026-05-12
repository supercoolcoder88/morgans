package articles

import (
	"database/sql"
	"fmt"
	"time"
)

type articlesStore struct {
	db *sql.DB
}

func NewArticlesStore(db *sql.DB) *articlesStore {
	return &articlesStore{
		db: db,
	}
}

type Article struct {
	ID          string `json:"id"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	AddedDate   string `json:"addedDate"`
	IsRead      bool   `json:"isRead"`
}

func (s *articlesStore) GetAllArticlesToday() ([]Article, error) {
	today := time.Now().Format("2006-01-02")

	query := `SELECT id, source, title, link, description, addedDate, isRead
		FROM Article
		WHERE DATE(addedDate) = ?`

	rows, err := s.db.Query(query, today)
	if err != nil {
		return nil, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article

		err := rows.Scan(
			&article.ID,
			&article.Source,
			&article.Title,
			&article.Link,
			&article.Description,
			&article.AddedDate,
			&article.IsRead,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning article: %w", err)
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, nil
}

func (s *articlesStore) GetArticle(id string) (*Article, error) {
	query := `SELECT id, source, title, link, description, addedDate, isRead
		FROM Article
		WHERE id = ?`

	var article Article

	err := s.db.QueryRow(query, id).Scan(
		&article.ID,
		&article.Source,
		&article.Title,
		&article.Link,
		&article.Description,
		&article.AddedDate,
		&article.IsRead,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error querying article: %w", err)
	}

	return &article, nil
}

func (s *articlesStore) UpdateIsRead(id string, isRead bool) error {
	query := `UPDATE Article SET isRead = ? WHERE id = ?`

	result, err := s.db.Exec(query, isRead, id)
	if err != nil {
		return fmt.Errorf("error updating article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (s *articlesStore) AddArticle(id, source, title, link, description string) error {
	query := `INSERT OR IGNORE INTO Article (id, source, title, link, description)
		VALUES (?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, id, source, title, link, description)
	if err != nil {
		return fmt.Errorf("error inserting article: %w", err)
	}

	return nil
}
