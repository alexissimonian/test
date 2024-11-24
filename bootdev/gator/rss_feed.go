package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	if len(feedURL) == 0 {
		return &RSSFeed{}, fmt.Errorf("Please provide a valid URL.")
	}

	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Problem creating the request for RSS Feed: %v\n", err)
	}

	request.Header.Add("User-Agent", "gator")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error executing the request while fetching feed: %v\n", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error reading received rss response: %v\n", err)
	}

	rssFeed := RSSFeed{}
	if err = xml.Unmarshal(responseData, &rssFeed); err != nil {
		return &RSSFeed{}, fmt.Errorf("Error parsing response into a RSS Feed struct: %v\n", err)
	}

	cleanFeed(&rssFeed)

	return &rssFeed, nil
}

func cleanFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
}
