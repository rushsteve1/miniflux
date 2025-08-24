// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "miniflux.app/v2/internal/ui"

import (
	"net/http"
	"time"

	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/html"
	"miniflux.app/v2/internal/http/route"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/ui/form"
	"miniflux.app/v2/internal/ui/session"
	"miniflux.app/v2/internal/ui/view"
	"miniflux.app/v2/internal/validator"
)

func (h *handler) updateEntry(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	entryID := request.RouteInt64Param(r, "entryID")
	entryBuilder := h.store.NewEntryQueryBuilder(user.ID)
	entryBuilder.WithEntryID(entryID)
	entryBuilder.WithoutStatus(model.EntryStatusRemoved)

	entry, err := entryBuilder.GetEntry()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	if entry == nil {
		html.NotFound(w, r)
		return
	}

	entryForm := form.NewEntryForm(r)

	sess := session.New(h.store, request.SessionID(r))
	view := view.New(h.tpl, r, sess)
	view.Set("form", entryForm)
	view.Set("entry", entry)
	view.Set("menu", "feeds")
	view.Set("user", user)
	view.Set("countUnread", h.store.CountUnreadEntries(user.ID))
	view.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(user.ID))

	date, err := time.Parse(time.DateOnly, entryForm.Date)
	if err != nil {
		view.Set("errorMessage", err.Error()) // TODO: translate this
		html.OK(w, r, view.Render("update_entry"))
		return
	}

	entryRequest := &model.EntryUpdateRequest{
		URL:    model.OptionalString(entryForm.URL),
		Title:  model.OptionalString(entryForm.Title),
		Author: model.OptionalString(entryForm.Author),
		Date:   model.OptionalField(date),
		Tags:   entryForm.Tags,
	}

	if validationErr := validator.ValidateEntryModification(entryRequest); validationErr != nil {
		view.Set("errorMessage", validationErr.Error()) // TODO: translate this
		html.OK(w, r, view.Render("update_entry"))
		return
	}

	entryRequest.Patch(entry)
	if err := h.store.UpdateEntry(entry); err != nil {
		html.ServerError(w, r, err)
		return
	}

	html.Redirect(w, r, route.Path(h.router, "feedEntry", "feedID", entry.FeedID, "entryID", entry.ID))
}
