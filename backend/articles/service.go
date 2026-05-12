package articles

import (
	"database/sql"
	"encoding/json"
	"log"

	llama "github.com/supercoolcoder88/llamacpp-go"
)

type service struct {
	db *sql.DB
}

type filterResponse struct {
	ids []string
}

var sources = map[string]string{
	"ycombinator":  "https://news.ycombinator.com/rss",
	"abc_vic":      "https://www.abc.net.au/news/feed/5470430/rss.xml",
	"abc_world":    "https://www.abc.net.au/news/feed/104217382/rss.xml",
	"abc_business": "https://www.abc.net.au/news/feed/104217374/rss.xml",
	"abc_top":      "https://www.abc.net.au/news/feed/10719986/rss.xml",
}

func NewService(db *sql.DB) *service {
	return &service{
		db: db,
	}
}

func (s *service) FetchArticles() error {
	log.Print("Fetching articles")
	articles := make(map[string]feedItem)

	// Fetch all the articles for each source
	for source, url := range sources {
		items := readRSSFeed(source, url)

		for _, item := range items {
			articles[item.GUID] = item
		}
	}

	// Pass the articles to AI for filtering
	ids, err := filterArticles(articles)
	if err != nil {
		return err
	}

	// Filter out the valid GUIDs
	var filteredItems []feedItem
	for _, id := range ids {
		_, exists := articles[id]

		if exists {
			filteredItems = append(filteredItems, articles[id])
		}
	}

	// Pass results from AI to database
	store := NewArticlesStore(s.db)

	for _, article := range articles {
		store.AddArticle(article.GUID, article.Source, article.Title, article.Link, article.Description)
	}

	return nil
}

func filterArticles(items map[string]feedItem) ([]string, error) {
	client := llama.New("http://localhost:9090") // Put this in env variable
	log.Print("Querying llama.cpp")
	// Build messages
	system := llama.Message{
		Role: "system",
		Content: `You are a JSON API,
			Look through the provided items and select only items related to:
			- Software technology
			- Programming
			- Finance
			- AI
			- Important political news
			Return ONLY valid JSON.
			Schema:
			{
			"ids": ["guid1", "guid2", "guid3"]
			}
			Do not include explanations, markdown, or extra text.`,
	}

	j, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	user := llama.Message{
		Role:    "user",
		Content: string(j),
	}

	res, err := client.ChatJSON("llama.cpp", []llama.Message{system, user})

	if err != nil {
		return nil, err
	}

	var filteredIds filterResponse

	if err := json.Unmarshal([]byte(res), &filteredIds); err != nil {
		return nil, err
	}

	log.Printf("llm returned %v IDs", len(filteredIds.ids))
	return filteredIds.ids, nil
}
