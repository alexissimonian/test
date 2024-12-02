package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexissimonian/test/bootdev/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAggregator(s *state, c command) error {
	if len(c.args) != 1 {
		return fmt.Errorf("Agg expects only one duration argument like 1s, 1m, 1h. Got %v args.\n", len(c.args))
	}

	timeBetweenRequests, err := time.ParseDuration(c.args[0])
	if err != nil {
		return fmt.Errorf("Invalid duration argument. Expects format like 1s, 1m, 1h. Err: %v\n", err)
	}
	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		feed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return fmt.Errorf("Error getting next feed to fetch: %v\n", err)
		}

		rssFeed, err := fetchFeed(context.Background(), feed.Url)
		if err != nil {
			return fmt.Errorf("Error fetching feed with url: %v. Err: %v\n", feed.Url, err)
		}

		if err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			ID:            feed.ID,
			LastFetchedAt: sql.NullTime{Time: time.Now().UTC()},
			UpdatedAt:     time.Now().UTC(),
		}); err != nil {
			return fmt.Errorf("Error marking feed as fetched: %v\n", err)
		}
		for _, item := range rssFeed.Channel.Item {
			publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				return fmt.Errorf("Problem parsing pubdate: %v\n", err)
			}
			if err = s.db.CreatePost(context.Background(), database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				Title:       item.Title,
				Url:         item.Link,
				Description: sql.NullString{String: item.Description, Valid: true},
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			}); err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					continue
				}
				fmt.Printf("Could not create post: %v\n", err)
			}
		}
	}
}

func handlerCreateFeed(s *state, c command, user database.User) error {
	if len(c.args) != 2 {
		return fmt.Errorf("Error, addfeed expects 2 arguments. Got: %v\n", len(c.args))
	}

	feedName := c.args[0]
	feedUrl := c.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("Error adding the feed to db: %v\n", err)
	}

	fmt.Printf("%v\n", feed)

	// follow feed just added
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Error following the feed you just added: %v\n", err)
	}

	return nil
}

func handlerListFeeds(s *state, c command) error {
	if len(c.args) != 0 {
		return fmt.Errorf("Error, feeds expects 0 arguments. Got: %v\n", len(c.args))
	}

	getFeedsRow, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting your feeds: %v\n", err)
	}

	for i := range getFeedsRow {
		fmt.Printf("%v\n%v\n%v\n",
			getFeedsRow[i].Name,
			getFeedsRow[i].Url,
			getFeedsRow[i].Username,
		)
	}

	return nil
}

func handlerFollowFeed(s *state, c command, user database.User) error {
	if len(c.args) != 1 {
		return fmt.Errorf("Follow expects one argument only. Got: %v\n", len(c.args))
	}

	feedURL := c.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("Error, this feed does not exist %v\n", err)
	}

	createFeedFollowRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Error when trying to make user: %v follow feed: %v\n", createFeedFollowRow.Username, createFeedFollowRow.Feedname)
	}

	fmt.Printf("%v is now following %v feed !\n", createFeedFollowRow.Username, createFeedFollowRow.Feedname)
	return nil
}

func handlerFollowingFeed(s *state, c command, currentUser database.User) error {
	if len(c.args) != 0 {

	}

	feedsFollowingByCurrentUser, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("Error getting feeds of user: %v. Err: %v\n", currentUser.Name, err)
	}

	for _, ffbcu := range feedsFollowingByCurrentUser {
		fmt.Printf("%v\n", ffbcu.Feedname)
	}

	return nil
}

func handlerUnfollowingFeed(s *state, c command, user database.User) error {
	if len(c.args) != 1 {
		return fmt.Errorf("Unfollow expects only one argument. Got: %v\n", len(c.args))
	}

	feedUrl := c.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("Something went wrong finding the feed to unfollow: %v\n", err)
	}

	err = s.db.RemoveFeedFollow(context.Background(), database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("Something went wrong unfollowing feed: %v for user: %v: err: %v", feed.Url, user.Name, err)
	}

	return nil
}

func handlerBrowsePost(s *state, c command, user database.User) error {
	if len(c.args) > 1 {
		return fmt.Errorf("Browse accepts no more than 1 args. Got: %v\n", len(c.args))
	}
	limit := 2
    if len(c.args) == 1 {
        l, err := strconv.ParseInt(c.args[0], 0, 32)
        if err != nil {
            return fmt.Errorf("Browse expects an optional limit that is, a number. Got: %v\n", c.args[0])
        }
        limit = int(l)
    }
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		return fmt.Errorf("Could not get posts from user: %v. Err: %v\n", user.Name, err)
	}

	for _, post := range posts {
		fmt.Printf("title: %v\ndescription: %v\nurl: %v\n", post.Title, post.Description.String, post.Url)
	}

	return nil
}
