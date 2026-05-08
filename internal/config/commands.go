package config

import (

	"fmt"
	"os"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.CommandMap[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.CommandMap == nil {
		c.CommandMap = make(map[string]func(*State, Command) error)
	}
	c.CommandMap[name] = f
}






func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		fmt.Printf("usage: login <username>\n")
		os.Exit(1)
	}
	username := cmd.Args[0]
	
	if err := s.Cfg.SetUser(username, GetConfigFilePath()); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	return nil
}
