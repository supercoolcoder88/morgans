package articles

import (
	"log"

	"github.com/mmcdole/gofeed"
)

type feedItem struct {
	GUID        string
	Title       string
	Link        string
	Description string
	Source      string
}

func readRSSFeed(source string, url string) []feedItem {
	fp := gofeed.NewParser()

	feed, _ := fp.ParseURL(url)

	log.Printf("number of links: %v", feed.Len())

	feedItems := []feedItem{}

	for _, i := range feed.Items {
		feedItems = append(feedItems, feedItem{
			GUID:        i.GUID,
			Title:       i.Title,
			Link:        i.Link,
			Description: i.Description,
			Source:      source,
		})
	}

	return feedItems
}
