// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package readingtime // import "miniflux.app/v2/internal/reader/readingtime"

import (
	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/urllib"
)

func shouldFetchOdyseeWatchTime(entry *model.Entry) bool {
	if !config.Opts.FetchOdyseeWatchTime() {
		return false
	}

	return urllib.DomainWithoutWWW(entry.URL) == "odysee.com"
}

func fetchOdyseeWatchTime(websiteURL string) (int, error) {
	return fetchWatchTime(websiteURL, `meta[property="og:video:duration"]`, false)
}
