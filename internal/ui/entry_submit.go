package ui

import (
	"cmp"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/crypto"
	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response/html"
	"miniflux.app/v2/internal/http/route"
	"miniflux.app/v2/internal/locale"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/reader/processor"
	"miniflux.app/v2/internal/storage"
	"miniflux.app/v2/internal/ui/form"
	"miniflux.app/v2/internal/ui/session"
	"miniflux.app/v2/internal/ui/view"
)

func (h *handler) submitEntry(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		html.ServerError(w, r, err)
		return
	}
	printer := locale.NewPrinter(request.UserLanguage(r))

	queryBuilder := storage.NewFeedQueryBuilder(h.store, user.ID)
	queryBuilder.WithManual(true)
	feeds, err := queryBuilder.GetFeeds()
	if err != nil {
		html.ServerError(w, r, err)
		return
	}

	sess := session.New(h.store, request.SessionID(r))
	v := view.New(h.tpl, r, sess)
	v.Set("feeds", feeds)
	v.Set("menu", "search")
	v.Set("user", user)
	v.Set("countUnread", h.store.CountUnreadEntries(user.ID))
	v.Set("countErrorFeeds", h.store.CountUserFeedsWithErrors(user.ID))
	v.Set("defaultUserAgent", config.Opts.HTTPClientUserAgent())
	v.Set("hasProxyConfigured", config.Opts.HasHTTPClientProxyURLConfigured())

	entryForm := form.NewEntryForm(r)
	if validationErr := entryForm.Validate(); validationErr != nil {
		v.Set("form", entryForm)
		v.Set("errorMessage", validationErr.Translate(user.Language))
		html.OK(w, r, v.Render("add_entry"))
		return
	}

	var date time.Time
	if entryForm.Date != "" {
		date, err = time.Parse(time.DateOnly, entryForm.Date)
		if err != nil {
			html.ServerError(w, r, err)
			return
		}
	}

	entry := model.NewEntry()
	entry.Title = entryForm.Title
	entry.URL = entryForm.URL
	entry.Hash = crypto.SHA256(entry.URL)
	entry.UserID = user.ID
	entry.Date = cmp.Or(date, time.Now())
	entry.Tags = model.CleanTags(entryForm.Tags)
	entry.CreatedAt = time.Now()
	entry.ChangedAt = time.Now()

	if entryForm.FeedID == 0 {
		title := printer.Print("feed.manual")
		if len(feeds) > 0 {
			title += " " + strconv.FormatInt(int64(len(feeds))+1, 10)
		}

		feedUniqueString := fmt.Sprintf("manual-%d", time.Now().Unix())
		newFeed := &model.Feed{
			Manual:   true,
			Disabled: true,
			Crawler:  true,
			UserID:   entry.UserID,
			Title:    title,

			// These are not relevant but must be unique
			FeedURL: feedUniqueString,
			SiteURL: feedUniqueString,
		}

		err = h.store.CreateFeed(newFeed)
		if err != nil {
			v.Set("form", entryForm)
			v.Set("errorMessage", err.Error()) // TODO: localize error message
			html.OK(w, r, v.Render("add_entry"))
			return
		}

		entry.FeedID = newFeed.ID
		entry.Feed = newFeed

	} else {
		queryBuilder := storage.NewFeedQueryBuilder(h.store, user.ID)
		queryBuilder.WithManual(true)
		queryBuilder.WithFeedID(entryForm.FeedID)
		feed, err := queryBuilder.GetFeed()
		if err != nil {
			html.ServerError(w, r, err)
			return
		}
		entry.FeedID = feed.ID
		entry.Feed = feed
	}

	if entryForm.Download {
		// Proccess the entry now to fill fields
		if err := processor.ProcessEntryWebPage(entry.Feed, entry, user); err != nil {
			html.ServerError(w, r, err)
			return
		}
	}

	if err := h.store.CreateEntry(entry); err != nil {
		v.Set("form", entryForm)
		v.Set("errorMessage", err.Error()) // TODO: localize error message
		html.OK(w, r, v.Render("add_entry"))
		return
	}

	html.Redirect(w, r, route.Path(h.router, "feedEntry", "feedID", entry.FeedID, "entryID", entry.ID))
}
