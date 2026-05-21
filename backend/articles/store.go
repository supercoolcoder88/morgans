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
	Link        string `json:"link"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AddedDate   string `json:"addedDate"`
	IsRead      bool   `json:"isRead"`
}

func (s *articlesStore) GetArticles(todayOnly bool) ([]Article, error) {
	var query string
	var rows *sql.Rows
	var err error

	if todayOnly {
		today := time.Now().Format("2006-01-02")
		query = `SELECT link, source, title, description, addedDate, isRead
			FROM Article
			WHERE DATE(addedDate) = ?
			ORDER BY addedDate DESC`
		rows, err = s.db.Query(query, today)
	} else {
		query = `SELECT link, source, title, description, addedDate, isRead
			FROM Article
			ORDER BY addedDate DESC`
		rows, err = s.db.Query(query)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying articles: %w", err)
	}
	defer rows.Close()

	var articles []Article

	for rows.Next() {
		var article Article

		err := rows.Scan(
			&article.Link,
			&article.Source,
			&article.Title,
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

func (s *articlesStore) GetArticle(link string) (*Article, error) {
	query := `SELECT link, source, title, description, addedDate, isRead
		FROM Article
		WHERE link = ?`

	var article Article

	err := s.db.QueryRow(query, link).Scan(
		&article.Link,
		&article.Source,
		&article.Title,
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

func (s *articlesStore) UpdateIsRead(link string, isRead bool) error {
	query := `UPDATE Article SET isRead = ? WHERE link = ?`

	result, err := s.db.Exec(query, isRead, link)
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

func (s *articlesStore) AddArticle(link, source, title, description string) error {
	query := `INSERT OR IGNORE INTO Article (link, source, title, description)
		VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(query, link, source, title, description)
	if err != nil {
		return fmt.Errorf("error inserting article: %w", err)
	}

	return nil
}
