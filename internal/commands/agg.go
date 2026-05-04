package commands

import (
	"context"
	"fmt"

	"github.com/davicbtoliveira/rss_aggregator/internal/rss"
)

func HandlerAgg(s *State, cmd Command) error {
	rssfeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v", rssfeed)
	return nil
}
