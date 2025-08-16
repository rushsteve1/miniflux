package form

import (
	"net/http"
	"strconv"
	"strings"

	"miniflux.app/v2/internal/locale"
)

type EntryForm struct {
	URL      string
	FeedID   int64
	Title    string
	Author   string
	Date     string
	Content  string
	Tags     []string
	Download bool
}

func (e *EntryForm) Validate() *locale.LocalizedError {
	// TODO: Implement validation logic
	return nil
}

func NewEntryForm(r *http.Request) *EntryForm {
	feedID, err := strconv.Atoi(r.FormValue("feed_id"))
	if err != nil {
		feedID = 0
	}

	return &EntryForm{
		URL:      r.FormValue("url"),
		FeedID:   int64(feedID),
		Download: r.FormValue("download") == "1",
		Title:    r.FormValue("title"),
		Author:   r.FormValue("author"),
		Date:     r.FormValue("date"),
		Content:  r.FormValue("content"),
		Tags:     strings.Split(r.FormValue("tags"), ","),
	}
}
