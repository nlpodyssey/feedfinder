// Copyright (c) 2020, NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package requests

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const timeout time.Duration = 2 * time.Minute

// Fetch performs an HTTP request. It includes minimal error handling.
func Fetch(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP response status %s", response.Status)
	}

	return response, nil
}

// CloseResponseBodyOrFatal closes the Body stream of the given HTTP response.
// It includes error handling, causing the program to abort if the operation
// fails.
func CloseResponseBodyOrFatal(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
