package commands

import (
	"context"
	"fmt"
)

func HandlerFollowing(s *State, cmd Command) error {
	following, err := s.Db.GetFeedFollowsForUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}

	for _, v := range following {
		fmt.Println(v.FeedName)
	}

	return nil
}
