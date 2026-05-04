package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/davicbtoliveira/rss_aggregator/internal/commands"
	"github.com/davicbtoliveira/rss_aggregator/internal/config"
	"github.com/davicbtoliveira/rss_aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Read()

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	state := commands.State{
		Db:  dbQueries,
		Cfg: &cfg,
	}

	cmds := commands.Commands{
		CommandHandler: make(map[string]func(*commands.State, commands.Command) error),
	}
	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerListUsers)
	cmds.Register("agg", commands.HandlerAgg)

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
