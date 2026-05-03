package commands

import (
	"github.com/davicbtoliveira/rss_aggregator/internal/config"
	"github.com/davicbtoliveira/rss_aggregator/internal/database"
)

type State struct {
	Db  *database.Queries
	Cfg *config.Config
}
