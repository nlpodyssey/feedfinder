// Copyright (c) 2022, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/SpecializedGeneralist/feedfinder"
	"log"
	"os"
)

func main() {
	initialURL := flag.String("url", "", "initial URL to probe")
	maxDepth := flag.Int("depth", 2, "maximum depth of linked pages to follow (>= 1)")
	flag.Parse()

	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(1)
	}

	feeds, err := feedfinder.FindFeeds(*initialURL, *maxDepth)
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
}
