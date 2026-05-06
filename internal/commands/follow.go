package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/davicbtoliveira/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("not enough arguments for follow command")
	}

	feed, err := s.Db.GetFeedFromUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	follow_feed, err := s.Db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Println(follow_feed.FeedName)
	fmt.Println(follow_feed.UserName)

	return nil
}
