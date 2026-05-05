package commands

import (
	"context"
	"fmt"
)

func HandlerFeeds(s *State, cmd Command) error {
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, v := range feeds {
		fmt.Println(v.Name)
		fmt.Println(v.Url)
		user, err := s.Db.GetUserById(context.Background(), v.UserID)
		if err != nil {
			return err
		}
		fmt.Println(user.Name)
	}

	return nil
}
