// commands.go
package main

import (

	"fmt"
	"context"
	"time"
	"gator/internal/config"
	"gator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) Run(s *state, cmd command) error {
	if handler, ok := c.commandMap[cmd.name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) Register(name string, f func(*state, command) error) {
	if c.commandMap == nil {
		c.commandMap = make(map[string]func(*state, command) error)
	}
	c.commandMap[name] = f
}






func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: login <username>\n")

	}
	username := cmd.args[0]
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, username)
	if err != nil {
		return fmt.Errorf("error: %v\n", err)
	}
	
	if err := s.cfg.SetUser(user.Name, config.GetConfigFilePath()); err != nil {
		return fmt.Errorf("error: %v\n", err)
	}
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: register <username>\n")
	}
	username := cmd.args[0]
	ctx := context.Background()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("[CreateUser]error: %v\n", err)
	}

	if err := s.cfg.SetUser(user.Name, config.GetConfigFilePath()); err != nil {
		return fmt.Errorf("[SetUser]error: %v\n", err)
	}
	fmt.Printf("User '%s' registered successfully with ID: %s\n", user.Name, user.ID)


	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.Reset(ctx); err != nil {
		return fmt.Errorf("[Reset]error: %v\n", err)
	}
	fmt.Println("All users have been deleted.")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("[GetAllUsers]error: %v\n", err)
	}
	currentUser := s.cfg.CurrentUserName
	// fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}

	}
	return nil
}
