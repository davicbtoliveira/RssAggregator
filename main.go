package main

import (
	"github.com/davicbtoliveira/rss_aggregator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser()

	config.ReadGatorconfigContent()
}
