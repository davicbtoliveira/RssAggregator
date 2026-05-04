package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	feed := RSSFeed{}
	if err := xml.Unmarshal(body, &feed); err != nil {
		return &feed, err
	}

	//	DECODE SCAPED HTML
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, v := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(v.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}

	return &feed, nil
}
