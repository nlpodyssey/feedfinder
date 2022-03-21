// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

import (
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"
)

// FindFeeds searches for RSS or Atom Feed URLs, starting from initialURL
// and recursively following links, up to maxDepth levels.
func FindFeeds(initialURL string, maxDepth int) ([]string, error) {
	if maxDepth < 1 {
		return nil, errors.New("invalid max depth value (must be >= 1)")
	}

	initialURL = strings.TrimSpace(initialURL)
	if len(initialURL) == 0 {
		return nil, errors.New("invalid blank initial URL")
	}

	log.Printf("[INFO] Initial URL: %#v, max depth: %d", initialURL, maxDepth)

	linksBucket := newURLsBucket(initialURL)
	feedsBucket := newURLsBucket()

	for depth, offset := 0, 0; depth <= maxDepth && offset < linksBucket.Len(); depth++ {
		log.Printf("[INFO] Depth %d", depth)

		initialLen := linksBucket.Len()
		curSize := initialLen - offset
		for i, curURL := range linksBucket.List[offset:initialLen] {
			if i%(curSize/10+1) == 0 {
				log.Printf("[INFO] Depth %d - URL %d of %d...", depth, i+1, curSize)
			}

			links, isFeed, err := processURL(curURL, depth == maxDepth)
			if err != nil {
				log.Printf("[ERROR] (%#v): %v", curURL, err)
				continue
			}
			linksBucket.Put(links...)
			if isFeed {
				feedsBucket.Put(curURL)
				log.Printf("[INFO] potential Feed found: %#v", curURL)
			}
		}
		offset = initialLen
	}

	return feedsBucket.List, nil
}

func processURL(curURL string, feedsOnly bool) (links []string, isFeed bool, err error) {
	if feedsOnly {
		matched, err := regexp.MatchString(`(rss|atom|rdf|xml|feed)`, strings.ToLower(curURL))
		if err != nil {
			return nil, false, err
		}
		if !matched {
			return nil, false, nil
		}
	}

	resp, err := httpGet(curURL)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		if e := resp.Body.Close(); e != nil {
			log.Printf("[ERROR] error closing response body: %v", e)
		}
	}()
	rawLinks, isFeed, err := parseBodyAndExtractLinks(resp.Body)
	if err != nil {
		return nil, false, err
	}

	resolvedLinks := resolveLinksReference(curURL, rawLinks)
	resolvedLinks = removeBlacklistedURLs(resolvedLinks)

	return resolvedLinks, isFeed, nil
}

func resolveLinksReference(rawBaseURL string, links []string) []string {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		log.Printf("[WARN] cannot parse URL %#v: %v", rawBaseURL, err)
		baseURL = nil
	}
	resolvedLinks := newURLsBucket()
	for _, link := range links {
		parsedURL, err := url.Parse(link)
		if err != nil {
			log.Printf("[WARN] cannot parse link URL %#v: %v", link, err)
			continue
		}
		resolvedLinks.Put(baseURL.ResolveReference(parsedURL).String())
	}
	return resolvedLinks.List
}
