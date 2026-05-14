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


func handlerAgg(s *state, cmd command) error {
	fmt.Printf("Running aggregate command: %s with args: %v\n", cmd.name, cmd.args)
	// rssFeed, err := fetchFeed(context.Background(), cmd.args[0])
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching feed: %v\n", err)
	}
	fmt.Printf("Fetched feed: %+v\n", rssFeed)

	return nil
}	

func handlerAddFeed(s *state, cmd command, user database.User) error {
	fmt.Printf("Running addfeed command: %s", cmd.name)
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <feed_name> <feed_url>\n")
	}

	var feedName string = cmd.args[0]
	var feedURL string = cmd.args[1]

	ctx := context.Background()

	params := database.CreateFeedParams{
		Name: feedName,
		Url: feedURL,
		UserID: user.ID,
	}

	s.db.CreateFeed(ctx, params)


	feed, err := s.db.GetFeedByURL(ctx, feedURL)
	if err != nil {
		fmt.Errorf("Error accessing feed.")
	}

	params_ff := database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	s.db.CreateFeedFollow(ctx, params_ff)


	fmt.Printf("Added feed '%s' with URL '%s' for user '%s'\n", feedName, feedURL, user.Name)

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return fmt.Errorf("[GetAllFeeds]error: %v\n", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("Registered Feeds:")
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("error fetching user for feed: %v\n", err)
		}
		fmt.Printf("* %s (%s), User: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	feedURL := cmd.args[0]

	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: follow <feed_id>\n")
	}


	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err !=  nil {
		return fmt.Errorf("error following feed: %w\n", err)
	}

	fmt.Printf("User '%s' is now following the '%s' feed\n", user.Name, feed.Name)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
			fmt.Errorf("LOGGING: %v\n", err)
		}
	fmt.Printf("%s follows these blogs:\n", user.Name)
	for _, item := range follows {
		feed, err := s.db.GetFeedByID(context.Background(), item.FeedID)
		if err != nil {
			fmt.Errorf("LOGGING: %v\n", err)
		}
		fmt.Printf("* %s\n", feed.Name)
	}
	return nil
	
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: follow <feed_id>\n")
	}
	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	params := database.UnfollowParams{
		UserID : user.ID,
		FeedID : feed.ID,
	}
	err = s.db.Unfollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}
	return nil
}