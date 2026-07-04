package source

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pemistahl/lingua-go"
)

var (
	ErrorNoFeedURL = errors.New("no feed url")
)

// Language is exactly what it sounds like. However, because it's implemented as a bitmask, it supports multilingualism.
//
// For example, `LanguageEnglish|LanguageFrench` means something is bilingual.
type Language uint8

const (
	LanguageEnglish Language = 1 << iota
	LanguageFrench
)

func (l Language) String() string {
	if l == LanguageEnglish {
		return "en"
	}
	return "fr"
}

// Source is a news source, including metadata for it, as well as its articles.
type Source struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	FeedURL       string    `json:"feed_url"`
	Filter        []string  `json:"filter,omitempty"`
	Language      Language  `json:"language"`
	Articles      []Article `json:"articles,omitempty"`
	ErrorOnUpdate bool
}

// FetchArticles refreshes the Articles on the Source, replacing them entirely.
func (s *Source) FetchArticles(detector lingua.LanguageDetector, userAgent string) error {
	fp := gofeed.NewParser()
	if userAgent != "" {
		fp.UserAgent = userAgent
	}

	if s.FeedURL == "" {
		return ErrorNoFeedURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	feed, err := fp.ParseURLWithContext(s.FeedURL, ctx)
	if err != nil {
		return err
	}

	s.Articles = make([]Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		// enforce source's filter
		if s.Filter != nil {
			skip := true
			titleLower := strings.ToLower(item.Title)
			descLower := strings.ToLower(item.Description)
			contentLower := strings.ToLower(item.Content)
			for _, keyword := range s.Filter {
				keywordLower := strings.ToLower(keyword)
				if strings.Contains(titleLower, keywordLower) || strings.Contains(descLower, keywordLower) || strings.Contains(contentLower, keywordLower) {
					skip = false
					break
				}
			}
			if skip {
				continue
			}
		}

		// get article language
		language := LanguageFrench
		linguaLang, detected := detector.DetectLanguageOf(item.Title)
		if linguaLang == lingua.English {
			language = LanguageEnglish
		} else if !detected && s.Language != (LanguageFrench|LanguageEnglish) {
			language = s.Language
		}

		// Some items might have pictures in them
		item.Description = stripImages(item.Description)

		// item.Published isn't always filled in, attempt to get it from elsewhere
		// On Neomedia this is where it is hidden
		if item.Published == "" && item.PublishedParsed == nil && item.UpdatedParsed == nil {
			if a10, ok := item.Extensions["a10"]; ok {
				if updated, ok := a10["updated"]; ok && len(updated) > 0 {
					t, err := time.Parse(time.RFC3339, updated[0].Value)
					if err == nil {
						item.PublishedParsed = &t
						item.Published = t.Format(time.RFC1123)
					}
				}
			}
		}

		s.Articles = append(s.Articles, Article{
			Item:             *item,
			Language:         language,
			SourceName:       s.Name,
			SourceURL:        s.URL,
			SelectedSourceID: fmt.Sprintf("%s:%s", s.ID, language.String()),
		})
	}

	return nil
}

func (s *Source) CheckboxIDs(language Language) []string {
	var ids []string
	if (s.Language&LanguageFrench != 0) && (language&LanguageFrench != 0) {
		ids = append(ids, s.ID+":fr")
	}
	if (s.Language&LanguageEnglish != 0) && (language&LanguageEnglish != 0) {
		ids = append(ids, s.ID+":en")
	}
	return ids
}

func stripImages(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}
	doc.Find("img").Remove()
	result, _ := doc.Html()
	return result
}

// Article contains every piece of information about the article, including the Language.
//
// The Language is detected intelligently when it is fetched. It's not a property of the rss feed item.
type Article struct {
	gofeed.Item
	Language         Language
	SourceName       string
	SourceURL        string
	SelectedSourceID string
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
		return LanguageFrench | LanguageEnglish
	}

	return LanguageFrench
}
