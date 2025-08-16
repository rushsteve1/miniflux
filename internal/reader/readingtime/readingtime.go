// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

// Package readingtime provides a function to estimate the reading time of an article.
package readingtime

import (
	"errors"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"miniflux.app/v2/internal/config"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/proxyrotator"
	"miniflux.app/v2/internal/reader/fetcher"
	"miniflux.app/v2/internal/reader/sanitizer"
)

// EstimateReadingTime returns the estimated reading time of an article in minute.
func textReadingTime(content string, defaultReadingSpeed, cjkReadingSpeed int) int {
	sanitizedContent := sanitizer.StripTags(content)
	truncationPoint := min(len(sanitizedContent), 50)

	if isCJK(sanitizedContent[:truncationPoint]) {
		return int(math.Ceil(float64(utf8.RuneCountInString(sanitizedContent)) / float64(cjkReadingSpeed)))
	}
	return int(math.Ceil(float64(len(strings.Fields(sanitizedContent))) / float64(defaultReadingSpeed)))
}

func isCJK(text string) bool {
	totalCJK := 0

	for _, r := range text[:min(len(text), 50)] {
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Hangul, r) ||
			unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Yi, r) ||
			unicode.Is(unicode.Bopomofo, r) {
			totalCJK++
		}
	}

	// if at least 50% of the text is CJK, odds are that the text is in CJK.
	return totalCJK > len(text)/50
}

func fetchWatchTime(websiteURL, query string, isoDate bool) (int, error) {
	requestBuilder := fetcher.NewRequestBuilder()
	requestBuilder.WithTimeout(config.Opts.HTTPClientTimeout())
	requestBuilder.WithProxyRotator(proxyrotator.ProxyRotatorInstance)

	responseHandler := fetcher.NewResponseHandler(requestBuilder.ExecuteRequest(websiteURL))
	defer responseHandler.Close()

	if localizedError := responseHandler.LocalizedError(); localizedError != nil {
		slog.Warn("Unable to fetch watch time", slog.String("website_url", websiteURL), slog.Any("error", localizedError.Error()))
		return 0, localizedError.Error()
	}

	doc, docErr := goquery.NewDocumentFromReader(responseHandler.Body(config.Opts.HTTPClientMaxBodySize()))
	if docErr != nil {
		return 0, docErr
	}

	duration, exists := doc.FindMatcher(goquery.Single(query)).Attr("content")
	if !exists {
		return 0, errors.New("duration not found")
	}

	ret := 0
	if isoDate {
		parsedDuration, err := parseISO8601Duration(duration)
		if err != nil {
			return 0, fmt.Errorf("unable to parse iso duration %s: %v", duration, err)
		}
		ret = int(parsedDuration.Minutes())
	} else {
		parsedDuration, err := strconv.ParseInt(duration, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("unable to parse duration %s: %v", duration, err)
		}
		ret = int(parsedDuration / 60)
	}
	return ret, nil
}

// EstimateReadingTime updates the reading time of an entry based on its content.
func EstimateReadingTime(entry *model.Entry, user *model.User) int {
	if !user.ShowReadingTime {
		slog.Debug("Skip reading time estimation for this user", slog.Int64("user_id", user.ID))
		return 0
	}

	// Define watch time fetching scenarios
	watchTimeScenarios := [...]struct {
		shouldFetch func(*model.Entry) bool
		fetchFunc   func(string) (int, error)
		platform    string
	}{
		{shouldFetchYouTubeWatchTimeForSingleEntry, fetchYouTubeWatchTimeForSingleEntry, "YouTube"},
		{shouldFetchNebulaWatchTime, fetchNebulaWatchTime, "Nebula"},
		{shouldFetchOdyseeWatchTime, fetchOdyseeWatchTime, "Odysee"},
		{shouldFetchBilibiliWatchTime, fetchBilibiliWatchTime, "Bilibili"},
	}

	readingTime := 0

	// Iterate through scenarios and attempt to fetch watch time
	for _, scenario := range watchTimeScenarios {
		if scenario.shouldFetch(entry) {
			if entry.ID == 0 || entry.ReadingTime == 0 {
				if watchTime, err := scenario.fetchFunc(entry.URL); err != nil {
					slog.Warn("Unable to fetch watch time",
						slog.String("platform", scenario.platform),
						slog.Int64("user_id", user.ID),
						slog.Int64("entry_id", entry.ID),
						slog.String("entry_url", entry.URL),
						slog.Any("error", err),
					)
				} else {
					readingTime = watchTime
				}
			}
			break
		}
	}

	// Fallback to text-based reading time estimation
	if readingTime == 0 {
		return textReadingTime(entry.Content, user.DefaultReadingSpeed, user.CJKReadingSpeed)
	}
	return readingTime
}

// parseISO8601Duration parses a subset of ISO8601 durations, mainly for youtube video.
func parseISO8601Duration(duration string) (time.Duration, error) {
	after, ok := strings.CutPrefix(duration, "PT")
	if !ok {
		return 0, errors.New("the period doesn't start with PT")
	}

	var d time.Duration
	num := ""

	for _, char := range after {
		var val float64
		var err error

		switch char {
		case 'Y', 'W', 'D':
			return 0, fmt.Errorf("the '%c' specifier isn't supported", char)
		case 'H':
			if val, err = strconv.ParseFloat(num, 64); err != nil {
				return 0, err
			}
			d += time.Duration(val) * time.Hour
			num = ""
		case 'M':
			if val, err = strconv.ParseFloat(num, 64); err != nil {
				return 0, err
			}
			d += time.Duration(val) * time.Minute
			num = ""
		case 'S':
			if val, err = strconv.ParseFloat(num, 64); err != nil {
				return 0, err
			}
			d += time.Duration(val) * time.Second
			num = ""
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			num += string(char)
			continue
		default:
			return 0, errors.New("invalid character in the period")
		}
	}
	return d, nil
}
