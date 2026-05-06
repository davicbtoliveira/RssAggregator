package commands

import (
	"context"
	"fmt"

	"github.com/davicbtoliveira/rss_aggregator/internal/database"
)

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Not enough arguments for unfollow command")
	}

	feed, err := s.Db.GetFeedFromUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	if err = s.Db.DeleteFollow(context.Background(), database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return err
	}

	return nil
}
