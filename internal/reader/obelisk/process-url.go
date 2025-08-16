package obelisk

import (
	"context"
	"errors"
	"io"
	nurl "net/url"
	"strings"

	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/reader/fetcher"

	"github.com/PuerkitoBio/goquery"
)

var errSkippedURL = errors.New("skip processing url")

//nolint:gocyclo,unparam
func (arc *Archiver) processURL(ctx context.Context, url string, parentURL string, embedded ...bool) ([]byte, string, error) {
	// Parse embedded value
	isEmbedded := len(embedded) != 0 && embedded[0]

	// Make sure this URL is not empty, data or hash. If yes, just skip it.
	url = strings.TrimSpace(url)
	if url == "" || strings.HasPrefix(url, "data:") || strings.HasPrefix(url, "#") {
		return nil, "", errSkippedURL
	}

	// Parse URL to make sure it's valid request URL. If not, there might be
	// some error while preparing document, so just skip this URL
	parsedURL, err := nurl.ParseRequestURI(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Hostname() == "" {
		return nil, "", errSkippedURL
	}

	// Check in cache to see if this URL already processed
	arc.RLock()
	cache, cacheExist := arc.Cache[url]
	arc.RUnlock()

	if cacheExist {
		return cache.Data, cache.ContentType, nil
	}

	responseHandler := fetcher.NewResponseHandler(arc.requestBuilder.ExecuteRequest(url))
	defer responseHandler.Close()

	// Get content type
	contentType := responseHandler.ContentType()
	contentType = strings.TrimSpace(contentType)
	if contentType == "" {
		contentType = "text/plain"
	}

	// Read content of response body. If the downloaded file is HTML
	// or CSS it need to be processed again
	bodyReader := responseHandler.Body(config.Opts.HTTPClientMaxBodySize())
	var bodyContent []byte

	switch {
	case contentType == "text/html" && isEmbedded:
		document, err := goquery.NewDocumentFromReader(bodyReader)
		if err != nil {
			return nil, "", err
		}

		newHTML, err := arc.processHTML(document, url, false)
		if err == nil {
			bodyContent = s2b(newHTML)
		} else {
			return nil, "", err
		}

	case contentType == "text/css":
		newCSS, err := arc.processCSS(ctx, bodyReader, parsedURL)
		if err == nil {
			bodyContent = s2b(newCSS)
		} else {
			return nil, "", err
		}

	default:
		bodyContent, err = io.ReadAll(bodyReader)
		if err != nil {
			return nil, "", err
		}
	}

	// Save data URL to cache
	arc.Lock()
	arc.Cache[url] = Asset{
		Data:        bodyContent,
		ContentType: contentType,
	}
	arc.Unlock()

	return bodyContent, contentType, nil
}
