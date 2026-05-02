package commands

import "fmt"

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("not enought arguments")
	}

	if err := s.Cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}

	fmt.Printf("User %v has been set\n", cmd.Args[0])

	return nil
}
