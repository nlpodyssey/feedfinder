// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// parseBodyAndExtractLinks attempts to read and parse the content from the
// given reader.
//
// If the content can be parsed as HTML, the function collects the unique URL
// values from all "<a>" tags, and "<link>" tags with an XML-compatible type
// only. The unique URLs are returned as links and isFeed is false.
//
// If parts of the content are recognized as compatible with an Atom or RSS
// Feed, the list of returned links is empty, and isFeed is true.
func parseBodyAndExtractLinks(r io.Reader) (links []string, isFeed bool, err error) {
	urls := newURLsBucket()
	tokenizer := html.NewTokenizer(r)

	isFeed = false
	tagTokensCounter := 0
	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			if err := tokenizer.Err(); err != io.EOF {
				return nil, false, err
			}
			return urls.List, isFeed, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			if tagTokensCounter == 0 {
				tagName := strings.ToLower(token.Data)
				if tagName == "rss" || tagName == "feed" || tagName == "rdf" {
					isFeed = true
				}
			}
			tagTokensCounter++

			if u, ok := processTagToken(token); ok {
				urls.Put(u)
			}
		}
	}
}

func processTagToken(token html.Token) (string, bool) {
	switch token.Data {
	case "a":
		return getNormalizedAttributeValue(token, "href")
	case "link":
		return processLink(token)
	default:
		return "", false
	}
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

	if !hasRel || !hasType || !hasHref || rel != "alternate" || !feedLinkTypes[typeAttr] {
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
