package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/davicbtoliveira/rss_aggregator/internal/database"
)

func HandlerBrowse(s *State, cmd Command) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("too many arguments for browse command")
	}

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}

	limit := 2
	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.Db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\nURL: %s\n\n", post.Title, post.Url)
	}

	return nil
}
