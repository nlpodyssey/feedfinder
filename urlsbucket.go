// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package feedfinder

// urlsBucket is a naive implementation of a sorted set.
type urlsBucket struct {
	List []string
	Map  map[string]struct{}
}

// newURLsBucket creates a new urlsBucket object.
func newURLsBucket(initialURLs ...string) *urlsBucket {
	ub := &urlsBucket{
		List: make([]string, 0, len(initialURLs)),
		Map:  make(map[string]struct{}, len(initialURLs)),
	}
	ub.Put(initialURLs...)
	return ub
}

// Len returns the size of the set.
func (ub *urlsBucket) Len() int {
	return len(ub.List)
}

// Put inserts one or more URLs in the set, only if they don't already exist.
func (ub *urlsBucket) Put(urls ...string) {
	for _, u := range urls {
		if _, exists := ub.Map[u]; exists {
			continue
		}
		ub.List = append(ub.List, u)
		ub.Map[u] = struct{}{}
	}
}
