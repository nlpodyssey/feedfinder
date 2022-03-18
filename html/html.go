// Copyright (c) 2020, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

// CollectURLCandidates returns a URL list of possible RSS or Atom feeds.
// URLs are read from the attributes of "<link>" and "<a>" elements and
// filtered by very simple rules.
func CollectURLCandidates(r io.Reader) ([]string, error) {
	urls := make([]string, 0)
	t := html.NewTokenizer(r)

	addURLIfNew := func(url string) {
		for _, l := range urls {
			if l == url {
				return
			}
		}
		urls = append(urls, url)
	}

	for {
		switch t.Next() {
		case html.ErrorToken:
			if err := t.Err(); err != io.EOF {
				return nil, err
			}
			return urls, nil
		case html.SelfClosingTagToken, html.StartTagToken:
			url, urlOk, err := processTagToken(t.Token())
			if err != nil {
				return nil, err
			}
			if urlOk {
				addURLIfNew(url)
			}
		}
	}
}

func processTagToken(token html.Token) (string, bool, error) {
	switch token.Data {
	case "a":
		url, ok, err := processA(token)
		if err != nil {
			return "", false, err
		}
		return url, ok, nil
	case "link":
		url, ok := processLink(token)
		return url, ok, nil
	default:
		return "", false, nil
	}
}

func processA(token html.Token) (string, bool, error) {
	href, hasHref := getNormalizedAttributeValue(token, "href")
	if !hasHref {
		return "", false, nil
	}

	matched, err := regexp.MatchString(`(rss|rdf|xml|atom)`, href)
	if err != nil {
		return "", false, err
	}
	if !matched {
		return "", false, nil
	}

	return href, true, nil
}

var feedLinkTypes = map[string]bool{
	"application/atom+xml":   true,
	"application/rss+xml":    true,
	"application/x-atom+xml": true,
	"application/x.atom+xml": true,
	"text/xml":               true,
}

func processLink(token html.Token) (string, bool) {
	rel, hasRel := getNormalizedAttributeValue(token, "rel")
	typeAttr, hasType := getNormalizedAttributeValue(token, "type")
	href, hasHref := getNormalizedAttributeValue(token, "href")

	if !hasRel || !hasType || !hasHref ||
		rel != "alternate" || !feedLinkTypes[typeAttr] {
		return "", false
	}

	return href, true
}

func getNormalizedAttributeValue(token html.Token, key string) (string, bool) {
	attr, hasAttr := findTokenAttribute(token, key)
	if !hasAttr {
		return "", false
	}
	val := normalizeAttributeValue(attr.Val)
	if len(val) == 0 {
		return "", false
	}
	return val, true
}

func findTokenAttribute(token html.Token, key string) (html.Attribute, bool) {
	for _, attr := range token.Attr {
		if attr.Key == key {
			return attr, true
		}
	}
	return html.Attribute{}, false
}

func normalizeAttributeValue(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
