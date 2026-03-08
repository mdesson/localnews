package source

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pemistahl/lingua-go"
)

// Language is exactly what it sounds like. However, because it's implemented as a bitmask, it supports multilingualism.
//
// For example, `LanguageEnglish|LanguageFrench` means something is bilingual.
type Language uint8

const (
	LanguageEnglish Language = 1 << iota
	LanguageFrench
)

// Source is a news source, including metadata for it, as well as its articles.
type Source struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	FeedURL  string    `json:"feed_url"`
	Filter   []string  `json:"filter,omitempty"`
	Language Language  `json:"language"`
	Articles []Article `json:"articles,omitempty"`
}

// FetchArticles refreshes the Articles on the Source, replacing them entirely.
func (s *Source) FetchArticles(detector lingua.LanguageDetector) error {
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

	s.Articles = make([]Article, len(feed.Items))
	for i, item := range feed.Items {
		language := LanguageFrench
		linguaLang, detected := detector.DetectLanguageOf(item.Title)
		if linguaLang == lingua.English {
			language = LanguageEnglish
		} else if !detected && s.Language != (LanguageFrench|LanguageEnglish) {
			language = s.Language
		}

		s.Articles[i] = Article{Item: *item, Language: language}
	}

	return nil
}

// Article contains every piece of information about the article, including the Language.
//
// The Language is detected intelligently when it is fetched. It's not a property of the rss feed item.
type Article struct {
	gofeed.Item
	Language Language
}

// UserLanguage extracts the user's language from an HTTP request.
//
// Note that bilingualism is a possible outcome of this operation.
func UserLanguage(r *http.Request) Language {
	inputLang := "fr"
	if c, err := r.Cookie("lang"); err == nil {
		if c.Value == "en" {
			inputLang = c.Value
		} else if c.Value == "bi" {
			inputLang = c.Value
		}
	} else {
		accept := r.Header.Get("Accept-Language")
		if strings.Contains(accept, "en") {
			inputLang = "en"
		}
	}

	if inputLang == "en" {
		return LanguageEnglish
	} else if inputLang == "bi" {
		return LanguageFrench & LanguageEnglish
	}

	return LanguageFrench
}
