package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/davicbtoliveira/rss_aggregator/internal/rss"
)

func scrapeFeeds(s *State) error {
	nextFetch, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	if err := s.Db.MarkFeedFetched(context.Background(), nextFetch.ID); err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), nextFetch.Url)
	if err != nil {
		return err
	}

	for _, v := range rssFeed.Channel.Item {
		fmt.Println(v.Title)
	}

	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough arguments to fetch feed")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("Collecting feeds every ", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}
