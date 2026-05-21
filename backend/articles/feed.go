package articles

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

type feedItem struct {
	Title       string
	Link        string
	Description string
	Source      string
}

func readRSSFeed(source string, url string) ([]feedItem, error) {
	start := time.Now()
	log.Printf("[feed] Fetching %s ...", source)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %w", source, err)
	}

	feedItems := make([]feedItem, 0, len(feed.Items))
	for _, i := range feed.Items {
		feedItems = append(feedItems, feedItem{
			Title:       i.Title,
			Link:        i.Link,
			Description: i.Description,
			Source:      source,
		})
	}

	log.Printf("[feed] %s: %d articles in %v", source, len(feedItems), time.Since(start))
	return feedItems, nil
}
