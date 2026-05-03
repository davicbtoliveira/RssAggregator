package commands

import (
	"context"
	"fmt"
)

func HandlerReset(s *State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("too many args for reset command")
	}

	if err := s.Db.ResetUsers(context.Background()); err != nil {
		fmt.Printf("error when reseting users database: %v\n", err)
	}

	return nil
}
