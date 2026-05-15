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
	IDs []int `json:"ids"`
}

type feedItemForLLM struct {
	Title       string `json:"t"`
	Description string `json:"d"`
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

func (s *service) ReadArticlesFromRSSFeeds() error {
	log.Print("Fetching articles")
	articles := make(map[int]feedItem)
	counter := 0 // using counter so that llm can link articles

	// Fetch all the articles for each source
	for source, url := range sources {
		items := readRSSFeed(source, url)

		for _, item := range items {
			articles[counter] = item
			counter++
		}
	}
	articleToFilter := formatFeedItemsForLLM(articles)
	// Pass the articles to AI for filtering
	ids, err := filterArticles(articleToFilter)

	if err != nil {
		return err
	}

	// Filter out the valid IDs
	var filteredItems []feedItem
	for _, id := range ids {
		item, exists := articles[id]

		if exists {
			filteredItems = append(filteredItems, item)
		}
	}

	// Pass results from AI to database
	store := NewArticlesStore(s.db)

	log.Printf("Adding %v items", len(filteredItems))
	for _, article := range filteredItems {
		store.AddArticle(article.Link, article.Source, article.Title, article.Description)
	}

	return nil
}

func (s *service) FetchArticlesToday() ([]Article, error) {
	store := NewArticlesStore(s.db)
	articles, err := store.GetAllArticlesToday()

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func formatFeedItemsForLLM(items map[int]feedItem) map[int]feedItemForLLM {
	formatted := make(map[int]feedItemForLLM)
	for id, item := range items {
		formatted[id] = feedItemForLLM{
			Title:       item.Title,
			Description: item.Description,
		}
	}
	return formatted
}

func filterArticles(items map[int]feedItemForLLM) ([]int, error) {
	client := llama.New("http://llama-cpp-server:9090")
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
			"ids": [0, 1, 2, 3]
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

	log.Printf("llm returned %v ids", len(filteredIds.IDs))
	return filteredIds.IDs, nil
}
