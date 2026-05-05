package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/davicbtoliveira/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Not enought arguments for addfeed")
	}

	user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
	if err != nil {
		return err
	}

	_, err = s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	followCmd := Command{
		Name: "follow",
		Args: cmd.Args[1:],
	}
	if err = HandlerFollow(s, followCmd); err != nil {
		return err
	}

	return nil
}
