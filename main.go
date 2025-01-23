package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ricardosilva86/blogaggregator/internal/config"
	"github.com/ricardosilva86/blogaggregator/internal/database"
	"github.com/ricardosilva86/blogaggregator/internal/utils"
	"os"
	"time"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	command map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.command[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if f, ok := c.command[cmd.name]; ok {
		return f(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> <args> (optional)")
	}

	c, err := config.Read()
	if err != nil {
		fmt.Println(fmt.Errorf("error reading config file: %w", err))
		os.Exit(1)
	}

	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		fmt.Println(fmt.Errorf("error opening database: %w", err))
		os.Exit(1)
	}

	dbQueries := database.New(db)
	s := &state{
		cfg: &c,
		db:  dbQueries,
	}

	cmds := &commands{
		command: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", middlewareLoggedIn(handlerFeeds))
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))

	args := os.Args
	err = cmds.run(s, command{
		name: args[1],
		args: []string(args[2:]),
	})
	if err != nil {
		fmt.Println(fmt.Errorf("error running command: %w", err))
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, c command) error {
		if s.cfg.CurrentUserName == "" {
			return errors.New("not logged in")
		}
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			fmt.Println("You can't login to an account that doesn't exist!")
			return err
		}
		// Call the original handler function
		return handler(s, c, user)
	}
}

func handlerLogin(s *state, c command) error {
	if len(c.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	_, err := s.db.GetUserByName(context.Background(), c.args[0])
	if err != nil {
		fmt.Println("You can't login to an account that doesn't exist!")
		os.Exit(1)
	}

	if err := s.cfg.SetUser(c.args[0]); err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Println("User set successfully")
	return nil
}

func handlerRegister(s *state, c command) error {
	if len(c.args) == 0 {
		return fmt.Errorf("no username provided")
	}
	name := c.args[0]
	userParams := database.CreateUserParams{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}

	_, err := s.db.GetUserByName(context.Background(), name)
	if err == nil {
		fmt.Println("user already exists")
		os.Exit(1)
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}

	fmt.Println(user)

	return nil
}

func handlerReset(s *state, c command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("error resetting users: %w", err)
	}

	fmt.Println("Users reset successfully")
	return nil
}

func handlerListUsers(s *state, c command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch all users: %w", err)
	}

	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerAgg(s *state, c command) error {
	if len(c.args) == 0 {
		return fmt.Errorf("no time provided")
	}
	t := c.args[0]
	timeBetweenRequests, err := time.ParseDuration(t)
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}
	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		fmt.Println("Scraping feeds...")
		err := utils.ScrapeFeeds(s.db)
		if err != nil {
			return fmt.Errorf("error scraping feeds: %w", err)
		}
	}

}

func handlerAddFeed(s *state, c command, user database.User) error {
	if len(c.args) == 0 {
		fmt.Println("missing name and url")
		os.Exit(1)
	} else if len(c.args) == 1 {
		fmt.Println("missing url")
		os.Exit(1)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      c.args[0],
		Url:       c.args[1],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	f, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Printf("error creating feed: %w\n", err)
		os.Exit(1)
	}

	// once the new feed is created
	// the user will automatically follow it
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    f.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Printf("failed to follow newly created feed: %w\n", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", f)

	return nil
}

func handlerFeeds(s *state, c command, user database.User) error {
	feeds, err := s.db.ListFeeds(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("failed to fetch feeds: %w\n", err)
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		fmt.Println(feed.Name_2)
	}
	return nil
}

func handleFollow(s *state, c command, user database.User) error {
	url := c.args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		fmt.Printf("error querying feed: %w\n", err)
		os.Exit(1)
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Printf("error following feed with url: %w\n", err)
	}

	fmt.Println(feedFollow)

	return nil
}

func handleFollowing(s *state, c command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("failed to fetch follows for user %s: %w\n", user.Name, err)
		return err
	}
	fmt.Printf("%v\n", feeds)
	for _, feed := range feeds {
		fmt.Println(feed.Feedname)
	}
	return nil
}

func handleUnfollow(s *state, c command, user database.User) error {
	feedFollow := database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    c.args[0],
	}
	if err := s.db.DeleteFeedFollow(context.Background(), feedFollow); err != nil {
		fmt.Printf("error unfollowing feed: %w\n", err)
		os.Exit(1)
	}
	return nil
}

func handleBrowse(s *state, c command, user database.User) error {
	listFeeds, err := s.db.ListFeeds(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error fetching posts for feed: %w\n", err)
	}

	for _, feed := range listFeeds {
		postParams := database.GetPostsForFeedOfUserParams{
			FeedID: feed.ID,
			UserID: user.ID,
		}
		posts, err := s.db.GetPostsForFeedOfUser(context.Background(), postParams)
		if err != nil {
			return fmt.Errorf("error fetching list of posts for feed: %s. Error is: %w\n", feed.Name, err)
		}
		for _, post := range posts {
			fmt.Println(post.Title, post.Url, post.PublishedAt)
		}
	}

	return nil
}
