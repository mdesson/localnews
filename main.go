package main

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type Language uint8

const (
	English Language = 1 << iota
	French
)

type Source struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	FeedURL  string        `json:"feed_url"`
	Filter   []string      `json:"filter,omitempty"`
	Language Language      `json:"language"`
	Articles []gofeed.Item `json:"articles,omitempty"`
}

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.neomedia.com/vaudreuil-soulanges/Rss/RssFeed")
	if err != nil {
		panic(err)
	}
	fmt.Println(feed)
}
