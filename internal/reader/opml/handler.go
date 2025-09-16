// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package opml // import "miniflux.app/v2/internal/reader/opml"

import (
	"fmt"
	"io"

	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/storage"
)

// Handler handles the logic for OPML import/export.
type Handler struct {
	store *storage.Storage
}

// Export exports user feeds to OPML.
func (h *Handler) Export(userID int64) (string, error) {
	feeds, err := h.store.Feeds(userID)
	if err != nil {
		return "", err
	}

	subscriptions := make([]subcription, 0, len(feeds))
	for _, feed := range feeds {
		title := "Uncategorized"
		if feed.Category != nil {
			title = feed.Category.Title
		}

		subscriptions = append(subscriptions, subcription{
			Title:        feed.Title,
			FeedURL:      feed.FeedURL,
			SiteURL:      feed.SiteURL,
			Description:  feed.Description,
			CategoryName: title,
		})
	}

	return serialize(subscriptions), nil
}

// Import parses and create feeds from an OPML import.
func (h *Handler) Import(userID int64, data io.Reader) error {
	subscriptions, err := parse(data)
	if err != nil {
		return err
	}

	for _, subscription := range subscriptions {
		if !h.store.FeedURLExists(userID, subscription.FeedURL) {
			var category *model.Category
			var err error

			if subscription.CategoryName != "" {
				category, err = h.store.CategoryByTitle(userID, subscription.CategoryName)
				if err != nil {
					return fmt.Errorf("opml: unable to search category by title: %w", err)
				}

				if category == nil {
					category, err = h.store.CreateCategory(userID, &model.CategoryCreationRequest{Title: subscription.CategoryName})
					if err != nil {
						return fmt.Errorf(`opml: unable to create this category: %q`, subscription.CategoryName)
					}
				}
			}

			feed := &model.Feed{
				UserID:      userID,
				Title:       subscription.Title,
				FeedURL:     subscription.FeedURL,
				SiteURL:     subscription.SiteURL,
				Description: subscription.Description,
				Category:    category,
			}

			if err := h.store.CreateFeed(feed); err != nil {
				return fmt.Errorf(`opml: unable to create this feed: %q`, subscription.FeedURL)
			}
		}
	}

	return nil
}

// NewHandler creates a new handler for OPML files.
func NewHandler(store *storage.Storage) *Handler {
	return &Handler{store: store}
}
