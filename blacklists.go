// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

import (
	"fmt"
	"net/url"
	"strings"
)

// HostBlacklist is a list of host names to be excluded from the process of
// searching for feeds.
//
// The values MUST be lowercase.
//
// To ensure precise matches, each entry in this list should be prefixed with
// a dot ".".
//
// It should contain domains for which you are sure no good feed URL can
// be extracted, or whose inclusion would cause an explosive amount of
// candidates, at least according to your specific set of initial URLs and
// crawling depth.
//
// Example: ".facebook.com" will cause the exclusion of URLs such as
// "https://facebook.com/..." or "https://m.facebook.com", but it will allow
// "https://thisisnotfacebook.com".
var HostBlacklist = []string{
	".facebook.com",
	".google.com",
	".instagram.com",
	".microsoft.com",
	".twitter.com",
	".youtube.com",
}

// URLPrefixBlacklist allows discarding invalid or undesired URLs based on
// their prefix.
//
// The values MUST be lowercase.
var URLPrefixBlacklist = []string{
	"javascript:",
	"mailto:",
	"tel:",
	"whatsapp:",
}

// removeBlacklistedURLs returns a new filtered list of URLs, discarding entries
// according to HostBlacklist and URLPrefixBlacklist. URLs that cannot be
// parsed are kept.
func removeBlacklistedURLs(urls []string) []string {
	filtered := make([]string, 0, len(urls))
outerLoop:
	for _, rawURL := range urls {
		lowerURL := strings.ToLower(rawURL)
		for _, p := range URLPrefixBlacklist {
			if strings.HasPrefix(lowerURL, p) {
				continue outerLoop
			}
		}

		if parsed, err := url.Parse(lowerURL); err == nil {
			host := fmt.Sprintf(".%s", parsed.Host)
			for _, h := range HostBlacklist {
				if strings.HasSuffix(host, h) {
					continue outerLoop
				}
			}
		}

		filtered = append(filtered, rawURL)
	}
	return filtered
}
