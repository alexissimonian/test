package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alexissimonian/test/bootdev/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, c command) error {
    //if len(c.args) != 1 {
    //    return fmt.Errorf("Agg expects one argument. Got: %v\n", len(c.args))
    //}
    //feedURL := c.args[0]
    rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return fmt.Errorf("Something went wrong getting your feed: %v\n", err)
    }

    fmt.Printf("%v\n", rssFeed)
	return nil
}

func handlerAddFeed(s *state, c command) error {
    if len(c.args) != 2 {
        return fmt.Errorf("Error, addfeed expects 2 arguments. Got: %v\n", len(c.args))
    }

    user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
    if err != nil {
        return fmt.Errorf("Could not get current active user: %v\n", err)
    }

    feedName := c.args[0]
    feedUrl := c.args[1]

    feed,  err := s.db.AddFeed(context.Background(), database.AddFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now().UTC(),
        UpdatedAt: time.Now().UTC(),
        Name: feedName,
        Url: feedUrl,
        UserID: user.ID,
    })

    if err != nil {
        return fmt.Errorf("Error adding the feed to db: %v\n", err)
    }

    fmt.Printf("%v\n", feed)

    return nil
}
