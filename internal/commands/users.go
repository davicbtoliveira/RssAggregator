package commands

import (
	"context"
	"fmt"
)

func HandlerListUsers(s *State, cmd Command) error {
	res, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range res {
		if u.Name == s.Cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", u.Name)
		} else {
			fmt.Printf("* %v\n", u.Name)
		}
	}
	return nil
}
