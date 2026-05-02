package commands

import (
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandHandler map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.CommandHandler[cmd.Name]
	if !ok {
		return fmt.Errorf("command %v doesn't exists", cmd.Name)
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandHandler[name] = f
}
