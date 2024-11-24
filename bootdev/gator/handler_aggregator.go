package main

import (
	"context"
	"fmt"
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
