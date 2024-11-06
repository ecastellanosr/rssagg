package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ecastellanosr/rssagg/internal/database"
	"github.com/google/uuid"
)

func scrapefeeds(s *state, limit int32, TimeBetweenRequests time.Duration) error {
	start := time.Now()
	ticker := time.NewTicker(TimeBetweenRequests)
	nfeeds, err := s.db.GetNumberOfFeeds(context.Background())
	if limit > int32(nfeeds) {
		limit = int32(nfeeds)
	}
	if err != nil {
		return fmt.Errorf("error while getting the number of feeds, %w", err)
	}
	for i := 0; ; <-ticker.C {
		ch := make(chan error)
		feeds, err := s.db.GetNextFeedToFetch(context.Background(), limit)
		if err != nil {
			return fmt.Errorf("error while fetching the next feed, %w", err)
		}
		fmt.Println("----------------------------------------------------------")
		for _, feed := range feeds {
			go scrapefeed(feed, s, ch)
			err = <-ch
			if err != nil {
				return err
			}
		}

		fmt.Println("----------------------------------------------------------")
		i++
		if i*int(limit) >= int(nfeeds) {
			elapsed := time.Since(start)
			fmt.Printf("Binomial took %s\n", elapsed)
			break
		}
	}
	return nil
}
func scrapefeed(feed database.GetNextFeedToFetchRow, s *state, ch chan error) {

	err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		ch <- fmt.Errorf("error while marking the feed as fetched, %w", err)
		return
	}
	rssfeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		ch <- fmt.Errorf("error while fetching the feeds items, %w", err)
		return
	}
	fmt.Println(rssfeed.Channel.Title)

	for _, item := range rssfeed.Channel.Item {
		valid_des := true
		if item.Description == "" {
			valid_des = false
		}

		valid_pubdate := true
		pubdate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			valid_pubdate = false
		}

		postparams := database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  valid_des,
			},
			Url: feed.Url,
			PublishedAt: sql.NullTime{
				Time:  pubdate,
				Valid: valid_pubdate,
			},
			FeedID: feed.ID,
		}
		post, err := s.db.CreatePost(context.Background(), postparams)
		if err != nil {
			ch <- fmt.Errorf("error while creating the post %w", err)
			return
		}
		fmt.Printf("Post Title: %v\n Post ID: %v\n Description: %v\n Published in: %v\n", post.Title, post.ID, post.Description, post.PublishedAt)
	}
	ch <- nil
}
