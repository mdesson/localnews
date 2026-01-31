package source

import (
	"context"
	"errors"
	"time"

	"github.com/mmcdole/gofeed"
)

type Language uint8

const (
	LanguageEnglish Language = 1 << iota
	LanguageFrench
)

type Source struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	FeedURL  string        `json:"feed_url"`
	Filter   []string      `json:"filter,omitempty"`
	Language Language      `json:"language"`
	Articles []gofeed.Item `json:"articles,omitempty"`
}

func (s *Source) FetchArticles() error {
	fp := gofeed.NewParser()
	if s.FeedURL == "" {
		return errors.New("no feed url")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	feed, err := fp.ParseURLWithContext(s.FeedURL, ctx)
	if err != nil {
		return err
	}

	s.Articles = make([]gofeed.Item, len(feed.Items))
	for i, item := range feed.Items {
		s.Articles[i] = *item
	}

	return nil
}
