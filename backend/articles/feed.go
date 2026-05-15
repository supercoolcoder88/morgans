package articles

import (
	"log"

	"github.com/mmcdole/gofeed"
)

type feedItem struct {
	Title       string
	Link        string
	Description string
	Source      string
}

func readRSSFeed(source string, url string) []feedItem {
	log.Printf("Reading feed from %s", source)
	fp := gofeed.NewParser()

	feed, _ := fp.ParseURL(url)

	log.Printf("number of articles read: %v", feed.Len())

	feedItems := []feedItem{}

	for _, i := range feed.Items {
		feedItems = append(feedItems, feedItem{
			Title:       i.Title,
			Link:        i.Link,
			Description: i.Description,
			Source:      source,
		})
	}

	return feedItems
}
