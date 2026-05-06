package middleware

import (
	"context"

	"github.com/davicbtoliveira/rss_aggregator/internal/commands"
	"github.com/davicbtoliveira/rss_aggregator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *commands.State, cmd commands.Command, user database.User) error) func(*commands.State, commands.Command) error {
	return func(s *commands.State, cmd commands.Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
