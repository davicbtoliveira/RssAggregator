package commands

import (
	"context"
	"fmt"
	"os"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not enough arguments")
	}

	if _, err := s.Db.GetUser(context.Background(), cmd.Args[0]); err != nil {
		fmt.Printf("user %v doesn't exists in the database\n", cmd.Args[0])
		os.Exit(1)
	}

	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("User %v has been set\n", cmd.Args[0])

	return nil
}
