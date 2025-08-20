package ui

import (
	"net/http"

	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/json"
)

func (h *handler) scrollEntry(w http.ResponseWriter, r *http.Request) {
	entryID := request.RouteInt64Param(r, "entryID")
	scrollPercent := request.FormFloat32Value(r, "percent")
	if err := h.store.UpdateEntryScrollPercent(entryID, scrollPercent); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.OK(w, r, "OK")
}
