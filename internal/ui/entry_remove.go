// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "miniflux.app/v2/internal/ui"

import (
	"net/http"

	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/html"
	"miniflux.app/v2/internal/http/route"
)

func (h *handler) removeEntry(w http.ResponseWriter, r *http.Request) {
	entryID := request.RouteInt64Param(r, "entryID")

	query := h.store.NewEntryQueryBuilder(request.UserID(r))
	query.WithEntryID(entryID)
	entry, err := query.GetEntry()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	if entry == nil {
		html.NotFound(w, r)
		return
	}

	if err := h.store.RemoveEntry(entry.ID); err != nil {
		html.ServerError(w, r, err)
		return
	}

	html.Redirect(w, r, route.Path(h.router, "feeds"))
}
