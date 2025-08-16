package obelisk

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
	"miniflux.app/v2/internal/reader/fetcher"
)

// Asset is asset that used in a web page.
type Asset struct {
	Data        []byte
	ContentType string
}

// Archiver is the core of obelisk, which used to download a
// web page then embeds its assets.
type Archiver struct {
	sync.RWMutex

	Cache map[string]Asset

	DisableJS     bool
	DisableCSS    bool
	DisableEmbeds bool
	DisableMedias bool

	requestBuilder *fetcher.RequestBuilder
}

// Archive starts archival process for the specified request.
// Returns the archival result, content type and error if there are any.
func (arc *Archiver) Archive(document *goquery.Document, url string) (string, error) {
	if arc.Cache == nil {
		arc.Cache = make(map[string]Asset)
	}

	// If it's HTML process it
	result, err := arc.processHTML(document, url, false)
	if err != nil {
		return "", err
	}

	return result, nil
}
