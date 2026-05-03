package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davicbtoliveira/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not enoght arguments")
	}

	if _, err := s.Db.GetUser(context.Background(), cmd.Args[0]); err == nil {
		fmt.Printf("user %v already exists in the database\n", cmd.Args[0])
		os.Exit(1)
	}

	s.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	})

	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("user %v registered\n", cmd.Args[0])

	return nil
}
