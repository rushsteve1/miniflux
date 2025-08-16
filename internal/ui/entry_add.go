// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "miniflux.app/v2/internal/ui"

import (
	"net/http"

	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/html"
	"miniflux.app/v2/internal/storage"
	"miniflux.app/v2/internal/ui/form"
	"miniflux.app/v2/internal/ui/session"
	"miniflux.app/v2/internal/ui/view"
)

func (h *handler) showAddEntryPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	queryBuilder := storage.NewFeedQueryBuilder(h.store, request.UserID(r))
	queryBuilder.WithManual(true)
	feeds, err := queryBuilder.GetFeeds()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	form := form.EntryForm{
		FeedID: request.QueryInt64Param(r, "feed_id", 0),
	}

	if form.FeedID == 0 && len(feeds) > 0 {
		form.FeedID = feeds[0].ID
	}

	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)
	view.Set("feeds", feeds)
	view.Set("menu", "search")
	view.Set("user", user)
	view.Set("countUnread", h.store.CountUnreadEntries(user.ID))
	view.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(user.ID))
	view.Set("defaultUserAgent", config.Opts.HTTPClientUserAgent())
	view.Set("form", form)
	view.Set("hasProxyConfigured", config.Opts.HasHTTPClientProxyURLConfigured())

	html.OK(w, r, view.Render("add_entry"))
}
