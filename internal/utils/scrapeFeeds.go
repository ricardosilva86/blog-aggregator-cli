package utils

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ricardosilva86/blogaggregator/internal/config"
	"github.com/ricardosilva86/blogaggregator/internal/database"
	"time"
)

// ScrapeFeeds will scrape all feeds in the database
func ScrapeFeeds(db *database.Queries) error {
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	newFeed, err := config.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error scraping feed: %w", err)
	}

	if _, err := db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("error marking feed fetched: %w", err)
	}

	for _, item := range newFeed.Channel.Item {
		longLayout := "Mon, 02 Jan 2006 15:04:05 -0700"
		shortLayout := "2006-Jan-02"
		pubDate := time.Now()
		if pubDate, err = time.Parse(longLayout, item.PubDate); err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}
		if pubDate, err = time.Parse(shortLayout, pubDate.Format(shortLayout)); err != nil {
			return fmt.Errorf("error parsing date: %w", err)
		}

		createPostParams := database.CreatePostParams{
			ID:          uuid.New(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pubDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		_, err = db.CreatePost(context.Background(), createPostParams)
		if err != nil {
			fmt.Printf("Error type: %T\nError message: %v\n", err, err)
		}
		fmt.Println(item.Title)
	}

	return nil
}
