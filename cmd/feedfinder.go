// Copyright (c) 2020, SpecializedGeneralist Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/SpecializedGeneralist/feedfinder"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) != 2 {
		printUsageAndExit()
	}

	url := os.Args[1]
	feeds, err := feedfinder.FindFeeds(url)
	if err != nil {
		log.Fatal(err)
	}
	for _, feed := range feeds {
		fmt.Println(feed)
	}
}

func printUsageAndExit() {
	log.Fatalf("feedfinder v%s\nUsage: %s <URL>",
		feedfinder.Version, os.Args[0])
}
