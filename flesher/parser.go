package flesher

import (
	"html"
	"html/template"

	"github.com/gunnihinn/evil-feed-reader/provider"
	"github.com/mmcdole/gofeed"
)

func Flesh(p provider.Provider, resource string) (FeedResult, error) {
	reader, err := p(resource)
	if err != nil {
		return feedResult{}, err
	}
	defer reader.Close()

	feed, err := gofeed.NewParser().Parse(reader)
	if err != nil {
		return feedResult{}, err
	}

	result := feedResult{
		title: feed.Title,
		url:   feed.Link,
	}

	entries := make([]EntryResult, 0)
	for _, item := range feed.Items {
		var content string
		if item.Description != "" {
			content = item.Description
		} else {
			content = item.Content
		}

		var published string
		if item.Updated != "" {
			published = item.Updated
		} else {
			published = item.Published
		}

		entries = append(entries, entryResult{
			title:     html.UnescapeString(item.Title),
			url:       item.Link,
			content:   template.HTML(content),
			published: published,
		})
	}

	result.entries = entries

	return result, nil
}

// Parser parses a feed.
type Parser func([]byte) (FeedResult, error)
