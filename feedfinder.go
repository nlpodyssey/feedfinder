// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

import (
	"encoding/xml"
	"github.com/nlpodyssey/feedfinder/html"
	"github.com/nlpodyssey/feedfinder/requests"
	"net/url"
)

// Version is the version of feedfinder.
const Version = "0.0.0"

// FindFeeds finds a list of RSS and Atom feeds from a web page.
func FindFeeds(url string) ([]string, error) {
	urls, err := fetchPageAndCollectURLCandidates(url)
	if err != nil {
		return nil, err
	}

	urls, err = resolveLinksReference(url, urls)
	if err != nil {
		return nil, err
	}

	urls, err = filterActualFeeds(urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func fetchPageAndCollectURLCandidates(url string) ([]string, error) {
	response, err := requests.Fetch(url)
	if err != nil {
		return nil, err
	}
	defer requests.CloseResponseBodyOrFatal(response)

	urls, err := html.CollectURLCandidates(response.Body)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func resolveLinksReference(rawBaseURL string, links []string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	resolvedLinks := make([]string, len(links))
	for i, link := range links {
		url, err := url.Parse(link)
		if err != nil {
			return nil, err
		}
		resolvedLinks[i] = baseURL.ResolveReference(url).String()
	}
	return resolvedLinks, nil
}

func filterActualFeeds(urls []string) ([]string, error) {
	result := make([]string, 0)
	for _, url := range urls {
		if isActualFeed(url) {
			result = append(result, url)
		}
	}
	return result, nil
}

func isActualFeed(url string) bool {
	response, err := requests.Fetch(url)
	if err != nil {
		return false
	}
	defer requests.CloseResponseBodyOrFatal(response)

	var rootNodeName xml.Name
	err = xml.NewDecoder(response.Body).Decode(&rootNodeName)
	if err != nil {
		return false
	}

	name := rootNodeName.Local
	return name == "rss" || name == "feed" || name == "rdf"
}
