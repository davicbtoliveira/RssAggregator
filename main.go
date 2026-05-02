package main

import (
	"fmt"
	"os"

	"github.com/davicbtoliveira/rss_aggregator/internal/commands"
	"github.com/davicbtoliveira/rss_aggregator/internal/config"
)

func main() {
	cfg := config.Read()

	state := commands.State{
		Cfg: &cfg,
	}

	cmds := commands.Commands{
		CommandHandler: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)

	arguments := os.Args
	if len(arguments) < 2 {
		fmt.Println("not enought arguments")
		os.Exit(1)
	}
	cmd := commands.Command{
		Name: arguments[1],
		Args: arguments[2:],
	}

	if err := cmds.Run(&state, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
